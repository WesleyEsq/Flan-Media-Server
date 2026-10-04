package tests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/WesleyEsq/Flan-Media-Server/internal/controller"
	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/middleware"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
)

func setupStreamServer(t *testing.T) (*http.ServeMux, *service.AuthService, *repository.UserRepository, *model.User, int64, string) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "stream_test.db")
	db, err := database.Connect(dbPath)
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	userRepo := repository.NewUserRepository(db)
	sourceRepo := repository.NewSourceRepository(db)
	videoRepo := repository.NewVideoRepository(db)
	bookRepo := repository.NewBookRepository(db)
	authSvc := service.NewAuthService(userRepo, db, "secret-key-for-stream-test-at-least-32-bytes")

	ctx := context.Background()

	// 1. Create admin user
	hash, _ := authSvc.HashPIN("1234")
	adminUser := createTestUser(t, userRepo, "streamadmin", hash, model.RoleAdmin, "mascot")

	// 2. Create media storage folder & test dummy video file
	mediaDir := filepath.Join(tempDir, "movies")
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		t.Fatalf("Failed to create media directory: %v", err)
	}
	_ = os.WriteFile(filepath.Join(mediaDir, ".flan-keep"), []byte("flan-keep"), 0644)

	// Create a dummy video file of 1024 bytes
	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	videoFilePath := filepath.Join(mediaDir, "test_movie.mp4")
	if err := os.WriteFile(videoFilePath, testData, 0644); err != nil {
		t.Fatalf("Failed to create dummy video file: %v", err)
	}

	// 3. Register source, video, and file in DB
	sourceID, err := sourceRepo.Create(ctx, &model.StorageSource{
		Name:       "Movies Source",
		MediaType:  model.MediaTypeVideo,
		FolderPath: mediaDir,
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create source: %v", err)
	}

	videoID, err := videoRepo.Create(ctx, &model.Video{
		SourceID:   sourceID,
		Title:      "Test Movie",
		VideoType:  model.VideoTypeMovie,
		FolderPath: "",
	})
	if err != nil {
		t.Fatalf("Failed to create video: %v", err)
	}

	fileID, err := videoRepo.CreateFile(ctx, &model.VideoFile{
		VideoID:         videoID,
		Title:           "test_movie.mp4",
		RelativePath:    "test_movie.mp4",
		FileSize:        1024,
		Format:          "mp4",
		DurationSeconds: 120,
	})
	if err != nil {
		t.Fatalf("Failed to create video file: %v", err)
	}

	mux := http.NewServeMux()
	streamCtrl := controller.NewStreamController(videoRepo, bookRepo, sourceRepo, userRepo, authSvc, tempDir)
	streamCtrl.RegisterRoutes(mux)

	return mux, authSvc, userRepo, adminUser, fileID, string(testData)
}

func TestZeroCopyRangeStreaming(t *testing.T) {
	mux, authSvc, _, user, fileID, expectedData := setupStreamServer(t)

	// Wrap mux with Authenticate middleware
	handler := middleware.Authenticate(authSvc)(mux)

	// 1. Test full content request
	req := httptest.NewRequest("GET", fmt.Sprintf("/stream/video/%d", fileID), nil)
	cookieVal, _ := authSvc.CreateSessionCookie(user)
	req.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for full content request, got %d", rec.Code)
	}
	if rec.Header().Get("Accept-Ranges") != "bytes" {
		t.Errorf("Missing Accept-Ranges: bytes header")
	}
	if rec.Body.Len() != 1024 {
		t.Errorf("Expected 1024 bytes body, got %d", rec.Body.Len())
	}

	// 2. Test Range request: bytes=0-15 (16 bytes)
	reqRange := httptest.NewRequest("GET", fmt.Sprintf("/stream/video/%d", fileID), nil)
	reqRange.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})
	reqRange.Header.Set("Range", "bytes=0-15")

	recRange := httptest.NewRecorder()
	handler.ServeHTTP(recRange, reqRange)

	if recRange.Code != http.StatusPartialContent {
		t.Fatalf("Expected 206 Partial Content, got %d", recRange.Code)
	}
	contentRange := recRange.Header().Get("Content-Range")
	if contentRange != "bytes 0-15/1024" {
		t.Errorf("Expected Content-Range: bytes 0-15/1024, got %q", contentRange)
	}
	if recRange.Body.Len() != 16 {
		t.Fatalf("Expected body length 16, got %d", recRange.Body.Len())
	}
	if string(recRange.Body.Bytes()) != expectedData[0:16] {
		t.Errorf("Partial content bytes mismatch")
	}

	// 3. Test Range request: bytes=512-1023 (second half)
	reqRange2 := httptest.NewRequest("GET", fmt.Sprintf("/stream/video/%d", fileID), nil)
	reqRange2.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})
	reqRange2.Header.Set("Range", "bytes=512-1023")

	recRange2 := httptest.NewRecorder()
	handler.ServeHTTP(recRange2, reqRange2)

	if recRange2.Code != http.StatusPartialContent {
		t.Fatalf("Expected 206 Partial Content, got %d", recRange2.Code)
	}
	if recRange2.Body.Len() != 512 {
		t.Fatalf("Expected 512 bytes, got %d", recRange2.Body.Len())
	}
	if string(recRange2.Body.Bytes()) != expectedData[512:1024] {
		t.Errorf("Second half partial content mismatch")
	}
}

func TestStreamUnauthorized(t *testing.T) {
	mux, authSvc, _, _, fileID, _ := setupStreamServer(t)
	handler := middleware.Authenticate(authSvc)(mux)

	// No session cookie provided
	req := httptest.NewRequest("GET", fmt.Sprintf("/stream/video/%d", fileID), nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized for unauthenticated stream request, got %d", rec.Code)
	}
}

func TestCSRFMiddleware(t *testing.T) {
	innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	csrfProtected := middleware.CSRFProtect(innerHandler)

	// 1. GET requests should always pass
	getReq := httptest.NewRequest("GET", "/api/test", nil)
	getRec := httptest.NewRecorder()
	csrfProtected.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Errorf("GET should bypass CSRF, got %d", getRec.Code)
	}

	// 2. POST with matching Origin should pass
	postReq := httptest.NewRequest("POST", "http://localhost:8080/api/action", nil)
	postReq.Host = "localhost:8080"
	postReq.Header.Set("Origin", "http://localhost:8080")
	postRec := httptest.NewRecorder()
	csrfProtected.ServeHTTP(postRec, postReq)
	if postRec.Code != http.StatusOK {
		t.Errorf("POST with valid Origin should pass, got %d", postRec.Code)
	}

	// 3. POST with Sec-Fetch-Site: cross-site should be rejected with 403 Forbidden
	crossReq := httptest.NewRequest("POST", "http://localhost:8080/api/action", nil)
	crossReq.Host = "localhost:8080"
	crossReq.Header.Set("Sec-Fetch-Site", "cross-site")
	crossRec := httptest.NewRecorder()
	csrfProtected.ServeHTTP(crossRec, crossReq)
	if crossRec.Code != http.StatusForbidden {
		t.Errorf("Cross-site POST should be rejected with 403, got %d", crossRec.Code)
	}

	// 4. POST with mismatched Origin should be rejected with 403 Forbidden
	badOriginReq := httptest.NewRequest("POST", "http://localhost:8080/api/action", nil)
	badOriginReq.Host = "localhost:8080"
	badOriginReq.Header.Set("Origin", "http://malicious-site.com")
	badOriginRec := httptest.NewRecorder()
	csrfProtected.ServeHTTP(badOriginRec, badOriginReq)
	if badOriginRec.Code != http.StatusForbidden {
		t.Errorf("Mismatched Origin POST should be rejected with 403, got %d", badOriginRec.Code)
	}
}

func TestRequireAdminMiddleware(t *testing.T) {
	_, authSvc, userRepo, adminUser, _, _ := setupStreamServer(t)

	// Create regular user in the same database
	regularUser := createTestUser(t, userRepo, "regular", "hash", model.RoleUser, "cat")

	adminEndpoint := middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("admin access granted"))
	}))

	// 1. Regular user on JSON API endpoint returns 403 Forbidden
	reqRegular := httptest.NewRequest("GET", "/api/admin/system", nil)
	reqRegular.Header.Set("Accept", "application/json")
	ctxRegular := context.WithValue(reqRegular.Context(), middleware.UserFromContext(reqRegular.Context()), regularUser)
	_ = ctxRegular

	// Inject regular user into context
	reqRegular = reqRegular.WithContext(context.WithValue(reqRegular.Context(), "flan_user", regularUser))
	// Note: in auth_filter.go, userContextKey is unexported, but Authenticate sets it from cookie:
	authHandler := middleware.Authenticate(authSvc)(adminEndpoint)

	reqRegWithCookie := httptest.NewRequest("GET", "/api/admin/system", nil)
	reqRegWithCookie.Header.Set("Accept", "application/json")
	regCookie, _ := authSvc.CreateSessionCookie(regularUser)
	reqRegWithCookie.AddCookie(&http.Cookie{Name: "flan_session", Value: regCookie})
	recReg := httptest.NewRecorder()
	authHandler.ServeHTTP(recReg, reqRegWithCookie)

	if recReg.Code != http.StatusForbidden {
		t.Errorf("Regular user should be rejected from admin endpoint with 403, got %d", recReg.Code)
	}

	// 2. Admin user receives 200 OK
	reqAdminWithCookie := httptest.NewRequest("GET", "/api/admin/system", nil)
	reqAdminWithCookie.Header.Set("Accept", "application/json")
	adminCookie, _ := authSvc.CreateSessionCookie(adminUser)
	reqAdminWithCookie.AddCookie(&http.Cookie{Name: "flan_session", Value: adminCookie})
	recAdmin := httptest.NewRecorder()
	authHandler.ServeHTTP(recAdmin, reqAdminWithCookie)

	if recAdmin.Code != http.StatusOK {
		t.Errorf("Admin user should have 200 OK on admin endpoint, got %d", recAdmin.Code)
	}
	body, _ := io.ReadAll(recAdmin.Body)
	if string(body) != "admin access granted" {
		t.Errorf("Expected 'admin access granted', got %q", string(body))
	}
}

func TestAvatarServing(t *testing.T) {
	mux, _, _, adminUser, _, _ := setupStreamServer(t)

	req := httptest.NewRequest("GET", fmt.Sprintf("/avatars/%d", adminUser.ID), nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for avatar, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "image/svg+xml" {
		t.Errorf("Expected image/svg+xml, got %q", ct)
	}

	body, _ := io.ReadAll(rec.Body)
	if len(body) == 0 {
		t.Fatalf("Avatar body should not be empty")
	}
}

func TestVLCPlaylistM3U(t *testing.T) {
	mux, authSvc, _, adminUser, fileID, _ := setupStreamServer(t)
	handler := middleware.Authenticate(authSvc)(mux)

	cookieVal, _ := authSvc.CreateSessionCookie(adminUser)
	req := httptest.NewRequest("GET", fmt.Sprintf("/stream/vlc/video/%d/playlist.m3u", fileID), nil)
	req.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for M3U playlist, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "audio/x-mpegurl; charset=utf-8" {
		t.Errorf("Expected audio/x-mpegurl, got %q", ct)
	}

	body, _ := io.ReadAll(rec.Body)
	strBody := string(body)
	if len(strBody) == 0 || strBody[:7] != "#EXTM3U" {
		t.Errorf("Expected #EXTM3U header in playlist, got: %s", strBody)
	}
}

func TestDownloadDisposition(t *testing.T) {
	mux, authSvc, _, adminUser, fileID, _ := setupStreamServer(t)
	handler := middleware.Authenticate(authSvc)(mux)
	cookieVal, _ := authSvc.CreateSessionCookie(adminUser)

	// 1. Without dl=1, should be inline
	reqInline := httptest.NewRequest("GET", fmt.Sprintf("/download/video/%d", fileID), nil)
	reqInline.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})
	recInline := httptest.NewRecorder()
	handler.ServeHTTP(recInline, reqInline)

	if recInline.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", recInline.Code)
	}
	dispInline := recInline.Header().Get("Content-Disposition")
	if dispInline[:6] != "inline" {
		t.Errorf("Expected inline disposition, got %q", dispInline)
	}

	// 2. With dl=1, should be attachment
	reqAttach := httptest.NewRequest("GET", fmt.Sprintf("/download/video/%d?dl=1", fileID), nil)
	reqAttach.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})
	recAttach := httptest.NewRecorder()
	handler.ServeHTTP(recAttach, reqAttach)

	if recAttach.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", recAttach.Code)
	}
	dispAttach := recAttach.Header().Get("Content-Disposition")
	if dispAttach[:10] != "attachment" {
		t.Errorf("Expected attachment disposition, got %q", dispAttach)
	}
}

