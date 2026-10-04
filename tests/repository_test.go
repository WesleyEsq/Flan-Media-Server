package tests

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
)

func setupRepos(t *testing.T) (*database.DB, *repository.UserRepository, *repository.SourceRepository, *repository.VideoRepository, *repository.BookRepository, *repository.ProgressRepository) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "repo_test.db")
	db, err := database.Connect(dbPath)
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	userRepo := repository.NewUserRepository(db)
	sourceRepo := repository.NewSourceRepository(db)
	videoRepo := repository.NewVideoRepository(db)
	bookRepo := repository.NewBookRepository(db)
	progressRepo := repository.NewProgressRepository(db)

	return db, userRepo, sourceRepo, videoRepo, bookRepo, progressRepo
}

func TestUserRepositoryCRUD(t *testing.T) {
	_, userRepo, _, _, _, _ := setupRepos(t)
	ctx := context.Background()

	// Initial count should be 0
	count, err := userRepo.Count(ctx)
	if err != nil || count != 0 {
		t.Fatalf("Expected count 0, got %d (err: %v)", count, err)
	}

	// Create user
	u := &model.User{
		Username:    "alice",
		DisplayName: "Alice Liddell",
		PINHash:     "hash123",
		Role:        model.RoleAdmin,
		AvatarIcon:  "cat",
	}
	id, err := userRepo.Create(ctx, u)
	if err != nil || id <= 0 {
		t.Fatalf("Failed to create user: %v", err)
	}
	u.ID = id

	// Get by ID
	fetched, err := userRepo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.Username != "alice" || fetched.DisplayName != "Alice Liddell" || fetched.Role != model.RoleAdmin {
		t.Errorf("Mismatch in fetched user: %+v", fetched)
	}

	// Get by Username
	byName, err := userRepo.GetByUsername(ctx, "alice")
	if err != nil || byName.ID != id {
		t.Errorf("GetByUsername failed: %+v (err: %v)", byName, err)
	}

	// Update user
	u.DisplayName = "Alice in Wonderland"
	u.AvatarIcon = "star"
	if err := userRepo.Update(ctx, u); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	updated, _ := userRepo.GetByID(ctx, id)
	if updated.DisplayName != "Alice in Wonderland" || updated.AvatarIcon != "star" {
		t.Errorf("Update did not persist: %+v", updated)
	}

	// Update PIN (bumps token_version)
	if err := userRepo.UpdatePIN(ctx, id, "newhash456"); err != nil {
		t.Fatalf("UpdatePIN failed: %v", err)
	}
	pinUpdated, _ := userRepo.GetByID(ctx, id)
	if pinUpdated.PINHash != "newhash456" || pinUpdated.TokenVersion != 2 {
		t.Errorf("PIN or token version not updated: %+v", pinUpdated)
	}

	// Increment token version
	if err := userRepo.IncrementTokenVersion(ctx, id); err != nil {
		t.Fatalf("IncrementTokenVersion failed: %v", err)
	}
	tokenUpdated, _ := userRepo.GetByID(ctx, id)
	if tokenUpdated.TokenVersion != 3 {
		t.Errorf("Expected token_version 3, got %d", tokenUpdated.TokenVersion)
	}

	// Delete user
	if err := userRepo.Delete(ctx, id); err != nil {
		t.Fatalf("Delete user failed: %v", err)
	}
	_, err = userRepo.GetByID(ctx, id)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}
}

func TestSourceRepositoryCRUD(t *testing.T) {
	_, _, sourceRepo, _, _, _ := setupRepos(t)
	ctx := context.Background()

	s := &model.StorageSource{
		Name:       "Movies Drive",
		MediaType:  model.MediaTypeVideo,
		FolderPath: "/mnt/storage/movies",
		IsActive:   true,
	}

	id, err := sourceRepo.Create(ctx, s)
	if err != nil || id <= 0 {
		t.Fatalf("Failed to create storage source: %v", err)
	}
	s.ID = id

	// Get by ID
	fetched, err := sourceRepo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.Name != "Movies Drive" || !fetched.IsActive {
		t.Errorf("Fetched source mismatch: %+v", fetched)
	}

	// Get by Path
	byPath, err := sourceRepo.GetByPath(ctx, "/mnt/storage/movies")
	if err != nil || byPath.ID != id {
		t.Errorf("GetByPath failed: %+v (err: %v)", byPath, err)
	}

	// SetActive
	if err := sourceRepo.SetActive(ctx, id, false); err != nil {
		t.Fatalf("SetActive failed: %v", err)
	}
	deactivated, _ := sourceRepo.GetByID(ctx, id)
	if deactivated.IsActive {
		t.Errorf("Expected IsActive to be false")
	}

	// UpdateName
	if err := sourceRepo.UpdateName(ctx, id, "Cinema Storage"); err != nil {
		t.Fatalf("UpdateName failed: %v", err)
	}
	renamed, _ := sourceRepo.GetByID(ctx, id)
	if renamed.Name != "Cinema Storage" {
		t.Errorf("Name not updated: got %s", renamed.Name)
	}

	// Delete
	if err := sourceRepo.Delete(ctx, id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = sourceRepo.GetByID(ctx, id)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestVideoRepositoryCRUDAndCascade(t *testing.T) {
	db, _, sourceRepo, videoRepo, _, _ := setupRepos(t)
	ctx := context.Background()

	// Create active source
	sourceID, err := sourceRepo.Create(ctx, &model.StorageSource{
		Name:       "Video Storage",
		MediaType:  model.MediaTypeVideo,
		FolderPath: "/media/video",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create source: %v", err)
	}

	// Create movie video
	movie := &model.Video{
		SourceID:       sourceID,
		Title:          "The Matrix",
		VideoType:      model.VideoTypeMovie,
		ReleaseYear:    1999,
		Overview:       "Neo discovers reality is a simulation.",
		FolderPath:     "The Matrix (1999)",
		MetadataLocked: true,
	}
	videoID, err := videoRepo.Create(ctx, movie)
	if err != nil || videoID <= 0 {
		t.Fatalf("Failed to create video: %v", err)
	}
	movie.ID = videoID

	// Create video file
	file := &model.VideoFile{
		VideoID:         videoID,
		Title:           "The Matrix.mp4",
		RelativePath:    "The Matrix (1999)/The Matrix.mp4",
		FileSize:        2147483648,
		MTime:           1700000000,
		Format:          "mp4",
		DurationSeconds: 8160,
	}
	fileID, err := videoRepo.CreateFile(ctx, file)
	if err != nil || fileID <= 0 {
		t.Fatalf("Failed to create video file: %v", err)
	}

	// List videos
	videos, err := videoRepo.List(ctx, "movie", false)
	if err != nil {
		t.Fatalf("List videos failed: %v", err)
	}
	if len(videos) != 1 || videos[0].Title != "The Matrix" {
		t.Fatalf("Expected 1 movie 'The Matrix', got %d items", len(videos))
	}

	// Search videos
	searchResults, err := videoRepo.Search(ctx, "simulation", false)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(searchResults) != 1 || searchResults[0].ID != videoID {
		t.Errorf("Expected search to find The Matrix via overview query")
	}

	// Verify GetFileWithContainer
	f, v, s, err := videoRepo.GetFileWithContainer(ctx, fileID)
	if err != nil {
		t.Fatalf("GetFileWithContainer failed: %v", err)
	}
	if f.ID != fileID || v.ID != videoID || s.ID != sourceID {
		t.Errorf("GetFileWithContainer mismatch: f=%d, v=%d, s=%d", f.ID, v.ID, s.ID)
	}

	// Delete video and test cascade trigger
	if err := videoRepo.Delete(ctx, videoID); err != nil {
		t.Fatalf("Failed to delete video: %v", err)
	}

	// Verify video_files row was deleted by trigger
	var fileCount int
	err = db.Reader.QueryRowContext(ctx, "SELECT COUNT(*) FROM video_files WHERE video_id = ?", videoID).Scan(&fileCount)
	if err != nil || fileCount != 0 {
		t.Errorf("Expected 0 video_files after cascade delete, got %d (err: %v)", fileCount, err)
	}
}

func TestBookRepositoryCRUDAndCascade(t *testing.T) {
	db, _, sourceRepo, _, bookRepo, _ := setupRepos(t)
	ctx := context.Background()

	sourceID, _ := sourceRepo.Create(ctx, &model.StorageSource{
		Name:       "Books Storage",
		MediaType:  model.MediaTypeBook,
		FolderPath: "/media/books",
		IsActive:   true,
	})

	book := &model.Book{
		SourceID:       sourceID,
		Title:          "Neuromancer",
		Author:         "William Gibson",
		Overview:       "The sky above the port was the color of television, tuned to a dead channel.",
		FolderPath:     "Neuromancer",
		MetadataLocked: true,
	}
	bookID, err := bookRepo.Create(ctx, book)
	if err != nil || bookID <= 0 {
		t.Fatalf("Failed to create book: %v", err)
	}

	bookFile := &model.BookFile{
		BookID:       bookID,
		Title:        "Neuromancer.epub",
		RelativePath: "Neuromancer/Neuromancer.epub",
		FileSize:     1048576,
		MTime:        1700000000,
		Format:       model.BookFormatEPUB,
	}
	bfID, err := bookRepo.CreateFile(ctx, bookFile)
	if err != nil || bfID <= 0 {
		t.Fatalf("Failed to create book file: %v", err)
	}

	// List books filtered by EPUB format
	books, err := bookRepo.List(ctx, "epub", false)
	if err != nil {
		t.Fatalf("List books failed: %v", err)
	}
	if len(books) != 1 || books[0].Title != "Neuromancer" {
		t.Fatalf("Expected 1 book 'Neuromancer', got %d", len(books))
	}

	// Search books
	results, err := bookRepo.Search(ctx, "Gibson", false)
	if err != nil || len(results) != 1 {
		t.Errorf("Search failed: got %d results (err: %v)", len(results), err)
	}

	// Delete book and test cascade trigger
	if err := bookRepo.Delete(ctx, bookID); err != nil {
		t.Fatalf("Delete book failed: %v", err)
	}

	var bfCount int
	err = db.Reader.QueryRowContext(ctx, "SELECT COUNT(*) FROM book_files WHERE book_id = ?", bookID).Scan(&bfCount)
	if err != nil || bfCount != 0 {
		t.Errorf("Expected 0 book_files after cascade delete, got %d (err: %v)", bfCount, err)
	}
}

func TestProgressRepositoryAndContinueWatching(t *testing.T) {
	_, userRepo, sourceRepo, videoRepo, _, progressRepo := setupRepos(t)
	ctx := context.Background()

	// 1. Create user
	u := &model.User{Username: "viewer", PINHash: "hash", Role: model.RoleUser, AvatarIcon: "ghost"}
	userID, _ := userRepo.Create(ctx, u)

	// 2. Create source, video, video_file
	sourceID, _ := sourceRepo.Create(ctx, &model.StorageSource{
		Name: "Media", MediaType: model.MediaTypeVideo, FolderPath: "/media", IsActive: true,
	})
	videoID, _ := videoRepo.Create(ctx, &model.Video{
		SourceID: sourceID, Title: "Interstellar", VideoType: model.VideoTypeMovie, FolderPath: "Interstellar",
	})
	fileID, _ := videoRepo.CreateFile(ctx, &model.VideoFile{
		VideoID: videoID, Title: "Interstellar.mp4", RelativePath: "Interstellar/Interstellar.mp4",
		DurationSeconds: 10140, Format: "mp4",
	})

	// 3. Upsert Progress at 45.5% (uncompleted)
	p := &model.Progress{
		UserID:       userID,
		MediaType:    model.MediaTypeVideo,
		FileID:       fileID,
		PositionData: "4613.7",
		Percentage:   45.5,
		IsFinished:   false,
	}
	if err := progressRepo.Upsert(ctx, p); err != nil {
		t.Fatalf("Failed to upsert progress: %v", err)
	}

	// 4. Retrieve progress
	saved, err := progressRepo.Get(ctx, userID, model.MediaTypeVideo, fileID)
	if err != nil {
		t.Fatalf("Get progress failed: %v", err)
	}
	if saved.Percentage != 45.5 || saved.IsFinished {
		t.Errorf("Saved progress mismatch: %+v", saved)
	}

	// 5. Query Continue Watching
	shelf, err := progressRepo.ListContinueWatching(ctx, userID, 10)
	if err != nil {
		t.Fatalf("ListContinueWatching failed: %v", err)
	}
	if len(shelf) != 1 || shelf[0].MediaID != videoID {
		t.Fatalf("Expected 1 item in Continue Watching shelf, got %d", len(shelf))
	}

	// 6. Mark finished (e.g. at 96%)
	p.Percentage = 96.0
	p.IsFinished = true
	if err := progressRepo.Upsert(ctx, p); err != nil {
		t.Fatalf("Failed to update finished progress: %v", err)
	}

	// 7. Shelf should now be empty because video is finished
	shelfAfter, err := progressRepo.ListContinueWatching(ctx, userID, 10)
	if err != nil {
		t.Fatalf("ListContinueWatching failed: %v", err)
	}
	if len(shelfAfter) != 0 {
		t.Errorf("Expected 0 items in Continue Watching for finished video, got %d", len(shelfAfter))
	}
}
