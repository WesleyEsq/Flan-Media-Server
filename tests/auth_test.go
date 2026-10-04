package tests

import (
	"context"
	"errors"
	"net/url"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
)

func setupTestDB(t *testing.T) (*database.DB, *repository.UserRepository, *service.AuthService) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.Connect(dbPath)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, db, "test-super-secret-session-key-32-chars-long!")

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db, userRepo, authSvc
}

func TestPINValidationAndHashing(t *testing.T) {
	_, _, authSvc := setupTestDB(t)

	// Test invalid lengths
	if _, err := authSvc.HashPIN("123"); !errors.Is(err, model.ErrInvalidInput) {
		t.Errorf("Expected ErrInvalidInput for 3-digit PIN, got %v", err)
	}
	if _, err := authSvc.HashPIN("1234567"); !errors.Is(err, model.ErrInvalidInput) {
		t.Errorf("Expected ErrInvalidInput for 7-digit PIN, got %v", err)
	}

	// Test valid 4-digit and 6-digit PIN
	hash4, err := authSvc.HashPIN("1234")
	if err != nil || len(hash4) == 0 {
		t.Fatalf("Failed to hash 4-digit PIN: %v", err)
	}
	hash6, err := authSvc.HashPIN("123456")
	if err != nil || len(hash6) == 0 {
		t.Fatalf("Failed to hash 6-digit PIN: %v", err)
	}
}

func createTestUser(t *testing.T, userRepo *repository.UserRepository, username, pinHash string, role model.Role, avatarIcon string) *model.User {
	u := &model.User{
		Username:   username,
		PINHash:    pinHash,
		Role:       role,
		AvatarIcon: avatarIcon,
	}
	id, err := userRepo.Create(context.Background(), u)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	u.ID = id
	u.TokenVersion = 1
	return u
}

func TestDualKeyLockout(t *testing.T) {
	_, userRepo, authSvc := setupTestDB(t)
	ctx := context.Background()

	hash, _ := authSvc.HashPIN("1234")
	user := createTestUser(t, userRepo, "tester", hash, model.RoleAdmin, "mascot")

	ip1 := "192.168.1.10"
	ip2 := "192.168.1.20"

	// 1st to 3rd failed attempts on ip1
	for i := 1; i <= 3; i++ {
		_, err := authSvc.AuthenticatePIN(ctx, ip1, user.ID, "9999")
		if !errors.Is(err, model.ErrUnauthorized) {
			t.Fatalf("Attempt %d: expected ErrUnauthorized, got %v", i, err)
		}
	}

	// 4th failed attempt triggers 2s delay
	start := time.Now()
	_, err := authSvc.AuthenticatePIN(ctx, ip1, user.ID, "9999")
	elapsed := time.Since(start)
	if !errors.Is(err, model.ErrUnauthorized) {
		t.Fatalf("4th attempt: expected ErrUnauthorized, got %v", err)
	}
	if elapsed < 1900*time.Millisecond {
		t.Errorf("Expected 4th attempt to incur ~2s pre-gate delay, elapsed: %v", elapsed)
	}

	// 5th failed attempt triggers lockout
	_, err = authSvc.AuthenticatePIN(ctx, ip1, user.ID, "9999")
	if !errors.Is(err, model.ErrUnauthorized) {
		t.Fatalf("5th attempt: expected ErrUnauthorized, got %v", err)
	}

	// 6th attempt should now be immediately rejected with ErrLockedOut without bcrypt calculation
	start = time.Now()
	_, err = authSvc.AuthenticatePIN(ctx, ip1, user.ID, "1234") // even with correct PIN!
	lockoutDuration := time.Since(start)
	if !errors.Is(err, model.ErrLockedOut) {
		t.Fatalf("Expected ErrLockedOut, got: %v", err)
	}
	if lockoutDuration > 100*time.Millisecond {
		t.Errorf("Lockout check should be instant in memory, took: %v", lockoutDuration)
	}

	// Dual-key verification: ip2 is NOT locked out for the same user!
	authenticatedUser, err := authSvc.AuthenticatePIN(ctx, ip2, user.ID, "1234")
	if err != nil {
		t.Fatalf("ip2 should not be locked out; failed with: %v", err)
	}
	if authenticatedUser.ID != user.ID {
		t.Errorf("User ID mismatch: got %d, want %d", authenticatedUser.ID, user.ID)
	}
}

func TestBcryptWorkerGate(t *testing.T) {
	_, userRepo, authSvc := setupTestDB(t)
	ctx := context.Background()

	hash, _ := authSvc.HashPIN("5555")
	user := createTestUser(t, userRepo, "gate_user", hash, model.RoleUser, "flan")

	// Launch multiple concurrent authentications to verify serialized execution without deadlocks
	var wg sync.WaitGroup
	errCh := make(chan error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ip := "10.0.0.1"
			// All pass correct PIN
			_, err := authSvc.AuthenticatePIN(ctx, ip, user.ID, "5555")
			if err != nil {
				errCh <- err
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("Concurrent AuthenticatePIN failed: %v", err)
	}
}

func TestSessionCookiesAndTokenVersionRevocation(t *testing.T) {
	_, userRepo, authSvc := setupTestDB(t)
	ctx := context.Background()

	hash, _ := authSvc.HashPIN("1234")
	user := createTestUser(t, userRepo, "session_user", hash, model.RoleAdmin, "cat")

	// Create session cookie
	cookieVal, exp := authSvc.CreateSessionCookie(user)
	if len(cookieVal) == 0 || exp.Before(time.Now()) {
		t.Fatalf("Invalid session cookie generated")
	}

	// Validate valid cookie
	validatedUser, renewalNeeded, err := authSvc.ValidateSessionCookie(ctx, cookieVal)
	if err != nil {
		t.Fatalf("Failed to validate authentic session cookie: %v", err)
	}
	if validatedUser.ID != user.ID {
		t.Errorf("User ID mismatch: got %d, want %d", validatedUser.ID, user.ID)
	}
	if renewalNeeded {
		t.Errorf("Renewal should not be needed for freshly issued session cookie")
	}

	// Tampered signature test
	tamperedCookie := cookieVal[:len(cookieVal)-4] + "ffff"
	_, _, err = authSvc.ValidateSessionCookie(ctx, tamperedCookie)
	if !errors.Is(err, model.ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized for tampered cookie, got %v", err)
	}

	// Revoke all sessions by incrementing token version
	err = userRepo.IncrementTokenVersion(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to increment token version: %v", err)
	}

	// Same cookie must now fail validation because token_version changed
	_, _, err = authSvc.ValidateSessionCookie(ctx, cookieVal)
	if !errors.Is(err, model.ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized after token version bump, got %v", err)
	}
}

func TestSignedVLCURLs(t *testing.T) {
	_, userRepo, authSvc := setupTestDB(t)
	ctx := context.Background()

	hash, _ := authSvc.HashPIN("1234")
	user := createTestUser(t, userRepo, "stream_user", hash, model.RoleUser, "star")

	signedURL := authSvc.GenerateSignedURL(model.MediaTypeVideo, 42, user.ID, user.TokenVersion)
	if signedURL == "" {
		t.Fatalf("Signed URL is empty")
	}

	parsed, err := url.Parse(signedURL)
	if err != nil {
		t.Fatalf("Failed to parse signed URL: %v", err)
	}

	q := parsed.Query()
	expStr := q.Get("exp")
	uStr := q.Get("u")
	sig := q.Get("sig")

	exp, _ := strconv.ParseInt(expStr, 10, 64)
	uID, _ := strconv.ParseInt(uStr, 10, 64)

	// Valid verification
	if !authSvc.VerifySignedURL(ctx, model.MediaTypeVideo, 42, uID, exp, sig) {
		t.Errorf("VerifySignedURL failed for valid signed URL")
	}

	// Tampered fileID
	if authSvc.VerifySignedURL(ctx, model.MediaTypeVideo, 999, uID, exp, sig) {
		t.Errorf("VerifySignedURL accepted tampered fileID")
	}

	// Expired timestamp
	if authSvc.VerifySignedURL(ctx, model.MediaTypeVideo, 42, uID, time.Now().Add(-1*time.Minute).Unix(), sig) {
		t.Errorf("VerifySignedURL accepted expired timestamp")
	}
}

func TestBootstrapToken(t *testing.T) {
	_, _, authSvc := setupTestDB(t)

	token := service.GenerateBootstrapToken()
	if len(token) != 6 {
		t.Fatalf("Bootstrap token should be 6 characters, got %d (%s)", len(token), token)
	}

	authSvc.SetBootstrapToken(token)

	if !authSvc.VerifyBootstrapToken(token) {
		t.Errorf("VerifyBootstrapToken failed with valid token")
	}

	if authSvc.VerifyBootstrapToken("invalid") {
		t.Errorf("VerifyBootstrapToken accepted invalid token")
	}

	authSvc.ClearBootstrapToken()
	if authSvc.VerifyBootstrapToken(token) {
		t.Errorf("VerifyBootstrapToken accepted token after ClearBootstrapToken")
	}
}
