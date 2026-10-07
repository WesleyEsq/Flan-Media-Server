package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type attemptRecord struct {
	failCount    int
	lockoutUntil time.Time
}

type AuthService struct {
	userRepo      *repository.UserRepository
	db            *database.DB
	sessionSecret []byte

	// In-memory dual-key PIN lockout tracking: key = "ip:userID"
	lockoutMu   sync.Mutex
	lockouts    map[string]*attemptRecord
	cleanupOnce sync.Once

	// Bcrypt CPU protection gate: single slot serialized channel
	bcryptGate chan struct{}

	// Multi-process token version sync via SQLite PRAGMA data_version
	tokenCacheMu      sync.RWMutex
	tokenCache        map[int64]int
	lastDataVersion   int64
	lastDataCheckTime time.Time

	// First-run bootstrap token
	bootstrapMu    sync.Mutex
	bootstrapToken string
}

func NewAuthService(userRepo *repository.UserRepository, db *database.DB, sessionSecret string) *AuthService {
	s := &AuthService{
		userRepo:      userRepo,
		db:            db,
		sessionSecret: []byte(sessionSecret),
		lockouts:      make(map[string]*attemptRecord),
		bcryptGate:    make(chan struct{}, 1),
		tokenCache:    make(map[int64]int),
	}

	// Initial token cache load
	_ = s.refreshTokenCache(context.Background())
	if ver, err := db.GetDataVersion(); err == nil {
		s.lastDataVersion = ver
	}
	return s
}

func (s *AuthService) SetBootstrapToken(token string) {
	s.bootstrapMu.Lock()
	defer s.bootstrapMu.Unlock()
	s.bootstrapToken = token
}

func (s *AuthService) GetBootstrapToken() string {
	s.bootstrapMu.Lock()
	defer s.bootstrapMu.Unlock()
	return s.bootstrapToken
}

func (s *AuthService) VerifyBootstrapToken(token string) bool {
	s.bootstrapMu.Lock()
	defer s.bootstrapMu.Unlock()
	if s.bootstrapToken == "" || token == "" {
		return false
	}
	return hmac.Equal([]byte(s.bootstrapToken), []byte(token))
}

func (s *AuthService) ClearBootstrapToken() {
	s.bootstrapMu.Lock()
	defer s.bootstrapMu.Unlock()
	s.bootstrapToken = ""
}

// GenerateBootstrapToken creates a 6-character alphanumeric string
func GenerateBootstrapToken() string {
	const charset = "abcdefghjkmnpqrstuvwxyz23456789" // easy to read, no 0/O, 1/l
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "flan77"
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// HashPIN creates bcrypt hash of numeric PIN
func (s *AuthService) HashPIN(pin string) (string, error) {
	if len(pin) < 4 || len(pin) > 6 {
		return "", fmt.Errorf("%w: PIN must be 4 to 6 numeric digits", model.ErrInvalidInput)
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// AuthenticatePIN validates user PIN with dual-key lockout tracking and bcrypt worker gate
func (s *AuthService) AuthenticatePIN(ctx context.Context, clientIP string, userID int64, pin string) (*model.User, error) {
	lockKey := fmt.Sprintf("%s:%d", clientIP, userID)

	// 1. Check in-memory lockout
	s.lockoutMu.Lock()
	rec, exists := s.lockouts[lockKey]
	now := time.Now()
	if exists && now.Before(rec.lockoutUntil) {
		remaining := time.Until(rec.lockoutUntil).Round(time.Second)
		s.lockoutMu.Unlock()
		return nil, fmt.Errorf("%w: try again in %v", model.ErrLockedOut, remaining)
	}
	if exists && !now.Before(rec.lockoutUntil) && rec.failCount >= 5 {
		// Lockout expired, reset fail count
		rec.failCount = 0
	}
	failCount := 0
	if exists {
		failCount = rec.failCount
	}
	s.lockoutMu.Unlock()

	// 2. Pre-gate artificial delay on 4th failed attempt
	if failCount == 3 {
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// 3. Acquire bcrypt comparison worker slot with 5-second timeout
	select {
	case s.bcryptGate <- struct{}{}:
		defer func() { <-s.bcryptGate }()
	case <-time.After(5 * time.Second):
		return nil, model.ErrBcryptBusy
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// 4. Retrieve user from DB
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.recordFailedAttempt(lockKey)
		return nil, model.ErrUnauthorized
	}

	// 5. Compare bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PINHash), []byte(pin)); err != nil {
		s.recordFailedAttempt(lockKey)
		return nil, model.ErrUnauthorized
	}

	// 6. Login succeeded: clear failed attempts
	s.lockoutMu.Lock()
	delete(s.lockouts, lockKey)
	s.lockoutMu.Unlock()

	return user, nil
}

func (s *AuthService) recordFailedAttempt(lockKey string) {
	s.lockoutMu.Lock()
	defer s.lockoutMu.Unlock()

	rec, exists := s.lockouts[lockKey]
	if !exists {
		rec = &attemptRecord{failCount: 0}
		s.lockouts[lockKey] = rec
	}
	rec.failCount++

	if rec.failCount >= 5 {
		// Lockout progression: 5 min base with exponential backoff capped at 60 min
		exponent := rec.failCount - 5
		lockoutMinutes := math.Min(60, 5*math.Pow(2, float64(exponent)))
		rec.lockoutUntil = time.Now().Add(time.Duration(lockoutMinutes) * time.Minute)
	}
}

// Session Cookie Management: format userID:role:tokenVersion:issuedAt:expiresAt:signature
func (s *AuthService) CreateSessionCookie(user *model.User) (string, time.Time) {
	now := time.Now().Unix()
	expiresAt := time.Now().Add(30 * 24 * time.Hour).Unix() // 30-day lifetime
	payload := fmt.Sprintf("%d:%s:%d:%d:%d", user.ID, user.Role, user.TokenVersion, now, expiresAt)
	sig := s.sign(payload)
	cookieVal := fmt.Sprintf("%s:%s", payload, sig)
	return cookieVal, time.Unix(expiresAt, 0)
}

func (s *AuthService) ValidateSessionCookie(ctx context.Context, cookieVal string) (*model.User, bool, error) {
	parts := strings.Split(cookieVal, ":")
	if len(parts) != 6 {
		return nil, false, model.ErrUnauthorized
	}

	userIDStr, roleStr, tokenVersionStr, issuedAtStr, expiresAtStr, sigHex := parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]
	payload := fmt.Sprintf("%s:%s:%s:%s:%s", userIDStr, roleStr, tokenVersionStr, issuedAtStr, expiresAtStr)

	// Verify constant-time HMAC signature
	if !s.verify(payload, sigHex) {
		return nil, false, model.ErrUnauthorized
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, false, model.ErrUnauthorized
	}
	tokenVersion, err := strconv.Atoi(tokenVersionStr)
	if err != nil {
		return nil, false, model.ErrUnauthorized
	}
	issuedAt, err := strconv.ParseInt(issuedAtStr, 10, 64)
	if err != nil {
		return nil, false, model.ErrUnauthorized
	}
	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return nil, false, model.ErrUnauthorized
	}

	// Verify expiration
	now := time.Now().Unix()
	if now > expiresAt {
		return nil, false, model.ErrUnauthorized
	}

	// Check token version against cached version (refreshed via PRAGMA data_version)
	currentVersion, err := s.getCurrentTokenVersion(ctx, userID)
	if err != nil || currentVersion != tokenVersion {
		return nil, false, model.ErrUnauthorized
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, false, model.ErrUnauthorized
	}

	// Check rolling renewal: if older than 7 days
	needsRenewal := (now - issuedAt) > (7 * 24 * 3600)

	return user, needsRenewal, nil
}

// VLC/MPV Signed URLs: payload = media_type:file_id:user_id:token_version:exp
func (s *AuthService) GenerateSignedURL(mediaType model.MediaType, fileID, userID int64, tokenVersion int) string {
	exp := time.Now().Add(4 * time.Hour).Unix()
	payload := fmt.Sprintf("%s:%d:%d:%d:%d", mediaType, fileID, userID, tokenVersion, exp)
	sig := s.sign(payload)
	return fmt.Sprintf("/download/%s/%d?exp=%d&u=%d&sig=%s", mediaType, fileID, exp, userID, sig)
}

func (s *AuthService) VerifySignedURL(ctx context.Context, mediaType model.MediaType, fileID, userID int64, exp int64, sigHex string) bool {
	if time.Now().Unix() > exp {
		return false
	}

	tokenVersion, err := s.getCurrentTokenVersion(ctx, userID)
	if err != nil {
		return false
	}

	payload := fmt.Sprintf("%s:%d:%d:%d:%d", mediaType, fileID, userID, tokenVersion, exp)
	return s.verify(payload, sigHex)
}

func (s *AuthService) getCurrentTokenVersion(ctx context.Context, userID int64) (int, error) {
	s.checkAndRefreshDataVersion(ctx)

	s.tokenCacheMu.RLock()
	defer s.tokenCacheMu.RUnlock()
	version, exists := s.tokenCache[userID]
	if !exists {
		// User might be newly created; direct lookup
		u, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return 0, err
		}
		return u.TokenVersion, nil
	}
	return version, nil
}

func (s *AuthService) checkAndRefreshDataVersion(ctx context.Context) {
	s.tokenCacheMu.Lock()
	defer s.tokenCacheMu.Unlock()

	currentDataVersion, err := s.db.GetDataVersion()
	if err != nil {
		return
	}

	if currentDataVersion != s.lastDataVersion {
		_ = s.refreshTokenCacheLocked(ctx)
		s.lastDataVersion = currentDataVersion
	}
}

func (s *AuthService) refreshTokenCache(ctx context.Context) error {
	s.tokenCacheMu.Lock()
	defer s.tokenCacheMu.Unlock()
	return s.refreshTokenCacheLocked(ctx)
}

func (s *AuthService) refreshTokenCacheLocked(ctx context.Context) error {
	users, err := s.userRepo.ListAll(ctx)
	if err != nil {
		return err
	}
	newMap := make(map[int64]int, len(users))
	for _, u := range users {
		newMap[u.ID] = u.TokenVersion
	}
	s.tokenCache = newMap
	return nil
}

func (s *AuthService) sign(payload string) string {
	mac := hmac.New(sha256.New, s.sessionSecret)
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *AuthService) verify(payload string, sigHex string) bool {
	expectedSig := s.sign(payload)
	expectedBytes, err1 := hex.DecodeString(expectedSig)
	actualBytes, err2 := hex.DecodeString(sigHex)
	if err1 != nil || err2 != nil {
		return false
	}
	return hmac.Equal(expectedBytes, actualBytes)
}
