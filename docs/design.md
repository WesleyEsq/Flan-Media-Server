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
* **Storage Paths:** Fixed directory roots at `./media/video` and `./media/books`, avoiding dynamic multi-library tables.
* **Stream Concurrency:** Maximum 3 active concurrent video streams, enforced by a playback lease governor to prevent I/O thrashing on attached USB drives.
* **Packaging:** Single standalone static binary compiled with `CGO_ENABLED=0` using `modernc.org/sqlite`, with web templates and assets embedded via `embed.FS`.

---

## 2. Software Architecture (Layered Controller-Service-Repository)

The server implements a clean 4-tier layered architecture:

```text
HTTP Request → [Middleware Filters: Auth, CSRF, Stream Governor]
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
* **Repository Layer (`internal/repository`):** Dedicated Data Access Objects (DAOs) executing pure SQL queries for the 6-table schema.
* **Service Layer (`internal/service`):** Business rules and workflow orchestration (`AuthService`, `MediaService`, `ScannerService`, `GovernorService`).
* **Controller Layer (`internal/controller`):** HTTP transport adapters, request decoders, status code mapping, ViewModel assembly, and route registration onto `http.ServeMux`.
* **Middleware Layer (`internal/middleware`):** Request filters: HMAC session authentication, CSRF/origin verification, and playback stream capacity enforcement.
* **View Layer (`web/templates`, `web/static`):** 9 server-rendered HTML templates and vanilla CSS/JS packaged directly into the binary via `embed.FS`.

---

## 3. Core Technical Decisions

### Sidebar Navigation & Utility Controls

* **Header Controls:** The top header contains the platform brand title on the left, and utility controls on the right: the software manual (`/manual`) and user profile avatar. No general navigation links reside in the header.
* **Persistent Sidebar:** All main navigation occurs through three buttons: `Video` (`/video`), `Books` (`/books`), and `Manage Server` (`/manage`). On screens 768px wide or smaller, CSS transitions this sidebar into a fixed 56px bottom navigation bar.
* **User Profile Dialog:** Selecting the avatar opens a dialog where users can update their display name, select a preset SVG avatar or upload a custom image (max 2MB), modify their numeric PIN, or log out.

### Authentication & Bootstrap Token

* **Login Screen:** Unauthenticated requests route to `/login`, presenting a focused tactile keypad: household member profile avatar buttons, numeric PIN input, and instant access button.
* **Bootstrap Setup Token:** On initial startup with an empty user database, Flan prints a 6-character bootstrap setup token to stdout/systemd journal. Creating the initial administrator account requires this token to prevent unauthorized access across the local network.
* **CLI Recovery:** If an administrator forgets their PIN, running `./flan --reset-admin` directly on the host generates a recovery token.

### Direct Streaming & External Player Integration

* **Direct File Delivery:** Media files stream in their native formats. Handlers delegate to Go's `http.ServeContent` with HTTP 206 Partial Content (Range requests), transferring bytes directly from filesystem cache to socket via Linux `sendfile`.
* **External Player Links (VLC / MPV):** Because browsers do not support certain audio codecs (such as AC3 or DTS), Flan provides short-lived HMAC-signed URLs (`/download/video/{file_id}?exp=<unix>&u=<user_id>&sig=<hmac>`) valid for 4 hours. These URLs allow external players to stream without needing the browser's session cookie.
* **Discrete Book Reading Progress:** PDF and EPUB reading progress is recorded as discrete states: `unread`, `reading`, and `finished`.

### Concurrency and Brute-Force Safeguards

1. **Playback Lease Governor:** Tracks active video sessions by `(session_id, file_id)` with a 30-second activity lease. Caps total concurrent streams at 3 (`HTTP 503` with `Retry-After: 30` when exceeded). Multiple range requests or seeks from the same session do not consume additional slots.
2. **PIN Lockout:** Tracks failed login attempts by IP address and User ID. 5 consecutive failed attempts trigger a 5-minute lockout with exponential backoff on subsequent failures. Concurrent bcrypt hashing operations are serialized (maximum 1 concurrent hash) to avoid CPU exhaustion on low-power devices.

---

## 4. Server Endpoints

### HTML Views (9 Server-Rendered Templates)

* `GET /login`: Tactile profile keypad & login screen.
* `GET /setup`: Initial admin initialization screen (requires bootstrap token).
* `GET /` or `GET /video`: Video catalog card grid.
* `GET /video/{id}`: Video container detail view with file list and signed download links.
* `GET /books`: Books catalog card grid.
* `GET /books/{id}`: Book container detail view with volume/format list.
* `GET /watch/{file_id}`: Video player page (Plyr) with playback resume.
* `GET /read/{file_id}`: Document viewer (PDF) or direct download prompt (EPUB).
* `GET /manage`: Administration console (system metrics, manual scan trigger, user profiles).
* `GET /manual`: System documentation and user manual.

### Media Delivery & Static Content

* `GET /stream/video/{file_id}`: Byte-range streaming via `sendfile` (HTTP 206).
* `GET /stream/book/{file_id}`: Serves book files with Range support.
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
  * `POST /api/sources/{id}/scan`: Scans source folder on-demand; returns gathered candidate items with cleaned titles.
  * `POST /api/sources/{id}/ingest`: Commits user-selected items and edited titles into the catalog (`metadata_locked = 1`).
  * `POST /api/scan`: Maintenance scan of existing catalog items to flag missing files or update mtimes.
  * `GET /api/scan/status`: Reports active maintenance scan progress and counts.
  * `POST /api/upload`: Direct streaming upload (accepts `source_id`, `media_type`, title, and file body).
* **Media Metadata & File Editing:**
  * `PUT /api/media/{type}/{id}`: Edits container title, release year, and overview (`metadata_locked = 1`).
  * `POST /api/media/{type}/{id}/cover`: Uploads or replaces custom container cover artwork.
  * `PUT /api/media/video/{id}/files`: Batch updates episode titles, sequence order indices, and hidden flags.
  * `DELETE /api/media/{type}/{id}`: Hides or removes container from catalog.
* **Progress:**
  * `GET /api/progress/{type}/{file_id}`: Retrieves saved playback position or reading status.
  * `POST /api/progress`: Saves current playback position or reading status.

---

## 5. Authorization Matrix

* **Public (Unauthenticated):**
  * `GET /login`, `POST /api/login`, `GET /static/*`, `GET /api/users`
  * `GET /setup`, `POST /api/setup` (first-run only when user count is 0)
  * `GET /healthz`
* **Authenticated User:**
  * Catalog: `GET /video`, `GET /books`, `GET /watch/*`, `GET /read/*`
  * Media delivery: `GET /stream/*`, `GET /covers/*`, `GET /avatars/*`
  * Direct download: `GET /download/*` (authenticated via session cookie)
  * Progress: `GET /api/progress/*`, `POST /api/progress`
  * Profile updates: `PUT /api/users/{id}/*` (self only)
* **External Players:**
  * `GET /download/*` (authenticated via signed query parameters `exp`, `u`, and `sig`)
* **Administrator Only:**
  * Account creation & deletion: `POST /api/users`, `DELETE /api/users/{id}`
  * Storage sources: `GET /api/sources`, `POST /api/sources`, `DELETE /api/sources/{id}`
  * Ingestion & uploads: `GET /api/ingestion/*`, `POST /api/ingestion/*`, `POST /api/upload`
  * Catalog & metadata editing: `POST /api/scan`, `PUT /api/media/*`, `POST /api/media/*/cover`, `PUT /api/media/video/*/files`, `DELETE /api/media/*`
  * User profile override: `PUT /api/users/{id}/*` (admin override)

---

## 6. Related Documentation

* [Directory Structure & Architecture](directories.md)
* [Database Schema & Queries](database.md)
* [Storage Architecture](storage.md)
* [Rate Limiting & Throttling](rate-limiting.md)
* [Security Threat Model](threat-model.md)
* [Compilation & Deployment](compilation.md)
* [Testing Strategy](testing.md)
