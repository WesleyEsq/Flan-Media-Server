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
│   ├── database/           # SQLite connection, wear-leveling pragmas, and 6-table schema
│   ├── model/              # [M] Domain models (video, book, user, progress)
│   ├── controller/         # [C] HTTP handlers, status codes, and route registration
│   ├── middleware/         # HTTP filters (auth, stream governor semaphore)
│   └── scraper/            # Simple local folder crawler (zero network dependencies)
├── web/                    # [V] Presentation layer embedded via embed.FS
│   ├── templates/          # 9 Server-rendered Go HTML templates (login, video, books, etc.)
│   └── static/             # Assets: CSS, modular JS, bundled sample avatars, player libs (Plyr)
├── data/                   # Persistent runtime directory (flan.db, data/avatars/)
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
| **`internal/config`** | Config | `config.go` | Reads `.env` and environment variables with safe defaults (ports, fixed directories, secrets). |
| **`internal/database`**| Persistence | `database.go` | Pure-Go SQLite connection (`modernc.org/sqlite`), WAL pragmas, `.flan-keep` check, pool limits (`SetMaxOpenConns(1)`), 6-table DDL migrations, and snapshots. |
| **`internal/model`** | **Model** | `user.go`, `video.go`, `book.go`, `progress.go`, `errors.go` | Domain structs, sentinel errors (`ErrNotFound`, `ErrDuplicate`), and pure SQL queries for containers and unified progress. |
| **`internal/controller`**| **Controller** | `auth.go`, `user.go`, `video.go`, `book.go`, `progress.go`, `avatar.go`, `page.go` | Self-contained route registration (`RegisterRoutes(mux)`), HTTP decoding, ViewModel assembly, zero-copy `sendfile` streaming, download routing, and custom avatar upload. |
| **`internal/middleware`**| Pipeline | `auth.go`, `governor.go` | Composable HTTP filters: HMAC session cookie verification and the concurrent stream semaphore. |
| **`internal/scraper`** | Local Scanner | `scanner.go` | 100% offline filesystem crawler: maps folder name to container title and discovers `poster.jpg`. |
| **`web/templates`** | **View** | 9 HTML templates | Sidebar-only layout, split Start screen, high-contrast cards with purple footers, and offline manual. |
| **`web/static`** | Static Assets | `css/style.css`, `js/*.js`, `vendor/plyr/*`, `assets/avatars/*` | Tactile styling with thick 2px black borders, bundled SVG preset avatars, and video player dependencies. |

---

## 3. Separation of Concerns & Rules

1. **Sidebar-Only Navigation:** Platform navigation is exclusively inside the left sidebar. The top header bar has zero navigation links.
2. **Decoupled Handlers:** Controllers never import `database/sql` directly; they accept narrow Go interfaces satisfied by models or mocks.
3. **Opaque Routes:** Public endpoints take integer IDs (`/stream/video/42`), never host filesystem paths.
4. **Embedded Views:** All templates, sample avatars, and assets in `web/` are compiled directly into the binary via `embed.FS` with zero CDN calls.
5. **Hermetic Testing:** Packages are designed for testability using `testing/fstest.MapFS` and `net/http/httptest`.

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Simplified Database Schema (6 Tables)](database.md)
* [Storage Architecture](storage.md)
* [Testing Guidelines](testing.md)
