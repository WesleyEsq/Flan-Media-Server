package controller

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
)

type AuthController struct {
	authService    *service.AuthService
	userRepo       *repository.UserRepository
	trustedProxies []string
}

func NewAuthController(authService *service.AuthService, userRepo *repository.UserRepository, trustedProxies []string) *AuthController {
	return &AuthController{
		authService:    authService,
		userRepo:       userRepo,
		trustedProxies: trustedProxies,
	}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/setup", c.handleSetup)
	mux.HandleFunc("POST /api/login", c.handleLogin)
	mux.HandleFunc("POST /api/logout", c.handleLogout)
}

func (c *AuthController) handleSetup(w http.ResponseWriter, r *http.Request) {
	// Only allow setup if there are zero users or bootstrap token is set
	count, err := c.userRepo.Count(r.Context())
	if err != nil || (count > 0 && c.authService.GetBootstrapToken() == "") {
		http.Error(w, "Setup already completed", http.StatusForbidden)
		return
	}

	var req struct {
		Token    string `json:"token"`
		Username string `json:"username"`
		PIN      string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.Token = strings.TrimSpace(req.Token)
	req.Username = strings.TrimSpace(req.Username)
	req.PIN = strings.TrimSpace(req.PIN)

	if !c.authService.VerifyBootstrapToken(req.Token) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bootstrap setup token"})
		return
	}

	pinHash, err := c.authService.HashPIN(req.PIN)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	adminUser := &model.User{
		Username:   req.Username,
		PINHash:    pinHash,
		Role:       model.RoleAdmin,
		AvatarIcon: "flan",
	}

	userID, err := c.userRepo.Create(r.Context(), adminUser)
	if err != nil {
		http.Error(w, "Failed to create administrator account", http.StatusInternalServerError)
		return
	}
	adminUser.ID = userID

	// Clear bootstrap token
	c.authService.ClearBootstrapToken()

	// Issue session cookie
	c.setSessionCookie(w, r, adminUser)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    adminUser.ToPublic(),
	})
}

func (c *AuthController) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64  `json:"user_id"`
		PIN    string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	clientIP := c.extractClientIP(r)
	user, err := c.authService.AuthenticatePIN(r.Context(), clientIP, req.UserID, req.PIN)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(err.Error(), "locked") {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid PIN"})
		return
	}

	c.setSessionCookie(w, r, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user.ToPublic(),
	})
}

func (c *AuthController) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookieName := "flan_session"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		cookieName = "__Host-flan_session"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (c *AuthController) setSessionCookie(w http.ResponseWriter, r *http.Request, user *model.User) {
	val, expires := c.authService.CreateSessionCookie(user)
	isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	cookieName := "flan_session"
	if isSecure {
		cookieName = "__Host-flan_session"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    val,
		Path:     "/",
		Expires:  expires,
		MaxAge:   30 * 24 * 3600,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (c *AuthController) extractClientIP(r *http.Request) string {
	remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteHost = r.RemoteAddr
	}

	if len(c.trustedProxies) > 0 {
		for _, trusted := range c.trustedProxies {
			if remoteHost == trusted {
				if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
					parts := strings.Split(xff, ",")
					client := strings.TrimSpace(parts[0])
					if client != "" {
						return client
					}
				}
				break
			}
		}
	}
	return remoteHost
}
