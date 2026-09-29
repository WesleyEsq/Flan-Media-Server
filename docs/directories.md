# Directory Structure & Package Guide

Flan Media Server is organized according to standard Go conventions and a clean **Model-View-Controller (MVC)** architectural pattern.

---

## 1. Project Directory Tree

```text
Flan-Media-Server/
├── cmd/
│   └── flan/               # Application entry point (main.go, CLI flags)
├── internal/               # Private application packages
│   ├── config/             # Strongly typed .env and environment variable loading
│   ├── database/           # SQLite connection, wear-leveling pragmas, and schema DDL
│   ├── model/              # [M] Domain models, sentinel errors, and SQL query stores
│   ├── controller/         # [C] HTTP handlers, status codes, and route registration
│   ├── middleware/         # HTTP filters (auth, rate limiting, stream governor)
│   ├── scraper/            # Local art finder and throttled TMDB background worker
│   └── storage/            # Disk space helper (statfs linux/other build tags)
├── web/                    # [V] Presentation layer embedded via embed.FS
│   ├── templates/          # 9 Server-rendered Go HTML templates
│   └── static/             # Assets: CSS, modular JS, mascots, and vendored player libs
├── docs/                   # System design, architecture, and threat model specifications
├── tests/                  # Hermetic unit and integration tests (fstest, mocks)
├── .env                    # Default environment configuration
├── LICENSE                 # Apache 2.0 License
└── README.md               # Quickstart and project overview
```

---

## 2. Package Responsibilities Matrix

| Package | MVC Layer | Key Files / Entities | Core Responsibility |
| :--- | :--- | :--- | :--- |
| **`cmd/flan`** | Entrypoint | `main.go` | Flag parsing (`--reset-admin`), config loading, DB bootstrap, dependency wiring, and graceful shutdown. |
| **`internal/config`** | Config | `config.go` | Reads `.env` and environment variables with safe defaults (fallback ports, directories, secrets). |
| **`internal/database`**| Persistence | `database.go` | Pure-Go SQLite connection (`modernc.org/sqlite`), WAL pragmas, `.flan-keep` check, pool limits (`SetMaxOpenConns(1)`), DDL migrations, and snapshots. |
| **`internal/model`** | **Model** | `user.go`, `library.go`, `movie.go`, `series.go`, `book.go`, `genre.go`, `progress.go`, `errors.go` | Domain structs, sentinel errors (`ErrNotFound`, `ErrDuplicate`), and pure SQL queries. Decoupled from HTTP transport. |
| **`internal/controller`**| **Controller** | `auth.go`, `user.go`, `library.go`, `movie.go`, `series.go`, `book.go`, `progress.go`, `upload.go`, `page.go` | Self-contained route registration (`RegisterRoutes(mux)`), HTTP decoding, ViewModel assembly, zero-copy `sendfile` streaming, and streaming uploads. |
| **`internal/middleware`**| Pipeline | `auth.go`, `ratelimit.go`, `governor.go` | Composable HTTP filters: HMAC session cookie verification, IP token buckets, brute-force delays, and the concurrent stream semaphore. |
| **`internal/scraper`** | Worker | `scanner.go`, `tmdb.go`, `cleaner.go` | In-place folder crawler, filename regex cleaning, local artwork detection, and throttled outbound TMDB queries (~2.8 req/s). |
| **`internal/storage`** | OS Abstraction | `disk_linux.go`, `disk_other.go` | Portable free-space verification via build tags (`unix.Statfs` on Linux; portable fallback for macOS/Windows). |
| **`web/templates`** | **View** | 9 `.html` files (`base.html`, `index.html`, etc.) | Server-rendered Go HTML templates. |
| **`web/static`** | Static Assets | `css/style.css`, `js/*.js`, `vendor/*` | Soft dark styling, vanilla client scripts, and embedded offline player dependencies (Plyr, ePub.js, JSZip). |

---

## 3. Separation of Concerns & Rules

1. **Decoupled Handlers:** Controllers never import `database/sql` directly; they accept narrow Go interfaces satisfied by models or mocks.
2. **Opaque Routes:** Public endpoints take integer IDs (`/stream/movie/42`), never host filesystem paths.
3. **Embedded Views:** All templates and static assets in `web/` are compiled directly into the binary via `embed.FS` with zero CDN calls.
4. **Hermetic Testing:** Packages are designed for testability using `testing/fstest.MapFS` and `net/http/httptest`.

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Database Schema & Wear-Leveling](database.md)
* [Storage Architecture](storage.md)
* [Testing Guidelines](testing.md)
