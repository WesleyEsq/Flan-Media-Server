package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
)

type contextKey string

const userContextKey = contextKey("flan_user")

func Authenticate(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookieName := "flan_session"
			if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
				cookieName = "__Host-flan_session"
			}

			cookie, err := r.Cookie(cookieName)
			if err != nil && cookieName == "__Host-flan_session" {
				// Fallback to standard cookie name if running behind non-host-prefix setup
				cookie, err = r.Cookie("flan_session")
			}

			if err == nil && cookie != nil && cookie.Value != "" {
				user, needsRenewal, err := authService.ValidateSessionCookie(r.Context(), cookie.Value)
				if err == nil && user != nil {
					ctx := context.WithValue(r.Context(), userContextKey, user)
					r = r.WithContext(ctx)

					if needsRenewal {
						newVal, expires := authService.CreateSessionCookie(user)
						isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
						http.SetCookie(w, &http.Cookie{
							Name:     cookieName,
							Value:    newVal,
							Path:     "/",
							Expires:  expires,
							MaxAge:   30 * 24 * 3600,
							HttpOnly: true,
							Secure:   isSecure,
							SameSite: http.SameSiteLaxMode,
						})
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func UserFromContext(ctx context.Context) *model.User {
	u, _ := ctx.Value(userContextKey).(*model.User)
	return u
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user == nil {
			if strings.HasPrefix(r.URL.Path, "/api/") || strings.Contains(r.Header.Get("Accept"), "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user == nil || user.Role != model.RoleAdmin {
			if strings.HasPrefix(r.URL.Path, "/api/") || strings.Contains(r.Header.Get("Accept"), "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "Administrator privilege required"})
				return
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
