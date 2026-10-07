package controller

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/middleware"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
)

type SystemController struct{}

func NewSystemController() *SystemController {
	return &SystemController{}
}

func (c *SystemController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/system/fs/browse", c.handleBrowseFS)
}

type FSEntry struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsDir     bool   `json:"is_dir"`
	HasMarker bool   `json:"has_marker"`
}

type FSMount struct {
	MountPoint string `json:"mount_point"`
	Device     string `json:"device"`
	FSType     string `json:"fs_type"`
}

type FSBrowseResponse struct {
	CurrentPath string    `json:"current_path"`
	ParentPath  string    `json:"parent_path"`
	Entries     []FSEntry `json:"entries"`
	Mounts      []FSMount `json:"mounts"`
}

func (c *SystemController) handleBrowseFS(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
		return
	}

	reqPath := r.URL.Query().Get("path")
	if reqPath == "" {
		// Default to /media if it exists, or /mnt, or root /
		if _, err := os.Stat("/media"); err == nil {
			reqPath = "/media"
		} else if _, err := os.Stat("/mnt"); err == nil {
			reqPath = "/mnt"
		} else {
			reqPath = "/"
		}
	}

	cleanPath := filepath.Clean(reqPath)
	info, err := os.Stat(cleanPath)
	if err != nil || !info.IsDir() {
		// Fallback to root
		cleanPath = "/"
	}

	// 1. Detect mounted storage volumes (parsing /proc/mounts on Linux)
	var mounts []FSMount
	if f, err := os.Open("/proc/mounts"); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		seen := make(map[string]bool)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 3 {
				device := fields[0]
				mountPoint := fields[1]
				fsType := fields[2]

				// Filter for real block devices or external storage mount locations
				isStorage := strings.HasPrefix(device, "/dev/") ||
					strings.HasPrefix(mountPoint, "/media") ||
					strings.HasPrefix(mountPoint, "/mnt") ||
					strings.HasPrefix(mountPoint, "/run/media")

				if isStorage && !strings.HasPrefix(device, "/dev/loop") && !seen[mountPoint] {
					seen[mountPoint] = true
					mounts = append(mounts, FSMount{
						MountPoint: mountPoint,
						Device:     device,
						FSType:     fsType,
					})
				}
			}
		}
	}

	// 2. Read subdirectories in current path
	var entries []FSEntry
	if dirEntries, err := os.ReadDir(cleanPath); err == nil {
		for _, de := range dirEntries {
			if !de.IsDir() {
				continue
			}
			name := de.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}

			fullSubPath := filepath.Join(cleanPath, name)
			hasMarker := database.VerifyMarker(fullSubPath) == nil

			entries = append(entries, FSEntry{
				Name:      name,
				Path:      fullSubPath,
				IsDir:     true,
				HasMarker: hasMarker,
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	parentPath := ""
	if cleanPath != "/" {
		parentPath = filepath.Dir(cleanPath)
	}

	resp := FSBrowseResponse{
		CurrentPath: cleanPath,
		ParentPath:  parentPath,
		Entries:     entries,
		Mounts:      mounts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
