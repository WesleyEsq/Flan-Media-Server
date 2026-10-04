package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/WesleyEsq/Flan-Media-Server/internal/middleware"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
)

type MediaController struct {
	mediaService   *service.MediaService
	scannerService *service.ScannerService
	sourceRepo     *repository.SourceRepository
	videoRepo      *repository.VideoRepository
	bookRepo       *repository.BookRepository
}

func NewMediaController(
	mediaService *service.MediaService,
	scannerService *service.ScannerService,
	sourceRepo *repository.SourceRepository,
	videoRepo *repository.VideoRepository,
	bookRepo *repository.BookRepository,
) *MediaController {
	return &MediaController{
		mediaService:   mediaService,
		scannerService: scannerService,
		sourceRepo:     sourceRepo,
		videoRepo:      videoRepo,
		bookRepo:       bookRepo,
	}
}

func (c *MediaController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/sources", c.handleListSources)
	mux.HandleFunc("POST /api/sources", c.handleAddSource)
	mux.HandleFunc("PUT /api/sources/{id}", c.handleUpdateSourceName)
	mux.HandleFunc("DELETE /api/sources/{id}", c.handleDeleteSource)
	mux.HandleFunc("POST /api/sources/{id}/scan", c.handleScanSource)
	mux.HandleFunc("POST /api/sources/{id}/ingest", c.handleIngestSource)
	mux.HandleFunc("POST /api/scan", c.handleMaintenanceScan)
	mux.HandleFunc("GET /api/scan/status", c.handleScanStatus)

	mux.HandleFunc("GET /api/progress/{type}/{file_id}", c.handleGetProgress)
	mux.HandleFunc("POST /api/progress", c.handleSaveProgress)
	mux.HandleFunc("GET /api/search", c.handleSearch)

	mux.HandleFunc("PUT /api/media/{type}/{id}", c.handleUpdateMedia)
	mux.HandleFunc("PUT /api/media/video/{id}", c.handleUpdateVideo)
	mux.HandleFunc("POST /api/media/{type}/{id}/cover", c.handleUpdateCover)
	mux.HandleFunc("DELETE /api/media/{type}/{id}/cover", c.handleDeleteCover)
	mux.HandleFunc("PUT /api/media/video/{id}/files", c.handleBatchUpdateVideoFiles)
	mux.HandleFunc("DELETE /api/media/{type}/{id}", c.handleDeleteMedia)
}

func (c *MediaController) handleListSources(w http.ResponseWriter, r *http.Request) {
	sources, err := c.sourceRepo.ListAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to list sources", http.StatusInternalServerError)
		return
	}

	for _, s := range sources {
		var stat syscall.Statfs_t
		if err := syscall.Statfs(s.FolderPath, &stat); err == nil {
			s.FreeSpaceBytes = uint64(stat.Bavail) * uint64(stat.Bsize)
			s.TotalSpaceBytes = uint64(stat.Blocks) * uint64(stat.Bsize)
			s.IsOnline = true
		} else {
			s.IsOnline = false
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sources)
}

func (c *MediaController) handleAddSource(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Name       string          `json:"name"`
		MediaType  model.MediaType `json:"media_type"`
		FolderPath string          `json:"folder_path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.FolderPath = strings.TrimSpace(req.FolderPath)
	if req.Name == "" || req.FolderPath == "" {
		http.Error(w, "Name and Folder Path are required", http.StatusBadRequest)
		return
	}

	// Verify or create folder
	if err := os.MkdirAll(req.FolderPath, 0755); err != nil {
		http.Error(w, "Failed to create/access target directory on disk", http.StatusBadRequest)
		return
	}

	// Auto-initialize .flan-keep marker if not present
	markerPath := filepath.Join(req.FolderPath, ".flan-keep")
	if _, err := os.Stat(markerPath); os.IsNotExist(err) {
		_ = os.WriteFile(markerPath, []byte("flan-keep"), 0644)
	}

	src := &model.StorageSource{
		Name:       req.Name,
		MediaType:  req.MediaType,
		FolderPath: req.FolderPath,
		IsActive:   true,
	}

	id, err := c.sourceRepo.Create(r.Context(), src)
	if err != nil {
		http.Error(w, "Failed to save storage source (path may already exist)", http.StatusBadRequest)
		return
	}
	src.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(src)
}

func (c *MediaController) handleUpdateSourceName(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := c.sourceRepo.UpdateName(r.Context(), id, strings.TrimSpace(req.Name)); err != nil {
		http.Error(w, "Failed to update source name", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (c *MediaController) handleDeleteSource(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := c.sourceRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Failed to remove source", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (c *MediaController) handleScanSource(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	source, err := c.sourceRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	candidates, err := c.scannerService.ScanSource(r.Context(), source)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candidates)
}

func (c *MediaController) handleIngestSource(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req model.IngestionPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := c.scannerService.IngestApproved(r.Context(), id, req.Approved); err != nil {
		http.Error(w, "Failed committing items to database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (c *MediaController) handleMaintenanceScan(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	status, err := c.scannerService.MaintenanceScan(r.Context())
	if err != nil {
		http.Error(w, "Maintenance scan failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (c *MediaController) handleScanStatus(w http.ResponseWriter, r *http.Request) {
	status := c.scannerService.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (c *MediaController) handleGetProgress(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mediaType := model.MediaType(r.PathValue("type"))
	fileIDStr := r.PathValue("file_id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	p, err := c.mediaService.GetProgress(r.Context(), user.ID, mediaType, fileID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (c *MediaController) handleSaveProgress(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		MediaType    model.MediaType `json:"media_type"`
		FileID       int64           `json:"file_id"`
		PositionData string          `json:"position_data"`
		Percentage   float64         `json:"percentage"`
		IsFinished   bool            `json:"is_finished"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	p := &model.Progress{
		UserID:       user.ID,
		MediaType:    req.MediaType,
		FileID:       req.FileID,
		PositionData: req.PositionData,
		Percentage:   req.Percentage,
		IsFinished:   req.IsFinished,
	}

	if err := c.mediaService.SaveProgress(r.Context(), p); err != nil {
		http.Error(w, "Failed to save progress", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (c *MediaController) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"videos": []*model.Video{},
			"books":  []*model.Book{},
		})
		return
	}

	videos, books, err := c.mediaService.Search(r.Context(), q, false)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"videos": videos,
		"books":  books,
	})
}

func (c *MediaController) handleUpdateVideo(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string             `json:"title"`
		ReleaseYear int                `json:"release_year"`
		VideoType   model.VideoType    `json:"video_type"`
		Overview    string             `json:"overview"`
		Files       []*model.VideoFile `json:"files"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	video, err := c.videoRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	video.Title = strings.TrimSpace(req.Title)
	video.ReleaseYear = req.ReleaseYear
	video.VideoType = req.VideoType
	video.Overview = req.Overview

	if err := c.mediaService.UpdateVideoMetadata(r.Context(), video, req.Files); err != nil {
		http.Error(w, "Failed to update video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (c *MediaController) handleUpdateMedia(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	mediaType := r.PathValue("type")
	if mediaType == "video" {
		c.handleUpdateVideo(w, r)
		return
	}

	if mediaType == "book" {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Title    string `json:"title"`
			Author   string `json:"author"`
			Overview string `json:"overview"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		book, err := c.bookRepo.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}

		if strings.TrimSpace(req.Title) != "" {
			book.Title = strings.TrimSpace(req.Title)
		}
		book.Author = strings.TrimSpace(req.Author)
		book.Overview = strings.TrimSpace(req.Overview)
		book.MetadataLocked = true

		if err := c.mediaService.UpdateBookMetadata(r.Context(), book); err != nil {
			http.Error(w, "Failed to update book", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
		return
	}

	http.Error(w, "Unsupported media type", http.StatusBadRequest)
}

func (c *MediaController) handleBatchUpdateVideoFiles(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Verify video exists
	if _, err := c.videoRepo.GetByID(r.Context(), id); err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	var req struct {
		Files []*model.VideoFile `json:"files"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	for _, f := range req.Files {
		f.VideoID = id
	}

	if err := c.videoRepo.BatchUpdateFiles(r.Context(), req.Files); err != nil {
		http.Error(w, "Failed to update video files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func downloadImageFromURL(ctx context.Context, rawURL string) ([]byte, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, fmt.Errorf("invalid URL scheme, must be http or https")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "FlanMediaServer/1.0 (Artwork Fetcher)")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed fetching image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d fetching image", resp.StatusCode)
	}

	// Limit to max 5MB
	data, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 5*1024*1024 {
		return nil, fmt.Errorf("image exceeds 5MB limit")
	}
	return data, nil
}

func (c *MediaController) handleUpdateCover(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	mediaType := r.PathValue("type")
	if mediaType != "video" && mediaType != "book" {
		http.Error(w, "Invalid media type", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Max 5 MB for cover artwork
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	var imgBytes []byte
	var ext string

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
			http.Error(w, "File exceeds 5MB limit", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("cover")
		if err == nil {
			defer file.Close()
			imgBytes, err = io.ReadAll(file)
			if err != nil {
				http.Error(w, "Failed reading image data", http.StatusInternalServerError)
				return
			}
		} else {
			coverURL := strings.TrimSpace(r.FormValue("cover_url"))
			if coverURL != "" {
				imgBytes, err = downloadImageFromURL(r.Context(), coverURL)
				if err != nil {
					http.Error(w, "Failed downloading image from URL: "+err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				http.Error(w, "Missing cover file or cover_url in form", http.StatusBadRequest)
				return
			}
		}
	} else if strings.HasPrefix(contentType, "application/json") {
		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.URL) == "" {
			http.Error(w, "Missing url in JSON payload", http.StatusBadRequest)
			return
		}
		imgBytes, err = downloadImageFromURL(r.Context(), strings.TrimSpace(req.URL))
		if err != nil {
			http.Error(w, "Failed downloading image from URL: "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		var err error
		imgBytes, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed reading request body", http.StatusInternalServerError)
			return
		}
	}

	if len(imgBytes) == 0 {
		http.Error(w, "Empty cover image", http.StatusBadRequest)
		return
	}

	// Validate magic bytes
	detected := http.DetectContentType(imgBytes)
	switch detected {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	default:
		http.Error(w, "Unsupported image format. Allowed formats: JPEG, PNG, WebP", http.StatusBadRequest)
		return
	}

	relPath, err := c.mediaService.SaveCoverArt(r.Context(), model.MediaType(mediaType), id, imgBytes, ext)
	if err != nil {
		http.Error(w, "Failed to save cover artwork", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success":    true,
		"cover_path": relPath,
	})
}

func (c *MediaController) handleDeleteCover(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	mediaType := r.PathValue("type")
	if mediaType != "video" && mediaType != "book" {
		http.Error(w, "Invalid media type", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := c.mediaService.RemoveCoverArt(r.Context(), model.MediaType(mediaType), id); err != nil {
		http.Error(w, "Failed to remove cover artwork", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (c *MediaController) handleDeleteMedia(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	mediaType := r.PathValue("type")
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if mediaType == "video" {
		err = c.videoRepo.Delete(r.Context(), id)
	} else {
		err = c.bookRepo.Delete(r.Context(), id)
	}

	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
