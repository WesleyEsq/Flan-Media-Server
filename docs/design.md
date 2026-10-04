# Master System Architecture

High-level architecture, operating constraints, and system design for Flan Media Server.

---

## 1. Operating Targets & System Constraints

Flan is designed for low-power host environments: repurposed PCs, small homelab devices, and single-board computers (such as Raspberry Pi and Orange Pi) serving 1 to 5 household members.

* **Resource Footprint:** Operates with low memory and CPU overhead using zero-copy file streaming, direct format playback, a pure-Go static binary with embedded assets, and a bounded SQLite page cache (~2 MB).
* **CPU Model:** Omits real-time video transcoding; all media streams directly to preserve host CPU and disk lifespan.
* **Navigation:** Consolidated strictly into a persistent left sidebar on desktop (`Video`, `Books`, and `Manage Server`), transforming into a fixed bottom bar on mobile screens. The top header contains brand identification and utility controls only.
* **Media Model:** Unifies video media (movies and series) into containers and file lists. Unifies books (single titles and multi-volume series) into containers and file lists.
* **Network Independence:** Fully offline operation with zero third-party metadata API calls. Directory names define titles; local `poster.jpg` files provide cover art.
* **Storage Architecture:** Multi-drive storage roots managed via the `storage_sources` table, pre-seeded with default roots `./media/video` and `./media/books` upon initial setup.
* **Direct Streaming:** Direct zero-copy file streaming via standard HTTP 206 Range requests delegating to kernel `sendfile` without CPU-heavy transcoding or complex proxy layers.
* **Packaging:** Single standalone static binary compiled with `CGO_ENABLED=0` using `modernc.org/sqlite`, with web templates and assets embedded via `embed.FS`.

---

## 2. Software Architecture (Layered Controller-Service-Repository)

The server implements a clean 4-tier layered architecture:

```text
HTTP Request → [Middleware Filters: Auth, CSRF]
                     ↓
             Controller (RegisterRoutes) 
                     ↓
                  Service (Domain Logic)
                     ↓
               Repository (SQL DAOs)
                     ↓
                 SQLite WAL

Views (web/templates): Rendered by PageController using domain ViewModels.
Domain Entities (internal/model): Pure data structures shared across layers.
```

* **Model Layer (`internal/model`):** Pure domain structs (`Video`, `VideoFile`, `Book`, `BookFile`, `User`, `Progress`) and sentinel errors (`ErrNotFound`, `ErrDuplicate`). Zero database imports or SQL execution.
* **Repository Layer (`internal/repository`):** Dedicated Data Access Objects (DAOs) executing pure SQL queries for the 7-table schema.
* **Service Layer (`internal/service`):** Business rules and workflow orchestration (`AuthService`, `MediaService`, `ScannerService`).
* **Controller Layer (`internal/controller`):** HTTP transport adapters, request decoders, status code mapping, ViewModel assembly, and route registration onto `http.ServeMux`.
* **Middleware Layer (`internal/middleware`):** Request filters: HMAC session authentication and CSRF/origin verification.
* **View Layer (`web/templates`, `web/static`):** 11 server-rendered HTML templates and vanilla CSS/JS packaged directly into the binary via `embed.FS`.

---

## 3. Core Technical Decisions

### Sidebar Navigation & Utility Controls

* **Header Controls:** The top header contains the platform brand title on the left, and utility controls on the right: the software manual (`/manual`) and user profile avatar. No general navigation links reside in the header.
* **Persistent Sidebar:** All main navigation occurs through three buttons: `Video` (`/video`), `Books` (`/books`), and `Manage Server` (`/manage`). On screens 768px wide or smaller, CSS transitions this sidebar into a fixed 64px bottom navigation bar.
* **User Profile Dialog:** Selecting the avatar opens a dialog where users can update their display name, select a preset SVG avatar or upload a custom image (max 2MB), modify their numeric PIN, or log out.

### Authentication & Bootstrap Token

* **Login Screen:** Unauthenticated requests route to `/login`, presenting a focused tactile keypad: household member profile avatar buttons, numeric PIN input, and instant access button.
* **Bootstrap Setup Token:** On initial startup with an empty user database, Flan prints a 6-character bootstrap setup token to stdout/systemd journal. Creating the initial administrator account requires this token to prevent unauthorized access across the local network.
* **CLI Recovery:** If an administrator forgets their PIN, running `./flan --reset-admin` directly on the host interactively prompts for and sets a new administrator PIN.

### Direct Streaming & External Player Integration

* **Direct File Delivery:** Media files stream in their native formats. Handlers delegate to Go's standard `http.ServeContent` with HTTP 206 Partial Content (Range requests), transferring bytes directly from filesystem cache to socket via Linux `sendfile` without userspace buffering or artificial lease limits.
* **External Player Links (VLC / MPV):** Because browsers do not support certain containers (such as MKV) or audio codecs (such as AC3 or DTS), Flan provides short-lived HMAC-signed URLs (`/download/{type}/{file_id}?exp=<unix>&u=<user_id>&sig=<hmac>`) valid for 4 hours. The HMAC signature binds `media_type:file_id:user_id:token_version:exp`, preventing cross-entity access and revoking links if the user PIN is changed.
* **Discrete Book Reading Progress:** PDF and EPUB reading progress is recorded as discrete states: `unread`, `reading`, and `finished`.

### Authentication Safeguards & Brute-Force Protection

* **PIN Lockout:** Tracks failed login attempts strictly in memory by Client IP and User ID. 5 consecutive failed attempts trigger a 5-minute lockout (`HTTP 429`) with exponential backoff (capped at 1 hour).
* **CPU Starvation Protection:** Concurrent bcrypt hashing operations are serialized via a single-slot worker queue (`chan struct{}`) with an artificial 2-second processing delay applied on the 4th failed attempt before acquiring the bcrypt slot.
* **Scan Deduplication:** Library scans are guarded by `singleflight.Group` to prevent redundant concurrent filesystem walks.

---

## 4. Server Endpoints

### HTML Views (11 Server-Rendered Templates)

* `GET /login`: Tactile profile keypad & login screen.
* `GET /setup`: Initial admin initialization screen (requires bootstrap token).
* `GET /` or `GET /video`: Video catalog card grid.
* `GET /video/{id}`: Video container detail view with file list and signed download links.
* `GET /books`: Books catalog card grid.
* `GET /books/{id}`: Book container detail view with volume/format list.
* `GET /search`: Full-canvas search page across videos and books with filter chips.
* `GET /watch/{file_id}`: Video player page (Plyr) with playback resume and subtitles.
* `GET /read/{file_id}`: Document viewer (PDF) or direct download prompt (EPUB).
* `GET /manage`: Administration console (system metrics, manual scan trigger, user profiles).
* `GET /manual`: System documentation and user manual.

### Media Delivery & Static Content

* `GET /stream/video/{file_id}`: Byte-range streaming via `sendfile` (HTTP 206).
* `GET /stream/book/{file_id}`: Serves book files with Range support.
* `GET /stream/subtitles/{file_id}/{track_id}`: Serves sidecar subtitles dynamically converted to WebVTT.
* `GET /download/{type}/{file_id}`: Direct download or external player playback (accepts session cookie or signed URL parameters).
* `GET /covers/{type}/{id}`: Cached cover artwork with long-lived client cache headers.
* `GET /avatars/{user_id}`: User profile avatar image.
* `GET /static/*`: Embedded CSS, JavaScript, and player assets.
* `GET /healthz`: Lightweight health check endpoint (`200 OK`).

### API Endpoints

* **Authentication:**
  * `POST /api/setup`: Initializes the admin account using the bootstrap token.
  * `POST /api/login`: Validates user PIN and issues HMAC session cookie.
  * `POST /api/logout`: Clears session cookie.
* **Playback & Progress:**
  * `GET /api/progress/{type}/{file_id}`: Retrieves saved playback position or reading status.
  * `POST /api/progress`: Saves current playback position or reading status.
* **User Management:**
  * `GET /api/users`: Public list of usernames and avatars for the login dropdown.
  * `POST /api/users`: Admin creates a new user profile.
  * `PUT /api/users/{id}`: Updates user display name or avatar selection.
  * `POST /api/users/{id}/avatar`: Uploads custom avatar image (PNG, JPEG, WebP, max 2MB).
  * `PUT /api/users/{id}/pin`: Changes user PIN and revokes prior tokens.
  * `DELETE /api/users/{id}`: Admin deletes a user account.
* **Storage Sources & Ingestion Pipeline:**
  * `GET /api/sources`: List configured storage roots with drive space stats.
  * `POST /api/sources`: Add a new storage folder (auto-initializes `.flan-keep` if missing).
  * `PUT /api/sources/{id}`: Update storage source friendly display name.
  * `DELETE /api/sources/{id}`: Remove a storage source.
  * `POST /api/sources/{id}/scan`: Scans source folder on-demand; returns gathered candidate items as JSON in response.
  * `POST /api/sources/{id}/ingest`: Commits user-selected items and edited titles into the catalog (`metadata_locked = 1`).
  * `POST /api/scan`: Maintenance scan of existing catalog items to flag missing files or update mtimes.
  * `GET /api/scan/status`: Reports active maintenance scan progress and counts.
  * `POST /api/upload`: Direct streaming upload (accepts `source_id`, `media_type`, title, and file body).
* **Media Metadata & File Editing:**
  * `PUT /api/media/{type}/{id}`: Edits container title, release year, and overview (`metadata_locked = 1`).
  * `POST /api/media/{type}/{id}/cover`: Uploads or replaces custom container cover artwork.
  * `PUT /api/media/video/{id}/files`: Batch updates episode titles, sequence order indices, and hidden flags.
  * `DELETE /api/media/{type}/{id}`: Hides or removes container from catalog.

---

## 5. Authorization Matrix

* **Public (Unauthenticated):**
  * `GET /login`, `POST /api/login`, `GET /static/*`, `GET /api/users`
  * `GET /setup`, `POST /api/setup` (first-run only when user count is 0)
  * `GET /healthz`
* **Authenticated User:**
  * Catalog & Views: `GET /`, `GET /video`, `GET /video/{id}`, `GET /books`, `GET /books/{id}`, `GET /search`, `GET /watch/*`, `GET /read/*`, `GET /manual`
  * Media delivery: `GET /stream/*`, `GET /covers/*`, `GET /avatars/*`
  * Direct download: `GET /download/*` (authenticated via session cookie)
  * Progress & Playback: `GET /api/progress/*`, `POST /api/progress`
  * Profile updates: `PUT /api/users/{id}/*` (self only)
* **External Players:**
  * `GET /download/*` (authenticated via signed query parameters `exp`, `u`, and `sig`)
* **Administrator Only:**
  * Administration Console: `GET /manage`
  * Account creation & deletion: `POST /api/users`, `DELETE /api/users/{id}`
  * Storage sources: `GET /api/sources`, `POST /api/sources`, `PUT /api/sources/{id}`, `DELETE /api/sources/{id}`
  * Ingestion & uploads: `POST /api/sources/{id}/scan`, `POST /api/sources/{id}/ingest`, `POST /api/scan`, `GET /api/scan/status`, `POST /api/upload`
  * Catalog & metadata editing: `PUT /api/media/*`, `POST /api/media/*/cover`, `PUT /api/media/video/*/files`, `DELETE /api/media/*`
  * User profile override: `PUT /api/users/{id}/*` (admin override)

---

## 6. Related Documentation

* [Directory Structure & Architecture](directories.md)
* [Database Schema & Queries](database.md)
* [Storage Architecture](storage.md)
* [Authentication Security & Rate Limiting](rate-limiting.md)
* [Security Threat Model](threat-model.md)
* [Compilation & Deployment](compilation.md)
* [Testing Strategy](testing.md)
