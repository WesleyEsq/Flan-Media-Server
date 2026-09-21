# Directories in the Project

The project is organized according to standard Go conventions and a clean **Model-View-Controller (MVC)** architectural pattern. This keeps domain data models, presentation templates, HTTP controllers, cross-cutting middlewares, and database initialization cleanly separated.

## Directory Structure

```text
Flan-Media-Server/
├── cmd/
│   └── flan/               # Server entry point (main.go)
├── internal/               # Private application packages
│   ├── config/             # .env and environment variable loading
│   ├── database/           # SQLite connection, pragmas, pool limits, and schema DDL
│   │   └── database.go     # Dedicated DB initialization and configuration file
│   ├── model/              # [M] Model Layer: Domain entities, stores, and queries
│   │   ├── errors.go       # Sentinel domain errors (ErrNotFound, ErrDuplicate, etc.)
│   │   ├── user.go         # User entity, lockout tracking, and PIN verification
│   │   ├── library.go      # Library entity, mount validation, and queries
│   │   ├── movie.go        # Movie entity, catalog queries, and file tracking
│   │   ├── series.go       # Series & Episode entities, season tree queries
│   │   ├── book.go         # Book entity, format checks, and catalog queries
│   │   ├── genre.go        # Genre & ItemGenre normalized mapping
│   │   └── progress.go     # VideoProgress & BookProgress upsert/query stores
│   ├── controller/         # [C] Controller Layer: HTTP handlers & self-contained routes
│   │   ├── auth.go         # /login, /setup, /api/auth routes and session cookies
│   │   ├── user.go         # /api/users CRUD endpoints and PIN reset routes
│   │   ├── library.go      # /api/libraries management and scan trigger routes
│   │   ├── movie.go        # /api/movies, /stream/movie/{id}, and /covers/movie/{id}
│   │   ├── series.go       # /api/series, /stream/episode/{id}, and series covers
│   │   ├── book.go         # /api/books, /stream/book/{id}, and document delivery
│   │   ├── progress.go     # /api/progress watch and reading position sync
│   │   ├── upload.go       # /api/upload streaming multipart file/series intake
│   │   └── page.go         # HTML template page views and ViewModel assembly
│   ├── middleware/         # Cross-cutting HTTP filters
│   │   ├── auth.go         # HMAC cookie authentication & user context injection
│   │   ├── ratelimit.go    # IP token bucket & progressive lockout throttles
│   │   └── governor.go     # In-memory semaphore concurrent stream governor
│   └── scraper/            # Local-first artwork finder & throttled TMDB worker
├── web/                    # [V] View Layer: Embedded web assets (via embed.FS)
│   ├── templates/          # Server-rendered Go HTML templates
│   │   ├── base.html       # Shell, navbar with e-manual dialog, and footer
│   │   ├── index.html      # Dashboard with Continue Watching/Reading shelves
│   │   ├── videos.html     # Movies and TV series catalog with genre pills
│   │   ├── show.html       # TV series detail view with season tabs and episodes
│   │   ├── books.html      # Book and document catalog view
│   │   ├── watch.html      # Plyr video player view with auto-resume prompt
│   │   ├── read.html       # Native PDF iframe embed & ePub.js book reader
│   │   ├── login.html      # "Who is watching?" profile selector & PIN keypad
│   │   ├── setup.html      # Initial onboarding wizard (admin & first library)
│   │   └── settings.html   # Storage health, library manager, and user settings
│   └── static/             # Static web files served under /static/
│       ├── css/            # Soft dark slate stylesheet with lavender accents (style.css)
│       ├── js/             # Modular vanilla JS (api.js, player.js, reader.js)
│       ├── vendor/         # Offline client dependencies (Plyr, ePub.js, JSZip)
│       └── assets/         # Curated profile avatar icons and mascots
├── docs/                   # Architectural specifications and diagrams
├── tests/                  # Unit and integration tests (TDD virtual fs / mocks)
├── .env                    # Default environment configuration
├── LICENSE                 # Apache 2.0 License
└── README.md               # Project overview and instructions
```

---

## Package Responsibilities

### `cmd/flan`

Contains `main.go`, the executable entry point:

- Parses CLI flags (including `--reset-admin` failsafe).
- Loads configuration via `internal/config`.
- Initializes the SQLite database via `internal/database`.
- Instantiates domain models (`internal/model`), middlewares (`internal/middleware`), and controllers (`internal/controller`).
- Mounts controller routes directly to Go 1.22's standard library multiplexer (`http.ServeMux`).
- Handles graceful shutdown on SIGINT/SIGTERM.

### `internal/config`

Reads `.env` and system environment variables into strongly-typed configuration structs with default fallbacks for ports, hosts, and memory targets.

### `internal/database` (Database Configuration & Initialization)

Houses the dedicated database initialization logic (`database.go`):

- Opens the SQLite database connection using pure-Go `modernc.org/sqlite`.
- Enforces single-connection pool limits (`db.SetMaxOpenConns(1)`, `db.SetMaxIdleConns(1)`) to guarantee minimal memory usage and zero write-lock contention.
- Applies low-memory wear-leveling pragmas (`WAL`, `synchronous = NORMAL`, `cache_size = -2000`, `busy_timeout = 5000`).
- Validates the `.flan-keep` marker file to prevent ghost databases on unmounted drives.
- Executes automatic DDL schema migrations on boot (creating `users`, `libraries`, `movies`, `series`, `episodes`, `books`, `genres`, `item_genres`, `video_progress`, and `book_progress`).
- Provides integrity checks (`PRAGMA quick_check;`) and daily zero-lock hot backups (`VACUUM INTO`).

### `internal/model` (The Model Layer)

Contains domain entities and data access stores:

- **`errors.go`**: Defines sentinel domain errors (`ErrNotFound`, `ErrDuplicate`, `ErrUnauthorized`, `ErrDiskLow`), decoupling transport from persistence.
- **`user.go`**: Handles user authentication queries, bcrypt PIN verification, and progressive lockout counters.
- **`library.go`**: Manages library paths, mount validation, and directory queries.
- **`movie.go`**: Handles movie catalog queries, filtering by genre, and relative path lookups.
- **`series.go`**: Manages TV series records, season groupings, and episode queries.
- **`book.go`**: Handles book catalog filtering by format (`epub`, `pdf`) and genre.
- **`genre.go`**: Manages the normalized `genres` and `item_genres` junction tables.
- **`progress.go`**: Manages atomic upserts for `video_progress` (seconds) and `book_progress` (EPUB CFI / PDF page counts).

### `internal/controller` (The Controller Layer)

Encapsulates HTTP transport and business orchestration:

- Each controller defines a `RegisterRoutes(mux *http.ServeMux)` method that binds its own endpoints.
- Decodes HTTP query parameters, path variables, and JSON/multipart bodies.
- Translates model sentinel errors into standard HTTP status codes (404, 409, 507).
- Assembles ViewModels and renders server-side Go HTML templates.
- Manages zero-copy streaming (`http.ServeContent` / `sendfile`) and streaming multipart uploads (`r.MultipartReader`).

### `internal/middleware` (HTTP Filters)

Provides composable standard library HTTP middleware:

- **`auth.go`**: Validates HMAC-signed session cookies and injects the authenticated `userID` and `role` into request contexts.
- **`ratelimit.go`**: Enforces IP token buckets (Zone C) and brute-force lockout intervals (Zone B).
- **`governor.go`**: Enforces the counting semaphore limiting concurrent video streams to protect mechanical drives (Zone A).

### `internal/scraper`

Operates as a background worker:

- Inspects local folders for existing artwork (`poster.jpg`, cover art, embedded EPUB art).
- Sanitizes filenames using regex patterns to extract canonical titles, years, seasons, and episodes.
- Queries TMDB at a metered rate (~2.8 req/s) and streams cover images directly to disk.

### `web/templates` (The View Layer)

Server-rendered Go HTML templates (`html/template`) styled with a soft dark slate background and lavender accents. Embedded into the binary via `embed.FS`.

### `web/static`

Static stylesheets (`style.css`), modular vanilla JavaScript (`api.js`, `player.js`, `reader.js`), vendored offline player libraries (Plyr, ePub.js, JSZip), and mascot SVGs. Embedded into the binary via `embed.FS`.

---

### Related Documentation

- [Master System Specifications](design.md)
- [Database Schema & Wear-Leveling Pragmas](database.md)
- [Storage Architecture & Drive Resiliency](storage.md)
- [Testing Strategy & TDD Guidelines](testing.md)
