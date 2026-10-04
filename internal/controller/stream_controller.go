package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/WesleyEsq/Flan-Media-Server/internal/middleware"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
	"github.com/WesleyEsq/Flan-Media-Server/web"
)

type StreamController struct {
	videoRepo   *repository.VideoRepository
	bookRepo    *repository.BookRepository
	sourceRepo  *repository.SourceRepository
	userRepo    *repository.UserRepository
	authService *service.AuthService
	dataDir     string
}

func NewStreamController(
	videoRepo *repository.VideoRepository,
	bookRepo *repository.BookRepository,
	sourceRepo *repository.SourceRepository,
	userRepo *repository.UserRepository,
	authService *service.AuthService,
	dataDir string,
) *StreamController {
	return &StreamController{
		videoRepo:   videoRepo,
		bookRepo:    bookRepo,
		sourceRepo:  sourceRepo,
		userRepo:    userRepo,
		authService: authService,
		dataDir:     dataDir,
	}
}

func (c *StreamController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /stream/video/{file_id}", c.handleStreamVideo)
	mux.HandleFunc("GET /stream/book/{file_id}", c.handleStreamBook)
	mux.HandleFunc("GET /stream/subtitles/{file_id}/{track_id}", c.handleStreamSubtitles)
	mux.HandleFunc("GET /download/{type}/{file_id}", c.handleDownload)
	mux.HandleFunc("GET /stream/vlc/{type}/{file_id}/playlist.m3u", c.handleVLCPlaylist)
	mux.HandleFunc("GET /covers/{type}/{id}", c.handleCover)
	mux.HandleFunc("GET /avatars/{user_id}", c.handleAvatar)
	mux.HandleFunc("POST /api/upload", c.handleDirectUpload)
}

func (c *StreamController) handleStreamVideo(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fileIDStr := r.PathValue("file_id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	file, video, source, err := c.videoRepo.GetFileWithContainer(r.Context(), fileID)
	if err != nil || file == nil || video == nil || source == nil {
		http.Error(w, "Video file not found", http.StatusNotFound)
		return
	}

	if !source.IsActive || file.IsMissing {
		http.Error(w, "Media currently unavailable", http.StatusServiceUnavailable)
		return
	}

	filePath := resolveDiskPath(source.FolderPath, video.FolderPath, file.RelativePath)
	if err := validateCanonicalPath(source.FolderPath, filePath); err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Unable to read media file", http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "Unable to read media file info", http.StatusInternalServerError)
		return
	}

	// Content-Type mapping
	switch strings.ToLower(file.Format) {
	case "mp4":
		w.Header().Set("Content-Type", "video/mp4")
	case "webm":
		w.Header().Set("Content-Type", "video/webm")
	case "mkv":
		w.Header().Set("Content-Type", "video/x-matroska")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Header().Set("Accept-Ranges", "bytes")
	// Delegate directly to http.ServeContent for zero-copy Linux sendfile Range delivery
	http.ServeContent(w, r, file.Title, stat.ModTime(), f)
}

func (c *StreamController) handleStreamBook(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fileIDStr := r.PathValue("file_id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	file, book, source, err := c.bookRepo.GetFileWithContainer(r.Context(), fileID)
	if err != nil || file == nil || book == nil || source == nil {
		http.Error(w, "Book file not found", http.StatusNotFound)
		return
	}

	if !source.IsActive || file.IsMissing {
		http.Error(w, "Media currently unavailable", http.StatusServiceUnavailable)
		return
	}

	filePath := resolveDiskPath(source.FolderPath, book.FolderPath, file.RelativePath)
	if err := validateCanonicalPath(source.FolderPath, filePath); err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Unable to read book file", http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "Unable to read book file info", http.StatusInternalServerError)
		return
	}

	if file.Format == model.BookFormatPDF {
		w.Header().Set("Content-Type", "application/pdf")
	} else {
		w.Header().Set("Content-Type", "application/epub+zip")
	}

	w.Header().Set("Accept-Ranges", "bytes")
	http.ServeContent(w, r, file.Title, stat.ModTime(), f)
}

func (c *StreamController) handleStreamSubtitles(w http.ResponseWriter, r *http.Request) {
	fileIDStr := r.PathValue("file_id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	file, video, source, err := c.videoRepo.GetFileWithContainer(r.Context(), fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	videoPath := resolveDiskPath(source.FolderPath, video.FolderPath, file.RelativePath)
	dir := filepath.Dir(videoPath)
	base := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))

	// Search for sidecar subtitles: e.g. base.en.srt, base.srt, base.vtt
	candidates := []string{
		filepath.Join(dir, base+".vtt"),
		filepath.Join(dir, base+".srt"),
		filepath.Join(dir, base+".en.vtt"),
		filepath.Join(dir, base+".en.srt"),
	}

	for _, p := range candidates {
		if data, err := os.ReadFile(p); err == nil {
			w.Header().Set("Content-Type", "text/vtt; charset=utf-8")
			if strings.HasSuffix(p, ".srt") {
				w.Write(service.ConvertSRTToWebVTT(data))
			} else {
				w.Write(data)
			}
			return
		}
	}

	http.Error(w, "Subtitles not found", http.StatusNotFound)
}

func (c *StreamController) handleDownload(w http.ResponseWriter, r *http.Request) {
	mediaType := model.MediaType(r.PathValue("type"))
	fileIDStr := r.PathValue("file_id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Verify authentication: either session cookie OR signed URL parameters
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		// Check signed URL
		q := r.URL.Query()
		expStr := q.Get("exp")
		uStr := q.Get("u")
		sig := q.Get("sig")

		if expStr == "" || uStr == "" || sig == "" {
			http.Error(w, "Unauthorized: valid session or signed parameters required", http.StatusUnauthorized)
			return
		}

		exp, err1 := strconv.ParseInt(expStr, 10, 64)
		u, err2 := strconv.ParseInt(uStr, 10, 64)
		if err1 != nil || err2 != nil {
			http.Error(w, "Invalid signature parameters", http.StatusBadRequest)
			return
		}

		if !c.authService.VerifySignedURL(r.Context(), mediaType, fileID, u, exp, sig) {
			http.Error(w, "Unauthorized: invalid or expired stream link", http.StatusUnauthorized)
			return
		}
	}

	var filePath, filename, format string
	if mediaType == model.MediaTypeVideo {
		file, video, source, err := c.videoRepo.GetFileWithContainer(r.Context(), fileID)
		if err != nil || file == nil || video == nil || source == nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		filePath = resolveDiskPath(source.FolderPath, video.FolderPath, file.RelativePath)
		filename = fmt.Sprintf("%s.%s", file.DisplayTitle(), file.Format)
		format = file.Format
	} else {
		file, book, source, err := c.bookRepo.GetFileWithContainer(r.Context(), fileID)
		if err != nil || file == nil || book == nil || source == nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		filePath = resolveDiskPath(source.FolderPath, book.FolderPath, file.RelativePath)
		filename = fmt.Sprintf("%s.%s", file.DisplayTitle(), file.Format)
		format = string(file.Format)
	}

	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not accessible on disk", http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "Unable to inspect file", http.StatusInternalServerError)
		return
	}

	// Allow byte-range requests for external streaming
	w.Header().Set("Accept-Ranges", "bytes")
	if r.URL.Query().Get("dl") == "1" || r.URL.Query().Get("download") == "1" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	} else {
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	}

	switch strings.ToLower(format) {
	case "mp4":
		w.Header().Set("Content-Type", "video/mp4")
	case "mkv":
		w.Header().Set("Content-Type", "video/x-matroska")
	case "webm":
		w.Header().Set("Content-Type", "video/webm")
	case "pdf":
		w.Header().Set("Content-Type", "application/pdf")
	case "epub":
		w.Header().Set("Content-Type", "application/epub+zip")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	http.ServeContent(w, r, filename, stat.ModTime(), f)
}

func (c *StreamController) handleCover(w http.ResponseWriter, r *http.Request) {
	mediaType := r.PathValue("type")
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var coverPath string
	if mediaType == "video" {
		v, err := c.videoRepo.GetByID(r.Context(), id)
		if err == nil && v != nil {
			coverPath = v.CoverPath
		}
	} else if mediaType == "book" {
		b, err := c.bookRepo.GetByID(r.Context(), id)
		if err == nil && b != nil {
			coverPath = b.CoverPath
		}
	}

	if coverPath == "" {
		http.NotFound(w, r)
		return
	}

	// Long-lived client caching
	w.Header().Set("Cache-Control", "public, max-age=604800")
	http.ServeFile(w, r, coverPath)
}

func (c *StreamController) handleAvatar(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.userRepo.GetByID(r.Context(), userID)
	if err != nil || user == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")

	if user.AvatarPath != "" {
		// Custom photo upload
		if stat, err := os.Stat(user.AvatarPath); err == nil && !stat.IsDir() {
			http.ServeFile(w, r, user.AvatarPath)
			return
		}
	}

	// Fallback to preset SVG from embedded StaticFS
	iconName := user.AvatarIcon
	if iconName == "" || iconName == "default" {
		iconName = "flan"
	}

	data, err := web.StaticFS.ReadFile("static/assets/avatars/" + iconName + ".svg")
	if err != nil {
		// Fallback to flan.svg
		data, err = web.StaticFS.ReadFile("static/assets/avatars/flan.svg")
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(data)
}

func (c *StreamController) handleVLCPlaylist(w http.ResponseWriter, r *http.Request) {
	mediaType := model.MediaType(r.PathValue("type"))
	fileIDStr := r.PathValue("file_id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	user := middleware.UserFromContext(r.Context())
	if user == nil {
		q := r.URL.Query()
		expStr := q.Get("exp")
		uStr := q.Get("u")
		sig := q.Get("sig")
		if expStr == "" || uStr == "" || sig == "" {
			http.Error(w, "Unauthorized: valid session or signed parameters required", http.StatusUnauthorized)
			return
		}
		exp, err1 := strconv.ParseInt(expStr, 10, 64)
		u, err2 := strconv.ParseInt(uStr, 10, 64)
		if err1 != nil || err2 != nil || !c.authService.VerifySignedURL(r.Context(), mediaType, fileID, u, exp, sig) {
			http.Error(w, "Unauthorized: invalid or expired stream link", http.StatusUnauthorized)
			return
		}
	}

	var title string
	if mediaType == model.MediaTypeVideo {
		f, err := c.videoRepo.GetFileByID(r.Context(), fileID)
		if err != nil || f == nil {
			http.NotFound(w, r)
			return
		}
		title = f.DisplayTitle()
	} else {
		http.Error(w, "M3U streaming only supported for video", http.StatusBadRequest)
		return
	}

	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := r.Host

	streamURL := fmt.Sprintf("%s://%s/download/%s/%d", scheme, host, mediaType, fileID)
	if r.URL.RawQuery != "" {
		streamURL += "?" + r.URL.RawQuery
	}

	w.Header().Set("Content-Type", "audio/x-mpegurl; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", "stream-"+fileIDStr+".m3u"))

	playlist := fmt.Sprintf("#EXTM3U\n#EXTINF:-1,%s\n%s\n", title, streamURL)
	w.Write([]byte(playlist))
}

func (c *StreamController) handleDirectUpload(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden: administrator only", http.StatusForbidden)
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Invalid multipart request", http.StatusBadRequest)
		return
	}

	var sourceID int64
	var title string
	var savedFilename string
	var finalPath string

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed reading multipart stream", http.StatusBadRequest)
			return
		}

		formName := part.FormName()
		if formName == "source_id" {
			data, _ := io.ReadAll(part)
			sourceID, _ = strconv.ParseInt(string(data), 10, 64)
			continue
		}
		if formName == "title" {
			data, _ := io.ReadAll(part)
			title = strings.TrimSpace(string(data))
			continue
		}

		if formName == "file" {
			origName := part.FileName()
			cleanBase := filepath.Base(origName)
			ext := strings.ToLower(filepath.Ext(cleanBase))

			// Verify extension
			validExts := map[string]bool{
				".mp4": true, ".mkv": true, ".webm": true,
				".epub": true, ".pdf": true,
			}
			if !validExts[ext] {
				http.Error(w, "Unsupported file format", http.StatusBadRequest)
				return
			}

			source, err := c.sourceRepo.GetByID(r.Context(), sourceID)
			if err != nil || source == nil {
				http.Error(w, "Target storage source not found", http.StatusBadRequest)
				return
			}

			// Disk space check via syscall.Statfs
			var stat syscall.Statfs_t
			if err := syscall.Statfs(source.FolderPath, &stat); err == nil {
				freeSpace := uint64(stat.Bavail) * uint64(stat.Bsize)
				if freeSpace < 1024*1024*1024 { // 1 GB minimum
					http.Error(w, "Insufficient storage space on destination drive (< 1 GB)", http.StatusInsufficientStorage)
					return
				}
			}

			// Direct-to-disk streaming to <target>/<cleanBase>.part
			partPath := filepath.Join(source.FolderPath, cleanBase+".part")
			finalPath = filepath.Join(source.FolderPath, cleanBase)

			dst, err := os.Create(partPath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to stage file on disk: %v", err), http.StatusInternalServerError)
				return
			}

			_, copyErr := io.Copy(dst, part)
			dst.Close()

			if copyErr != nil || r.Context().Err() != nil {
				_ = os.Remove(partPath)
				http.Error(w, "Upload interrupted or aborted", http.StatusInternalServerError)
				return
			}

			// Atomically commit file
			if err := os.Rename(partPath, finalPath); err != nil {
				_ = os.Remove(partPath)
				http.Error(w, "Failed committing file to disk", http.StatusInternalServerError)
				return
			}

			savedFilename = cleanBase
			break
		}
	}

	if savedFilename == "" {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}

	if title == "" {
		title, _ = service.CleanTitleAndYear(savedFilename)
	}

	// Index directly into catalog
	source, err := c.sourceRepo.GetByID(r.Context(), sourceID)
	if err == nil && source != nil {
		stat, _ := os.Stat(finalPath)
		fileSize := int64(0)
		mtime := int64(0)
		if stat != nil {
			fileSize = stat.Size()
			mtime = stat.ModTime().Unix()
		}
		ext := strings.ToLower(filepath.Ext(savedFilename))
		format := strings.TrimPrefix(ext, ".")

		if source.MediaType == model.MediaTypeVideo {
			v := &model.Video{
				SourceID:       sourceID,
				Title:          title,
				VideoType:      model.VideoTypeMovie,
				FolderPath:     savedFilename,
				MetadataLocked: true,
			}
			vID, err := c.videoRepo.Create(r.Context(), v)
			if err == nil {
				vf := &model.VideoFile{
					VideoID:      vID,
					Title:        title,
					RelativePath: "",
					FileSize:     fileSize,
					MTime:        mtime,
					Format:       format,
				}
				_, _ = c.videoRepo.CreateFile(r.Context(), vf)
			}
		} else {
			b := &model.Book{
				SourceID:       sourceID,
				Title:          title,
				FolderPath:     savedFilename,
				MetadataLocked: true,
			}
			bID, err := c.bookRepo.Create(r.Context(), b)
			if err == nil {
				bf := &model.BookFile{
					BookID:       bID,
					Title:        title,
					RelativePath: "",
					FileSize:     fileSize,
					MTime:        mtime,
					Format:       model.BookFormat(format),
				}
				_, _ = c.bookRepo.CreateFile(r.Context(), bf)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"filename": savedFilename,
		"title":    title,
	})
}

func resolveDiskPath(sourceRoot, containerFolder, relativePath string) string {
	if relativePath == "" {
		return filepath.Join(sourceRoot, containerFolder)
	}
	return filepath.Join(sourceRoot, containerFolder, relativePath)
}

func validateCanonicalPath(sourceRoot, targetPath string) error {
	evalRoot, err := filepath.EvalSymlinks(sourceRoot)
	if err != nil {
		evalRoot = sourceRoot
	}
	evalTarget, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		evalTarget = targetPath
	}

	evalRoot = filepath.Clean(evalRoot)
	evalTarget = filepath.Clean(evalTarget)

	if !strings.HasPrefix(evalTarget, evalRoot) {
		return errors.New("path traversal detected")
	}
	return nil
}
