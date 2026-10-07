package controller

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/WesleyEsq/Flan-Media-Server/internal/middleware"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
	"github.com/WesleyEsq/Flan-Media-Server/web"
)

type PageController struct {
	mediaService *service.MediaService
	authService  *service.AuthService
	userRepo     *repository.UserRepository
	sourceRepo   *repository.SourceRepository
	videoRepo    *repository.VideoRepository
	bookRepo     *repository.BookRepository
	templates    map[string]*template.Template
}

func NewPageController(
	mediaService *service.MediaService,
	authService *service.AuthService,
	userRepo *repository.UserRepository,
	sourceRepo *repository.SourceRepository,
	videoRepo *repository.VideoRepository,
	bookRepo *repository.BookRepository,
) (*PageController, error) {
	c := &PageController{
		mediaService: mediaService,
		authService:  authService,
		userRepo:     userRepo,
		sourceRepo:   sourceRepo,
		videoRepo:    videoRepo,
		bookRepo:     bookRepo,
		templates:    make(map[string]*template.Template),
	}

	funcMap := template.FuncMap{
		"slice": func(args ...string) []string {
			return args
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"humanizeBytes": func(bytes int64) string {
			const unit = 1024
			if bytes < unit {
				return fmt.Sprintf("%d B", bytes)
			}
			div, exp := int64(unit), 0
			for n := bytes / unit; n >= unit; n /= unit {
				div *= unit
				exp++
			}
			return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
		},
	}

	// Parse all page templates paired with layout.html
	pageNames := []string{
		"login.html", "setup.html", "video.html", "video_detail.html",
		"books.html", "book_detail.html", "search.html", "watch.html",
		"read.html", "manage.html", "manual.html",
	}

	for _, name := range pageNames {
		tmpl, err := template.New("base").Funcs(funcMap).ParseFS(web.TemplatesFS, "templates/layout.html", "templates/"+name)
		if err != nil {
			return nil, fmt.Errorf("parse template %s: %w", name, err)
		}
		c.templates[name] = tmpl
	}

	return c, nil
}

func (c *PageController) RegisterRoutes(mux *http.ServeMux) {
	// Static embedded assets
	staticSub, _ := fs.Sub(web.StaticFS, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	// Health check
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Public HTML views
	mux.HandleFunc("GET /login", c.handleLoginPage)
	mux.HandleFunc("GET /setup", c.handleSetupPage)

	// Authenticated HTML views
	mux.Handle("GET /{$}", middleware.RequireAuth(http.HandlerFunc(c.handleRoot)))
	mux.Handle("GET /video", middleware.RequireAuth(http.HandlerFunc(c.handleVideoCatalog)))
	mux.Handle("GET /video/{id}", middleware.RequireAuth(http.HandlerFunc(c.handleVideoDetail)))
	mux.Handle("GET /books", middleware.RequireAuth(http.HandlerFunc(c.handleBooksCatalog)))
	mux.Handle("GET /books/{id}", middleware.RequireAuth(http.HandlerFunc(c.handleBookDetail)))
	mux.Handle("GET /search", middleware.RequireAuth(http.HandlerFunc(c.handleSearchPage)))
	mux.Handle("GET /watch/{file_id}", middleware.RequireAuth(http.HandlerFunc(c.handleWatchPage)))
	mux.Handle("GET /read/{file_id}", middleware.RequireAuth(http.HandlerFunc(c.handleReadPage)))
	mux.Handle("GET /manual", middleware.RequireAuth(http.HandlerFunc(c.handleManualPage)))

	// Administrator HTML view
	mux.Handle("GET /manage", middleware.RequireAuth(middleware.RequireAdmin(http.HandlerFunc(c.handleManagePage))))
}

func (c *PageController) handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/video", http.StatusSeeOther)
}

func (c *PageController) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	// If first run with zero users, redirect to setup
	count, _ := c.userRepo.Count(r.Context())
	if count == 0 {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	users, _ := c.userRepo.ListAll(r.Context())
	data := map[string]interface{}{
		"PageTitle": "Login",
		"Users":     users,
	}
	c.render(w, "login.html", data)
}

func (c *PageController) handleSetupPage(w http.ResponseWriter, r *http.Request) {
	count, _ := c.userRepo.Count(r.Context())
	if count > 0 && c.authService.GetBootstrapToken() == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"PageTitle": "Initial Setup",
	}
	c.render(w, "setup.html", data)
}

func (c *PageController) handleVideoCatalog(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	filterType := r.URL.Query().Get("type")
	if filterType == "" {
		filterType = "all"
	}

	videos, _ := c.mediaService.ListVideos(r.Context(), filterType, false)
	continueWatching, _ := c.mediaService.GetContinueWatching(r.Context(), user.ID)

	data := map[string]interface{}{
		"PageTitle":        "Videos & Movies",
		"ActiveTab":        "video",
		"CurrentUser":      user,
		"Videos":           videos,
		"ContinueWatching": continueWatching,
		"SelectedType":     filterType,
	}
	c.render(w, "video.html", data)
}

func (c *PageController) handleVideoDetail(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	video, files, err := c.mediaService.GetVideo(r.Context(), id, user.ID, false)
	if err != nil || video == nil {
		http.NotFound(w, r)
		return
	}

	signedURLFunc := func(fileID int64) string {
		return c.authService.GenerateSignedURL(model.MediaTypeVideo, fileID, user.ID, user.TokenVersion)
	}

	data := map[string]interface{}{
		"PageTitle":     video.Title,
		"ActiveTab":     "video",
		"CurrentUser":   user,
		"Video":         video,
		"Files":         files,
		"SignedURLFunc": signedURLFunc,
	}
	c.render(w, "video_detail.html", data)
}

func (c *PageController) handleBooksCatalog(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	formatFilter := r.URL.Query().Get("format")
	if formatFilter == "" {
		formatFilter = "all"
	}

	books, _ := c.mediaService.ListBooks(r.Context(), formatFilter, false)
	jumpBackIn, _ := c.mediaService.GetJumpBackInBooks(r.Context(), user.ID)

	data := map[string]interface{}{
		"PageTitle":      "Books & Publications",
		"ActiveTab":      "books",
		"CurrentUser":    user,
		"Books":          books,
		"JumpBackIn":     jumpBackIn,
		"SelectedFormat": formatFilter,
	}
	c.render(w, "books.html", data)
}

func (c *PageController) handleBookDetail(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	book, files, err := c.mediaService.GetBook(r.Context(), id, user.ID, false)
	if err != nil || book == nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"PageTitle":   book.Title,
		"ActiveTab":   "books",
		"CurrentUser": user,
		"Book":        book,
		"Files":       files,
	}
	c.render(w, "book_detail.html", data)
}

func (c *PageController) handleSearchPage(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	data := map[string]interface{}{
		"PageTitle":   "Search Library",
		"ActiveTab":   "search",
		"CurrentUser": user,
	}
	c.render(w, "search.html", data)
}

func (c *PageController) handleWatchPage(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	fileIDStr := r.PathValue("file_id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	file, video, _, err := c.videoRepo.GetFileWithContainer(r.Context(), fileID)
	if err != nil || file == nil || video == nil {
		http.NotFound(w, r)
		return
	}

	// Fetch current progress
	var posSec float64
	var formattedPos string
	prog, err := c.mediaService.GetProgress(r.Context(), user.ID, model.MediaTypeVideo, fileID)
	if err == nil && prog != nil {
		if s, err := strconv.ParseFloat(prog.PositionData, 64); err == nil {
			posSec = s
			m := int(s) / 60
			sec := int(s) % 60
			formattedPos = fmt.Sprintf("%02d:%02d", m, sec)
		}
	}

	signedURL := c.authService.GenerateSignedURL(model.MediaTypeVideo, fileID, user.ID, user.TokenVersion)

	files, _ := c.videoRepo.ListFilesByVideoID(r.Context(), video.ID, user.ID, false)
	var prevFile, nextFile *model.VideoFile
	for i, f := range files {
		if f.ID == fileID {
			if i > 0 {
				prevFile = files[i-1]
			}
			if i < len(files)-1 {
				nextFile = files[i+1]
			}
			break
		}
	}

	data := map[string]interface{}{
		"PageTitle":         file.DisplayTitle(),
		"ActiveTab":         "video",
		"CurrentUser":       user,
		"Video":             video,
		"File":              file,
		"PrevFile":          prevFile,
		"NextFile":          nextFile,
		"PositionSeconds":   posSec,
		"FormattedPos":      formattedPos,
		"SignedDownloadURL": signedURL,
	}
	c.render(w, "watch.html", data)
}

func (c *PageController) handleReadPage(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	fileIDStr := r.PathValue("file_id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	file, book, _, err := c.bookRepo.GetFileWithContainer(r.Context(), fileID)
	if err != nil || file == nil || book == nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"PageTitle":   file.DisplayTitle(),
		"ActiveTab":   "books",
		"CurrentUser": user,
		"Book":        book,
		"File":        file,
	}
	c.render(w, "read.html", data)
}

func (c *PageController) handleManagePage(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	sources, _ := c.sourceRepo.ListAll(r.Context())
	users, _ := c.userRepo.ListAll(r.Context())

	data := map[string]interface{}{
		"PageTitle":   "Manage Server",
		"ActiveTab":   "manage",
		"CurrentUser": user,
		"Sources":     sources,
		"Users":       users,
	}
	c.render(w, "manage.html", data)
}

func (c *PageController) handleManualPage(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	data := map[string]interface{}{
		"PageTitle":   "Software Manual",
		"ActiveTab":   "manual",
		"CurrentUser": user,
	}
	c.render(w, "manual.html", data)
}

func (c *PageController) render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := c.templates[name]
	if !ok {
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.ExecuteTemplate(w, "base", data)
}
