# Directory Structure & Package Architecture

Flan Media Server is organized according to an enterprise-grade layered architecture: **Controller $\to$ Service $\to$ Repository $\to$ Model**. 

This structure provides an intuitive separation of concerns—clean domain entities, isolated SQL queries, dedicated business services, and thin HTTP handlers—while remaining idiomatic Go with zero external framework overhead.

---

## 1. Project Directory Tree

```text
Flan-Media-Server/
├── cmd/
│   └── flan/
│       └── main.go               # Entrypoint & DI container (wires config, DB, repos, services, controllers)
├── internal/                     # Private application packages
│   ├── config/                   # Configuration loading (.env, environment variables)
│   │   └── config.go
│   ├── database/                 # SQLite connection pooling, WAL pragmas, and schema migrations
│   │   ├── connection.go
│   │   └── migrations.go
│   ├── model/                    # Domain entities / records (zero DB or SQL dependencies)
│   │   ├── user.go
│   │   ├── video.go
│   │   ├── book.go
│   │   ├── progress.go
│   │   └── errors.go
│   ├── repository/               # Data Access Layer / DAOs (owns all raw SQL queries)
│   │   ├── user_repo.go
│   │   ├── video_repo.go
│   │   ├── book_repo.go
│   │   └── progress_repo.go
│   ├── service/                  # Business Logic Layer (auth, catalog queries, scanner, rate limits)
│   │   ├── auth_service.go       # PIN validation, lockout tracking, session HMAC, signed VLC URLs
│   │   ├── media_service.go      # Media container retrieval, detail view assembly, cover art resolution
│   │   ├── scanner_service.go    # Offline filesystem crawler & database reconciliation engine
│   │   └── governor_service.go   # 3-stream concurrent playback lease governor
│   ├── controller/               # HTTP Presentation Layer (request decoding, status codes, view rendering)
│   │   ├── auth_controller.go
│   │   ├── media_controller.go
│   │   ├── stream_controller.go  # Zero-copy kernel sendfile streaming (HTTP 206 Range)
│   │   ├── user_controller.go
│   │   └── page_controller.go    # Server-rendered HTML page endpoints
│   └── middleware/               # HTTP Request Pipeline / Filters
│       ├── auth_filter.go        # Session cookie verification & role checks
│       ├── csrf_filter.go        # Cross-origin protection & origin validation
│       └── governor_filter.go    # Playback lease validation for /stream/ and /download/
├── web/                          # Embedded presentation assets (templates and static files)
│   ├── embed.go                  # Package web: exports embedded filesystem (embed.FS)
│   ├── templates/                # 9 Server-rendered Go HTML templates (login, video, books, manage, etc.)
│   └── static/                   # CSS, modular JS, bundled SVG avatars, Plyr video player
├── data/                         # Persistent runtime directory (flan.db, data/covers/, data/avatars/)
├── media/                        # Fixed media storage directories (./media/video, ./media/books)
├── docs/                         # Architecture, database schema, security, and operations documentation
├── tests/                        # Unit and integration tests (fstest, mocks)
├── .env                          # Default environment configuration
├── LICENSE                       # Apache 2.0 License
└── README.md                     # Quickstart and project overview
```

---

## 2. Package Responsibilities

* **`cmd/flan` (Application Entrypoint):**
  * Parses command-line flags (`--reset-admin`).
  * Loads environment configuration.
  * Connects and initializes the database.
  * Explicitly constructs and injects dependencies.
  * Starts the HTTP listener and handles graceful shutdown.

* **`internal/config` (Configuration Properties):**
  * Loads `.env` file and environment variables into a strongly typed `Config` struct with fallback defaults.

* **`internal/database` (Persistence Infrastructure):**
  * Establishes SQLite connection pools (dedicated writer handle and parallel reader handles).
  * Executes WAL pragmas and verifies directory markers (`.flan-keep`).
  * Applies schema migrations using SQLite's `PRAGMA user_version`.

* **`internal/model` (Domain Entities):**
  * Defines pure state structs (`Video`, `VideoFile`, `Book`, `BookFile`, `User`, `Progress`).
  * Declares sentinel domain errors (`ErrNotFound`, `ErrDuplicate`).
  * Strictly forbidden from importing `database/sql` or performing I/O.

* **`internal/repository` (Data Access Objects / DAOs):**
  * Contains all raw SQL queries (`SELECT`, `INSERT`, `UPDATE`, `DELETE`).
  * Maps SQL result rows into `model` structs or slices.
  * Isolates database-specific logic from the rest of the application.

* **`internal/service` (Business Logic & Workflows):**
  * Implements domain rules: password hashing and verification, lockout backoff, 3-stream lease management, and filesystem crawler reconciliation.
  * Coordinates multiple repositories when executing business transactions.

* **`internal/controller` (HTTP Presentation Layer):**
  * Defines route handlers registered onto `http.ServeMux`.
  * Decodes URL parameters, query strings, and JSON request bodies.
  * Invokes the relevant service and serializes JSON responses or renders HTML templates.
  * Manages zero-copy byte-range video delivery via `http.ServeContent`.

* **`internal/middleware` (HTTP Request Interceptors):**
  * Chained filters wrapping `http.Handler`: session authentication, CSRF origin verification, and stream capacity checks.

* **`web` (Static Asset Bundle):**
  * Uses Go's `embed.FS` to bundle `templates/` and `static/` directly into the binary.

---

## 3. Java-to-Go Architectural Mapping

For engineers familiar with Java enterprise frameworks (such as Spring Boot, Quarkus, or Micronaut), here is how concepts translate:

* **Domain Records / Entities:**
  * Java: `public record Video(long id, String title) {}`
  * Go: `type Video struct { ID int64; Title string }` in `internal/model`

* **Data Access Objects (DAOs):**
  * Java: `@Repository public class VideoRepository { ... }`
  * Go: `type VideoRepository struct { db *sql.DB }` in `internal/repository`

* **Business Services:**
  * Java: `@Service public class VideoService { ... }`
  * Go: `type VideoService struct { repo VideoRepository }` in `internal/service`

* **HTTP Controllers:**
  * Java: `@RestController public class VideoController { ... }`
  * Go: `type VideoController struct { service VideoService }` in `internal/controller`

* **Request Filters:**
  * Java: `HandlerInterceptor` or `Filter`
  * Go: `func(http.Handler) http.Handler` in `internal/middleware`

* **Configuration Management:**
  * Java: `@ConfigurationProperties`
  * Go: `internal/config/config.go`

* **Schema Migrations:**
  * Java: Flyway / Liquibase
  * Go: `internal/database/migrations.go` using SQLite's `PRAGMA user_version`

### Constructor Dependency Injection

Dependencies are wired explicitly using constructor factory functions (`New...`) in `cmd/flan/main.go`:

```go
// 1. Database
db := database.Connect(cfg.DBPath)

// 2. Repositories (Data Access)
videoRepo := repository.NewVideoRepository(db)
userRepo := repository.NewUserRepository(db)

// 3. Services (Business Logic)
authService := service.NewAuthService(userRepo, cfg.SessionSecret)
mediaService := service.NewMediaService(videoRepo)

// 4. Controllers (HTTP Handlers)
videoController := controller.NewVideoController(mediaService)
authController := controller.NewAuthController(authService)

// 5. Route Registration
mux := http.NewServeMux()
videoController.RegisterRoutes(mux)
authController.RegisterRoutes(mux)
```

### Interface-Based Decoupling

Layers can be decoupled using Go interfaces to facilitate unit testing:

```go
package service

type VideoRepository interface {
    GetByID(ctx context.Context, id int64) (*model.Video, error)
    ListAll(ctx context.Context, query string) ([]model.Video, error)
}

type MediaService struct {
    repo VideoRepository
}

func NewMediaService(repo VideoRepository) *MediaService {
    return &MediaService{repo: repo}
}
```

This pattern allows substituting mock implementations during unit tests, equivalent to Mockito in Java.

---

## 4. Architectural Invariants

1. **Clean Model Layer:** Structs in `internal/model` must never import `database/sql` or execute database queries.
2. **Repository Isolation:** Only `internal/repository` and `internal/database` may import `database/sql`. Controllers and services never execute raw SQL.
3. **Thin Controllers:** Handlers in `internal/controller` only parse HTTP inputs, invoke the appropriate service, and serialize the response or render the view.
4. **Sidebar Navigation:** Platform navigation is exclusively inside the left sidebar (bottom bar on mobile). The top header contains no navigation links.
5. **Opaque Routes:** Public media endpoints take integer IDs (`/stream/video/42`), never host filesystem paths.
6. **Embedded Views:** All templates and static assets in `web/` are compiled directly into the binary via `embed.FS`.

---

## 5. Related Documentation

* [Master System Architecture](design.md)
* [Database Schema & Queries](database.md)
* [Storage Architecture](storage.md)
* [Testing Strategy](testing.md)
