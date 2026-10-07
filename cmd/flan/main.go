package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/WesleyEsq/Flan-Media-Server/internal/config"
	"github.com/WesleyEsq/Flan-Media-Server/internal/controller"
	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/middleware"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
)

const version = "1.0.0"

func main() {
	resetAdminFlag := flag.Bool("reset-admin", false, "Interactively reset the administrator PIN")
	portFlag := flag.String("port", "", "Server HTTP port")
	versionFlag := flag.Bool("version", false, "Print server version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Flan Media Server v%s\n", version)
		os.Exit(0)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	if *portFlag != "" {
		cfg.Port = *portFlag
	}

	// 1. App Tier Mount Marker Check (./data/.flan-keep)
	dataMarker := filepath.Join(cfg.DataDir, ".flan-keep")
	if _, err := os.Stat(dataMarker); os.IsNotExist(err) {
		// If data directory doesn't exist yet, initialize it
		_ = os.MkdirAll(cfg.DataDir, 0755)
		if err := os.WriteFile(dataMarker, []byte("flan-keep"), 0644); err != nil {
			log.Fatalf("FATAL: Failed to initialize app tier marker at %s: %v", dataMarker, err)
		}
	}

	// 2. Connect Database
	db, err := database.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("FATAL: Database connection failed: %v", err)
	}
	defer db.Close()

	// 3. Construct Repositories
	userRepo := repository.NewUserRepository(db)
	sourceRepo := repository.NewSourceRepository(db)
	videoRepo := repository.NewVideoRepository(db)
	bookRepo := repository.NewBookRepository(db)
	progressRepo := repository.NewProgressRepository(db)

	// 4. Construct Services
	authService := service.NewAuthService(userRepo, db, cfg.SessionSecret)
	mediaService := service.NewMediaService(videoRepo, bookRepo, progressRepo, sourceRepo, cfg.DataDir)
	scannerService := service.NewScannerService(videoRepo, bookRepo, sourceRepo)

	// Handle --reset-admin CLI flag
	if *resetAdminFlag {
		handleResetAdmin(userRepo, authService)
		return
	}

	// 5. Seed default storage sources if empty
	seedDefaultStorage(sourceRepo, cfg.MediaDir)

	// 6. Check if first run (zero users)
	ctx := context.Background()
	userCount, err := userRepo.Count(ctx)
	if err == nil && userCount == 0 {
		token := service.GenerateBootstrapToken()
		authService.SetBootstrapToken(token)

		fmt.Println("==================================================================")
		fmt.Println("           FLAN MEDIA SERVER - FIRST BOOT INITIALIZATION           ")
		fmt.Println("==================================================================")
		fmt.Printf(" Welcome! No user profiles found. An initial setup is required.\n")
		fmt.Printf(" 1. Visit setup page:  http://localhost:%s/setup\n", cfg.Port)
		fmt.Printf(" 2. Enter this one-time Bootstrap Setup Token:  \n\n")
		fmt.Printf("               >>>   %s   <<<\n\n", token)
		fmt.Println(" This token protects your server from unauthorized first-run claims.")
		fmt.Println("==================================================================")
	}

	// 7. Construct Controllers
	authController := controller.NewAuthController(authService, userRepo, cfg.TrustedProxies)
	userController := controller.NewUserController(userRepo, authService, cfg.DataDir)
	mediaController := controller.NewMediaController(mediaService, scannerService, sourceRepo, videoRepo, bookRepo)
	streamController := controller.NewStreamController(videoRepo, bookRepo, sourceRepo, userRepo, authService, cfg.DataDir)
	systemController := controller.NewSystemController()
	pageController, err := controller.NewPageController(mediaService, authService, userRepo, sourceRepo, videoRepo, bookRepo)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize page controller: %v", err)
	}

	// 8. Register Routes
	mux := http.NewServeMux()
	authController.RegisterRoutes(mux)
	userController.RegisterRoutes(mux)
	mediaController.RegisterRoutes(mux)
	streamController.RegisterRoutes(mux)
	systemController.RegisterRoutes(mux)
	pageController.RegisterRoutes(mux)

	// 9. Attach Middleware Pipeline
	handler := middleware.Authenticate(authService)(mux)
	handler = middleware.CSRFProtect(handler)

	// 10. Start HTTP Server with Deadlines
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("Flan Media Server listening on http://0.0.0.0:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listener error: %v", err)
		}
	}()

	// 11. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Flan Media Server gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server exited cleanly.")
}

func handleResetAdmin(userRepo *repository.UserRepository, authService *service.AuthService) {
	ctx := context.Background()
	users, err := userRepo.ListAll(ctx)
	if err != nil || len(users) == 0 {
		fmt.Println("No user accounts found in database. Start the server and visit /setup.")
		os.Exit(1)
	}

	var admin *model.User
	for _, u := range users {
		if u.Role == model.RoleAdmin {
			admin = u
			break
		}
	}

	if admin == nil {
		fmt.Println("No administrator user found in database.")
		os.Exit(1)
	}

	fmt.Printf("Resetting PIN for administrator '%s' (User ID: %d)\n", admin.Username, admin.ID)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter new 4 to 6-digit numeric PIN: ")
	pin, _ := reader.ReadString('\n')
	pin = strings.TrimSpace(pin)

	pinHash, err := authService.HashPIN(pin)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := userRepo.UpdatePIN(ctx, admin.ID, pinHash); err != nil {
		fmt.Printf("Failed to update PIN: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Administrator PIN successfully updated. All prior sessions have been revoked.")
}

func seedDefaultStorage(sourceRepo *repository.SourceRepository, mediaDir string) {
	ctx := context.Background()
	sources, err := sourceRepo.ListAll(ctx)
	if err != nil || len(sources) > 0 {
		return
	}

	defaultVideoDir := filepath.Join(mediaDir, "video")
	defaultBooksDir := filepath.Join(mediaDir, "books")

	for _, dir := range []string{defaultVideoDir, defaultBooksDir} {
		_ = os.MkdirAll(dir, 0755)
		marker := filepath.Join(dir, ".flan-keep")
		if _, err := os.Stat(marker); os.IsNotExist(err) {
			_ = os.WriteFile(marker, []byte("flan-keep"), 0644)
		}
	}

	_, _ = sourceRepo.Create(ctx, &model.StorageSource{
		Name:       "Default Videos",
		MediaType:  model.MediaTypeVideo,
		FolderPath: defaultVideoDir,
		IsActive:   true,
	})

	_, _ = sourceRepo.Create(ctx, &model.StorageSource{
		Name:       "Default Books",
		MediaType:  model.MediaTypeBook,
		FolderPath: defaultBooksDir,
		IsActive:   true,
	})
}
