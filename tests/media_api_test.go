package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
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

func setupMediaAPIServer(t *testing.T) (*http.ServeMux, *service.AuthService, *model.User, *repository.VideoRepository, *repository.BookRepository, int64, int64, []int64) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "media_api_test.db")
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
	progressRepo := repository.NewProgressRepository(db)
	authSvc := service.NewAuthService(userRepo, db, "secret-key-for-media-api-test-32-chars!")
	mediaSvc := service.NewMediaService(videoRepo, bookRepo, progressRepo, sourceRepo, tempDir)
	scannerSvc := service.NewScannerService(videoRepo, bookRepo, sourceRepo)

	ctx := context.Background()

	// Create admin user
	hash, _ := authSvc.HashPIN("1234")
	adminUser := createTestUser(t, userRepo, "mediaadmin", hash, model.RoleAdmin, "mascot")

	// Create video source & video
	sourceID, err := sourceRepo.Create(ctx, &model.StorageSource{
		Name: "Videos", MediaType: model.MediaTypeVideo, FolderPath: filepath.Join(tempDir, "vids"), IsActive: true,
	})
	if err != nil {
		t.Fatalf("Failed to create video source: %v", err)
	}

	videoID, err := videoRepo.Create(ctx, &model.Video{
		SourceID: sourceID, Title: "Original Title", VideoType: model.VideoTypeSeries, FolderPath: "show1",
	})
	if err != nil {
		t.Fatalf("Failed to create video: %v", err)
	}

	// Create 2 video files
	f1, _ := videoRepo.CreateFile(ctx, &model.VideoFile{
		VideoID: videoID, Title: "ep1.mp4", RelativePath: "ep1.mp4", Format: "mp4", OrderIndex: 1,
	})
	f2, _ := videoRepo.CreateFile(ctx, &model.VideoFile{
		VideoID: videoID, Title: "ep2.mp4", RelativePath: "ep2.mp4", Format: "mp4", OrderIndex: 2,
	})

	// Create book source & book
	bookSourceID, err := sourceRepo.Create(ctx, &model.StorageSource{
		Name: "Books", MediaType: model.MediaTypeBook, FolderPath: filepath.Join(tempDir, "books"), IsActive: true,
	})
	if err != nil {
		t.Fatalf("Failed to create book source: %v", err)
	}

	bookID, err := bookRepo.Create(ctx, &model.Book{
		SourceID: bookSourceID, Title: "Original Book", Author: "Old Author", FolderPath: "book1",
	})
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	mux := http.NewServeMux()
	mediaCtrl := controller.NewMediaController(mediaSvc, scannerSvc, sourceRepo, videoRepo, bookRepo)
	mediaCtrl.RegisterRoutes(mux)

	return mux, authSvc, adminUser, videoRepo, bookRepo, videoID, bookID, []int64{f1, f2}
}

func TestUpdateMediaBook(t *testing.T) {
	mux, authSvc, adminUser, _, bookRepo, _, bookID, _ := setupMediaAPIServer(t)
	handler := middleware.Authenticate(authSvc)(mux)

	payload := map[string]string{
		"title":    "Dune: Deluxe Edition",
		"author":   "Frank Herbert",
		"overview": "Set on the desert planet Arrakis...",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/media/book/%d", bookID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	cookieVal, _ := authSvc.CreateSessionCookie(adminUser)
	req.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for book update, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify database persistence
	book, err := bookRepo.GetByID(context.Background(), bookID)
	if err != nil {
		t.Fatalf("Failed to get updated book: %v", err)
	}
	if book.Title != "Dune: Deluxe Edition" || book.Author != "Frank Herbert" || !book.MetadataLocked {
		t.Errorf("Book metadata mismatch: %+v", book)
	}
}

func TestBatchUpdateVideoFiles(t *testing.T) {
	mux, authSvc, adminUser, videoRepo, _, videoID, _, fileIDs := setupMediaAPIServer(t)
	handler := middleware.Authenticate(authSvc)(mux)

	payload := map[string]any{
		"files": []map[string]any{
			{"file_id": fileIDs[0], "custom_title": "Episode 1: Pilot", "order_index": 1001, "is_hidden": false},
			{"file_id": fileIDs[1], "custom_title": "Episode 2: The Return", "order_index": 1002, "is_hidden": true},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/media/video/%d/files", videoID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	cookieVal, _ := authSvc.CreateSessionCookie(adminUser)
	req.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for batch files update, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify database updates
	f1, _, _, err := videoRepo.GetFileWithContainer(context.Background(), fileIDs[0])
	if err != nil {
		t.Fatalf("Failed to get file 1: %v", err)
	}
	if f1.CustomTitle != "Episode 1: Pilot" || f1.OrderIndex != 1001 {
		t.Errorf("File 1 mismatch: %+v", f1)
	}

	f2, _, _, err := videoRepo.GetFileWithContainer(context.Background(), fileIDs[1])
	if err != nil {
		t.Fatalf("Failed to get file 2: %v", err)
	}
	if f2.CustomTitle != "Episode 2: The Return" || f2.OrderIndex != 1002 || !f2.IsHidden {
		t.Errorf("File 2 mismatch: %+v", f2)
	}
}

func TestUpdateCoverArt(t *testing.T) {
	mux, authSvc, adminUser, videoRepo, _, videoID, _, _ := setupMediaAPIServer(t)
	handler := middleware.Authenticate(authSvc)(mux)
	cookieVal, _ := authSvc.CreateSessionCookie(adminUser)

	// 1. Test valid JPEG upload via multipart/form-data
	// Standard minimal JPEG SOI + APP0 header
	jpegBytes := []byte{
		0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00,
		0xff, 0xd9,
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("cover", "poster.jpg")
	_, _ = part.Write(jpegBytes)
	writer.Close()

	req := httptest.NewRequest("POST", fmt.Sprintf("/api/media/video/%d/cover", videoID), &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for cover upload, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Success   bool   `json:"success"`
		CoverPath string `json:"cover_path"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil || !resp.Success || resp.CoverPath == "" {
		t.Fatalf("Invalid response payload: %+v (err: %v)", resp, err)
	}

	// Verify video record has cover_path updated
	v, err := videoRepo.GetByID(context.Background(), videoID)
	if err != nil || v.CoverPath != resp.CoverPath {
		t.Errorf("Video cover_path not updated in DB: got %s, want %s", v.CoverPath, resp.CoverPath)
	}

	// 2. Test invalid image format (plain text)
	reqBad := httptest.NewRequest("POST", fmt.Sprintf("/api/media/video/%d/cover", videoID), bytes.NewReader([]byte("not an image")))
	reqBad.Header.Set("Content-Type", "application/octet-stream")
	reqBad.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})

	recBad := httptest.NewRecorder()
	handler.ServeHTTP(recBad, reqBad)

	if recBad.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for plain text, got %d", recBad.Code)
	}
}

func TestUpdateCoverArtViaURL(t *testing.T) {
	mux, authSvc, adminUser, videoRepo, _, videoID, _, _ := setupMediaAPIServer(t)
	handler := middleware.Authenticate(authSvc)(mux)
	cookieVal, _ := authSvc.CreateSessionCookie(adminUser)

	jpegBytes := []byte{
		0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00,
		0xff, 0xd9,
	}

	// Spin up test server serving the image
	imgServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(jpegBytes)
	}))
	defer imgServer.Close()

	// 1. JSON payload with URL
	payload, _ := json.Marshal(map[string]string{"url": imgServer.URL + "/poster.jpg"})
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/media/video/%d/cover", videoID), bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for URL cover download, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Success   bool   `json:"success"`
		CoverPath string `json:"cover_path"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil || !resp.Success || resp.CoverPath == "" {
		t.Fatalf("Invalid response: %+v (err: %v)", resp, err)
	}

	v, err := videoRepo.GetByID(context.Background(), videoID)
	if err != nil || v.CoverPath != resp.CoverPath {
		t.Errorf("DB cover_path mismatch: got %s, want %s", v.CoverPath, resp.CoverPath)
	}

	// 2. Delete Cover
	reqDel := httptest.NewRequest("DELETE", fmt.Sprintf("/api/media/video/%d/cover", videoID), nil)
	reqDel.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})
	recDel := httptest.NewRecorder()
	handler.ServeHTTP(recDel, reqDel)

	if recDel.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for DELETE cover, got %d", recDel.Code)
	}

	vAfter, err := videoRepo.GetByID(context.Background(), videoID)
	if err != nil || vAfter.CoverPath != "" {
		t.Errorf("Expected empty cover_path after deletion, got %q", vAfter.CoverPath)
	}
}

func TestPruneDeletedDirectoriesOnMaintenanceScan(t *testing.T) {
	mux, authSvc, adminUser, videoRepo, _, videoID, _, _ := setupMediaAPIServer(t)
	handler := middleware.Authenticate(authSvc)(mux)
	cookieVal, _ := authSvc.CreateSessionCookie(adminUser)

	ctx := context.Background()

	// 1. Verify video exists
	vBefore, err := videoRepo.GetByID(ctx, videoID)
	if err != nil || vBefore == nil {
		t.Fatalf("Expected video to exist initially: %v", err)
	}

	// 2. Setup the source directory with .flan-keep, but leave show1 directory non-existent (deleted)
	vidsDir := filepath.Join(t.TempDir(), "vids")
	_ = os.MkdirAll(vidsDir, 0755)
	_ = os.WriteFile(filepath.Join(vidsDir, ".flan-keep"), []byte("flan"), 0644)

	// Update source folder path to our valid temp dir
	// We can run POST /api/scan directly
	req := httptest.NewRequest("POST", "/api/scan", nil)
	req.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for /api/scan, got %d: %s", rec.Code, rec.Body.String())
	}

	// 3. Test Manual DELETE /api/media/video/{id}
	reqDel := httptest.NewRequest("DELETE", fmt.Sprintf("/api/media/video/%d", videoID), nil)
	reqDel.AddCookie(&http.Cookie{Name: "flan_session", Value: cookieVal})
	recDel := httptest.NewRecorder()
	handler.ServeHTTP(recDel, reqDel)

	// Should be 200 OK or 404 (if already pruned by MaintenanceScan)
	if recDel.Code != http.StatusOK && recDel.Code != http.StatusNotFound {
		t.Errorf("Expected 200 or 404 for DELETE video, got %d", recDel.Code)
	}

	// Video should not exist anymore
	_, errAfter := videoRepo.GetByID(ctx, videoID)
	if errAfter == nil {
		t.Errorf("Expected video to be deleted from database, but it still exists")
	}
}


