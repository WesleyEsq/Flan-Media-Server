# Master System Architecture

High-level architectural design and core technical decisions for Flan Media Server.

---

## 1. System Constraints & Operating Targets

Flan is built specifically for underpowered Linux devices (Raspberry Pi Zero, 1–5, low-spec ARM/x86 boards) serving 1 to 5 household users.

| Parameter | Specification | Engineering Rationale |
| :--- | :--- | :--- |
| **Memory Budget** | 15–20 MB RAM | Soft ceiling via `GOMEMLIMIT=16MiB` and aggressive GC via `GOGC=30`. |
| **CPU Budget** | Low-power ARM cores | Real-time transcoding is omitted; all media streams directly. |
| **Concurrent Streams** | Max 3 active | Enforced by counting semaphore to prevent USB mechanical drive thrashing. |
| **Connections** | 10–15 idle connections | Lightweight Go goroutines (~2 KB per connection; <50 KB overhead). |
| **Deployment** | Single static binary | Pure-Go driver (`modernc.org/sqlite`, `CGO_ENABLED=0`), assets embedded via `embed.FS`. |

---

## 2. Software Architecture (MVC)

The codebase implements a standard Model-View-Controller pattern using Go 1.22+ standard library features:

```text
HTTP Request → [Middleware Pipeline] → Controller (RegisterRoutes) → Model (Store) → SQLite WAL
                                      ↓
                               ViewModel Assembly → View (html/template) → Client Response
```

* **Model Layer (`internal/model`):** Domain structs, business validation, sentinel errors (`ErrNotFound`, `ErrDuplicate`), and pure SQL queries.
* **View Layer (`web/templates`, `web/static`):** 9 server-rendered Go HTML templates and vanilla CSS/JS packaged directly into the binary via `embed.FS`.
* **Controller Layer (`internal/controller`):** HTTP transport adapters, request decoders, status code mapping, ViewModel assembly, and route registration onto `http.ServeMux`.
* **Middleware Layer (`internal/middleware`):** Cross-cutting concerns: HMAC cookie authentication, 5-zone rate limiting, and the stream governor semaphore.

---

## 3. Key Technical Decisions

### Direct Play & Kernel Zero-Copy Streaming

* **No Transcoding:** Videos must be pre-encoded in web-compatible formats (MP4 with H.264/AAC, WebM with VP9/Opus, or web-safe MKV). Non-web containers trigger a fallback warning offering direct file download.
* **Linux `sendfile`:** Video delivery delegates to Go's standard library `http.ServeContent` with HTTP 206 Partial Content (Range requests). Media bytes move straight from kernel page cache to the network socket without traversing Go user-space memory.

### SQLite with Micro-SD Wear-Leveling

* **Driver:** Pure-Go SQLite (`modernc.org/sqlite`) allows instant cross-compilation without C toolchains.
* **Single Connection Pool:** `db.SetMaxOpenConns(1)` eliminates write-lock contention and conserves RAM. Background scans yield cooperatively (`runtime.Gosched()`) to keep UI requests responsive.
* **Pragmas:** WAL journal mode, `synchronous = NORMAL` (preserves micro-SD lifespan), `cache_size = -2000` (~2 MB RAM cap), and `busy_timeout = 5000`.

### Storage Tier Decoupling

* **Fast Tier (NVMe/SD):** Holds the application executable, SQLite database (`flan.db`), and cover art cache (`data/covers/`). UI queries and catalog browsing remain instantaneous.
* **Bulk Tier (USB HDD):** Holds heavy media files. Mechanical hard drives spin down when idle and wake only upon video stream requests (`/stream/...`).
* **Ghost Mount Defense:** A hidden marker file (`.flan-keep`) prevents creating an accidental empty database on the root drive when external mounts fail.

### Session Authentication & Instant Revocation

* **HMAC-SHA256 Cookies:** Sessions contain `userID:role:tokenVersion:issuedAt:signature`, verified in constant time.
* **Zero Database Hits:** Authenticated page views validate against an in-memory `token_version` map (`map[int]int`). Resetting a user's PIN increments their version, immediately invalidating old sessions across all devices without database hits.
* **Persistent Secret:** Auto-generated 32-byte secret stored in `data/.session_secret` (`0600` permissions).

---

## 4. Server Endpoints Summary

### Web Pages (Server-Rendered HTML)

| Method | Route | Description |
| :--- | :--- | :--- |
| `GET` | `/` | Dashboard with Continue Watching/Reading shelves and recently added media |
| `GET` | `/videos` | Movies and TV series catalog with genre filter pills |
| `GET` | `/show/{id}` | TV series detail view with season tabs and episode lists |
| `GET` | `/books` | Books and documents catalog with format and genre filters |
| `GET` | `/watch/{type}/{id}` | Video player page (Plyr) with auto-resume prompt (`type`: `movie` \| `episode`) |
| `GET` | `/read/{id}` | Document reader (native PDF iframe or ePub.js book viewer) |
| `GET` | `/login` | Profile selector ("Who is watching?") and numeric PIN keypad |
| `GET` | `/setup` | First-time onboarding wizard (permanently disabled once admin exists) |
| `GET` | `/settings` | System status, storage health, library scans, and user management |

### Media Delivery & Streaming

| Method | Route | Description |
| :--- | :--- | :--- |
| `GET` | `/stream/{type}/{id}` | Zero-copy byte range delivery via `sendfile` (HTTP 206) |
| `GET` | `/covers/{type}/{id}` | Serves locally cached cover images with long-lived browser caching |
| `GET` | `/static/*` | Serves embedded CSS, JS, player assets, and mascot SVG avatars |

### JSON Management & State APIs

| Category | Method | Route | Description |
| :--- | :--- | :--- | :--- |
| **Auth** | `POST` | `/api/setup` | Initializes admin account and default libraries |
| | `POST` | `/api/login` | Validates PIN, checks lockout, and issues signed cookie |
| | `POST` | `/api/logout` | Clears session cookie |
| **Users** | `GET` | `/api/users` | Lists profile tiles for selector and settings |
| | `POST` | `/api/users` | Admin creates a new profile |
| | `PUT` | `/api/users/{id}` | Updates profile avatar icon and accent color |
| | `PUT` | `/api/users/{id}/pin` | Changes PIN and increments token version |
| | `DELETE` | `/api/users/{id}` | Admin deletes a user profile |
| **Libraries** | `GET` | `/api/libraries` | Lists configured libraries with mount disk space stats |
| | `POST` | `/api/libraries` | Admin registers a new storage path and media type |
| | `DELETE` | `/api/libraries/{id}` | Admin removes a library directory from catalog |
| | `POST` | `/api/scan` | Triggers a scan across all libraries (single-flight) |
| | `POST` | `/api/libraries/{id}/scan` | Triggers a scan for a specific library |
| | `POST` | `/api/upload` | Streaming multipart upload direct to library disk path |
| **Catalog** | `GET` | `/api/movies` | Lists standalone movies (supports genre filter) |
| | `GET` | `/api/series` | Lists TV series cards (grouped by series title) |
| | `GET` | `/api/series/{id}` | Retrieves series details, seasons, and episodes |
| | `GET` | `/api/books` | Lists books (supports format and genre filters) |
| | `POST` | `/api/media/{type}/{id}/match` | Admin manual metadata override ("Fix Match") |
| | `DELETE` | `/api/media/{type}/{id}` | Admin removes an item from catalog |
| **Progress** | `GET` | `/api/progress/{type}/{id}` | Retrieves saved playback or reading position |
| | `POST` | `/api/progress` | Syncs current playback or reading progress |

---

## 5. Architectural References

* **Data & Storage:** [Database Schema](database.md) • [Storage Architecture](storage.md) • [Scraping Engine](scraper.md)
* **Security & Traffic:** [Rate Limiting](rate-limiting.md) • [Threat Model](threat-model.md)
* **Client & UI:** [Design System](client/design-system.md) • [Components](client/components.md) • [Pages](client/pages.md)
* **Operations:** [Compilation & Deployment](compilation.md) • [Testing Strategy](testing.md)
* **Diagrams:** [Data Flow](diagrams/data-flow.md) • [User Flows](diagrams/user-flows.md)
