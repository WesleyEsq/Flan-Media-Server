package tests

import (
	"encoding/json"
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

func TestSystemFSBrowse(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "system_test.db")
	db, err := database.Connect(dbPath)
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, db, "system-secret-key-32-chars-long!")

	hash, _ := authSvc.HashPIN("1234")
	adminUser := createTestUser(t, userRepo, "sysadmin", hash, model.RoleAdmin, "mascot")
	normalUser := createTestUser(t, userRepo, "standarduser", hash, model.RoleUser, "mascot")

	// Create test folder structure
	browseTestDir := filepath.Join(tempDir, "browse_root")
	_ = os.MkdirAll(filepath.Join(browseTestDir, "Anime"), 0755)
	_ = os.MkdirAll(filepath.Join(browseTestDir, "Movies"), 0755)
	_ = os.WriteFile(filepath.Join(browseTestDir, "Anime", ".flan-keep"), []byte("flan-keep"), 0644)

	systemCtrl := controller.NewSystemController()
	mux := http.NewServeMux()
	systemCtrl.RegisterRoutes(mux)

	// Wrap with authentication middleware
	handler := middleware.Authenticate(authSvc)(mux)

	adminCookie, _ := authSvc.CreateSessionCookie(adminUser)
	memberCookie, _ := authSvc.CreateSessionCookie(normalUser)

	t.Run("Forbidden for non-admin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/system/fs/browse?path="+browseTestDir, nil)
		req.AddCookie(&http.Cookie{
			Name:  "flan_session",
			Value: memberCookie,
		})
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden, got %d", w.Code)
		}
	})

	t.Run("Admin can browse directory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/system/fs/browse?path="+browseTestDir, nil)
		req.AddCookie(&http.Cookie{
			Name:  "flan_session",
			Value: adminCookie,
		})
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d: %s", w.Code, w.Body.String())
		}

		var resp controller.FSBrowseResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if resp.CurrentPath != browseTestDir {
			t.Errorf("Expected current_path %s, got %s", browseTestDir, resp.CurrentPath)
		}

		if len(resp.Entries) != 2 {
			t.Fatalf("Expected 2 subdirectories (Anime, Movies), got %d", len(resp.Entries))
		}

		var foundAnime, foundMovies bool
		for _, e := range resp.Entries {
			if e.Name == "Anime" {
				foundAnime = true
				if !e.HasMarker {
					t.Errorf("Expected Anime to have .flan-keep marker")
				}
			}
			if e.Name == "Movies" {
				foundMovies = true
				if e.HasMarker {
					t.Errorf("Expected Movies to NOT have marker")
				}
			}
		}

		if !foundAnime || !foundMovies {
			t.Errorf("Entries did not contain expected directories: %+v", resp.Entries)
		}
	})
}
