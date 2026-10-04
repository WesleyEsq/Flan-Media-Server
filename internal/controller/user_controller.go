package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/WesleyEsq/Flan-Media-Server/internal/middleware"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
)

type UserController struct {
	userRepo    *repository.UserRepository
	authService *service.AuthService
	dataDir     string
}

func NewUserController(userRepo *repository.UserRepository, authService *service.AuthService, dataDir string) *UserController {
	return &UserController{
		userRepo:    userRepo,
		authService: authService,
		dataDir:     dataDir,
	}
}

func (c *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/users", c.handleListUsers)
	mux.HandleFunc("POST /api/users", c.handleCreateUser)
	mux.HandleFunc("PUT /api/users/self", c.handleUpdateSelf)
	mux.HandleFunc("DELETE /api/users/{id}", c.handleDeleteUser)
	mux.HandleFunc("POST /api/users/{id}/avatar", c.handleUploadAvatar)
}

func (c *UserController) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.userRepo.ListAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	var pubList []model.UserPublic
	for _, u := range users {
		pubList = append(pubList, u.ToPublic())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pubList)
}

func (c *UserController) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.UserFromContext(r.Context())
	if currentUser == nil || currentUser.Role != model.RoleAdmin {
		http.Error(w, "Forbidden: administrator privilege required", http.StatusForbidden)
		return
	}

	var req struct {
		Username   string     `json:"username"`
		PIN        string     `json:"pin"`
		Role       model.Role `json:"role"`
		AvatarIcon string     `json:"avatar_icon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.PIN = strings.TrimSpace(req.PIN)
	if req.Username == "" || req.PIN == "" {
		http.Error(w, "Username and PIN are required", http.StatusBadRequest)
		return
	}

	pinHash, err := c.authService.HashPIN(req.PIN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Role != model.RoleAdmin && req.Role != model.RoleUser {
		req.Role = model.RoleUser
	}
	if req.AvatarIcon == "" {
		req.AvatarIcon = "flan"
	}

	newUser := &model.User{
		Username:   req.Username,
		PINHash:    pinHash,
		Role:       req.Role,
		AvatarIcon: req.AvatarIcon,
	}

	id, err := c.userRepo.Create(r.Context(), newUser)
	if err != nil {
		http.Error(w, "Failed to create user (username may already exist)", http.StatusBadRequest)
		return
	}
	newUser.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    newUser.ToPublic(),
	})
}

func (c *UserController) handleUpdateSelf(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.UserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		DisplayName string `json:"display_name"`
		AvatarIcon  string `json:"avatar_icon"`
		PIN         string `json:"pin,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	currentUser.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.AvatarIcon != "" {
		currentUser.AvatarIcon = req.AvatarIcon
	}

	if err := c.userRepo.Update(r.Context(), currentUser); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	if req.PIN != "" {
		pinHash, err := c.authService.HashPIN(strings.TrimSpace(req.PIN))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := c.userRepo.UpdatePIN(r.Context(), currentUser.ID, pinHash); err != nil {
			http.Error(w, "Failed to update PIN", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    currentUser.ToPublic(),
	})
}

func (c *UserController) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.UserFromContext(r.Context())
	if currentUser == nil || currentUser.Role != model.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if targetID == currentUser.ID {
		http.Error(w, "Cannot delete your own account", http.StatusBadRequest)
		return
	}

	if err := c.userRepo.Delete(r.Context(), targetID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (c *UserController) handleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.UserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || (targetID != currentUser.ID && currentUser.Role != model.RoleAdmin) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Limit to 2 MB
	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
	if err := r.ParseMultipartForm(2 * 1024 * 1024); err != nil {
		http.Error(w, "File exceeds 2 MB limit", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Missing avatar file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Magic byte inspection for JPEG, PNG, WebP
	headerBytes := make([]byte, 512)
	n, _ := file.Read(headerBytes)
	contentType := http.DetectContentType(headerBytes[:n])

	var ext string
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	default:
		http.Error(w, "Only JPEG, PNG, and WebP images are allowed", http.StatusBadRequest)
		return
	}

	avatarsDir := filepath.Join(c.dataDir, "avatars")
	_ = os.MkdirAll(avatarsDir, 0755)

	targetFilename := fmt.Sprintf("%d%s", targetID, ext)
	targetPath := filepath.Join(avatarsDir, targetFilename)

	dst, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Write first read bytes and rest
	if _, err := dst.Write(headerBytes[:n]); err != nil {
		http.Error(w, "Failed writing file", http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed writing file", http.StatusInternalServerError)
		return
	}

	// Update DB record
	user, err := c.userRepo.GetByID(r.Context(), targetID)
	if err == nil {
		user.AvatarPath = targetPath
		_ = c.userRepo.Update(r.Context(), user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"filename": header.Filename,
	})
}
