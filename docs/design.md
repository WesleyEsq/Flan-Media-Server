# Software Specifications

Explains the software design of Flan Media Server in one summarized file.

## 1. Overview and Goals

Flan Media Server is a lightweight media streaming server written in Go, specifically built for underpowered Linux devices like single-board computers (SBCs). The primary engineering constraint is keeping the operational footprint around 15 to 20mb RAM while serving 1 to 5 users and 2 to 3 active streams.

The server operates through direct streaming without on-the-fly transcoding, kernel-level zero-copy data transfer, an embedded SQLite database, and server-rendered Go templates styled with a comfortable, soft dark theme.

---

## 2. Technology and Dependencies

+ **Language:** Go (1.22 or higher).
+ **Network & HTTP:** Go standard library net/http. No third-party web frameworks are used.
+ **Web Client:** Go standard library html/template for multi-page server rendering, paired with modular vanilla JavaScript and CSS.
+ **Asset Packaging:** Go standard library embed.FS to package templates, styles, scripts, and images directly into the single binary executable.
+ **Database:** SQLite3 managed through database/sql. Uses Write-Ahead Logging (WAL) and limited page caching to keep memory low.
+ **Authentication:** Password/PIN hashing using bcrypt and cryptographically secure session tokens via crypto/rand.
+ **External Dependencies:** Kept to an absolute minimum, adhering to Apache 2.0, MIT, or BSD licensing.

---

## 3. User Experience & Design System

The visual design is aimed at regular people who want an approachable, comfortable interface rather than a clinical dashboard or a harsh pitch-black screen.

### Color Palette: Soft Dark with Lavender
+ **Background:** Deep charcoal and slate (#16161e and #1a1b26) to reduce eye strain in low-light environments without the harsh contrast of pure black.
+ **Accents & Highlights:** Soft muted lavender (#bb9af7 and #7aa2f7) for buttons, active tabs, and playback indicators.
+ **Cards & Surfaces:** Slightly lighter slate (#24283b) with subtle rounded corners and clean borders.
+ **Typography:** Clean system sans-serif fonts (Inter, system-ui) with high-contrast, off-white text (#c0caf5) for easy reading on TVs, tablets, and phones.

### Core UX Flows
+ **"Who is watching?" Screen:** On startup, returning users see large, friendly profile tiles. Clicking a tile opens a clean numeric PIN prompt.
+ **Videos vs. Books Switcher:** Prominent top-level tabs allow immediate switching between video streams and document reading.
+ **Visual Progress Bars:** Poster cards display a subtle lavender progress bar showing how much time remains (for example, "35m left").
+ **Responsive Media Player:** Fullscreen HTML5 video player with keyboard shortcuts, touch scrub bars, and playback speed adjustments.

---

## 4. First-Time Setup Wizard

When the server runs for the first time with an empty database:

1. Any request to the server detects zero users and redirects to /setup.
2. The setup page prompts the administrator for:
   + Admin username (e.g. Wesley).
   + 4 to 6-digit numeric PIN.
   + Initial media library directory path on the host.
3. When submitted, the server:
   + Validates the directory path and checks read permissions.
   + Hashes the PIN using bcrypt.
   + Generates a random 16-character emergency recovery key.
   + Saves the admin account and primary library into SQLite.
   + Launches an initial background library scan.
   + Issues an authenticated session cookie and displays the recovery key with a button to enter the catalog.
4. **Permanent Lockout:** As soon as one user exists in the database, the /setup route is permanently deactivated.

---

## 5. PIN Authentication and Account Recovery

### Profile-Based PINs
To ensure the client is fast and friendly on TVs and mobile devices, users log in using a 4 to 6-digit PIN.

### Brute-Force Protection
A 4-digit PIN has only 10,000 combinations. The server enforces strict progressive rate limiting:
+ Failed attempts are tracked in memory per profile and client IP.
+ After 5 consecutive failures, the profile is locked for 5 minutes.
+ Subsequent failures double the lockout duration.

### Three-Tier Account Recovery Model
Because the server operates locally without an internet email connection, recovery is handled locally:
1. **Standard Users:** The administrator can change or reset any regular user's PIN directly in the admin settings page.
2. **Admin Emergency Key:** If the admin forgets their PIN, a "Forgot PIN" link on the login page accepts the 16-character emergency recovery key generated during setup.
3. **Command-Line Failsafe:** If the admin loses both their PIN and their emergency key, physical or SSH access to the host allows running:
   ```bash
   ./flan --reset-admin
   ```
   This interactive CLI command prompts for a new admin PIN and updates SQLite directly.

---

## 6. File Intake Pipeline

Flan Media Server supports two distinct methods for introducing media into the system:

### Method A: Local In-Place Library Scanning (Existing Collections)
For users who already have large organized folders on external drives or NAS mounts:
+ Admins can define multiple library paths in settings (e.g. /mnt/storage/movies for videos, /mnt/storage/books for documents).
+ The server crawls the directory tree recursively.
+ Files remain in place on disk. The server records the file size, format, and path into SQLite without moving or modifying files.
+ Memory usage during scanning is negligible (a few kilobytes).

### Method B: Admin Web Upload (Streaming Intake)
For headless remote servers where the administrator does not have local terminal or file sharing access:
+ The web client allows dragging and dropping files or entire directory trees (using HTML5 webkitdirectory).
+ **Admin Only:** File uploads are restricted to the admin profile.
+ **Pre-Upload Disk Check:** Before accepting an upload, the server queries free disk space via statfs. If free space is below a safety threshold (e.g. 2gb), the upload is rejected with HTTP 507 Insufficient Storage.
+ **Zero-Memory Streaming:** Incoming multipart files are read using r.MultipartReader and streamed directly from the network socket to destination files on disk using io.Copy in 32kb buffers. Memory usage does not spike, even when uploading 10gb+ files.

---

## 7. Memory and Concurrency Model (~15 to 20mb RAM)

### Low Memory Enforcement
1. **Zero-Copy Streaming:** Video streaming uses Go's http.ServeContent, delegating transfers to the Linux sendfile system call. Media bytes move directly from the kernel page cache to the network socket, keeping the Go heap free of large media buffers.
2. **Runtime Memory Tuning:**
   + GOMEMLIMIT=16MiB triggers garbage collection before memory approaches 16mb.
   + GOGC=30 runs garbage collection in frequent, smaller increments.
3. **Read-Only Binary Assets:** HTML templates, CSS, and JS files are embedded into the binary using embed.FS, so they reside in read-only memory rather than the application heap.
4. **Goroutines:** Each connection runs in a 2kb goroutine. Handling 10 idle connections and 2 to 3 active streams consumes under 50kb of connection memory.

---

## 8. Database Schema and SQLite Configuration

SQLite3 is configured with Write-Ahead Logging and strict cache limits:

```sql
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = -2000;
PRAGMA temp_store = MEMORY;
PRAGMA foreign_keys = ON;
```

### Tables

```sql
-- Library Folders Configured by Admin
CREATE TABLE IF NOT EXISTS libraries (
    library_id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name          TEXT NOT NULL,
    path          TEXT NOT NULL UNIQUE,
    media_type    TEXT NOT NULL CHECK(media_type IN ('video', 'book')),
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- User Profiles
CREATE TABLE IF NOT EXISTS users (
    user_id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username           TEXT NOT NULL UNIQUE,
    pin_hash           TEXT NOT NULL,
    role               TEXT NOT NULL CHECK(role IN ('admin', 'user')),
    recovery_key_hash  TEXT,
    avatar_color       TEXT DEFAULT '#bb9af7',
    created_at         DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Active User Sessions
CREATE TABLE IF NOT EXISTS sessions (
    session_token TEXT PRIMARY KEY,
    user_id       INTEGER NOT NULL,
    expires_at    DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Media Catalog Items
CREATE TABLE IF NOT EXISTS media_items (
    media_id      INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id    INTEGER NOT NULL,
    title         TEXT NOT NULL,
    type          TEXT NOT NULL CHECK(type IN ('video', 'book')),
    format        TEXT NOT NULL,
    file_path     TEXT NOT NULL UNIQUE,
    file_size     INTEGER NOT NULL,
    duration      INTEGER DEFAULT 0,
    cover_path    TEXT,
    overview      TEXT,
    release_year  INTEGER,
    added_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (library_id) REFERENCES libraries(library_id) ON DELETE CASCADE
);

-- Playback and Reading Progress
CREATE TABLE IF NOT EXISTS playback_progress (
    user_id           INTEGER NOT NULL,
    media_id          INTEGER NOT NULL,
    position_seconds  INTEGER NOT NULL DEFAULT 0,
    total_seconds     INTEGER NOT NULL DEFAULT 0,
    is_finished       INTEGER NOT NULL DEFAULT 0,
    updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, media_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_media_type ON media_items(type);
CREATE INDEX IF NOT EXISTS idx_progress_updated ON playback_progress(updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(session_token);
```

---

## 9. Server Endpoints

### Onboarding & Authentication
+ GET /setup : First-time setup wizard (disabled once users exist).
+ POST /api/setup : Initializes admin account and primary library.
+ GET /login : Renders profile selector and PIN entry screen.
+ POST /api/login : Validates PIN and issues session cookie.
+ POST /api/logout : Clears session cookie and invalidates token.
+ POST /api/recover : Validates admin emergency key and resets admin PIN.

### Web Pages (Go Templates)
+ GET / : Main catalog view (videos, books, continue watching).
+ GET /watch/{id} : Video player page with auto-resume.
+ GET /read/{id} : Document viewer for PDF and EPUB files.
+ GET /settings : Server settings, library management, and user profiles.

### Media & Static Streaming
+ GET /stream/{id} : Streams video/documents using HTTP 206 Partial Content and Linux sendfile.
+ GET /static/* : Serves embedded CSS, JS, and image assets.

### Management API (JSON)
+ GET /api/media : Lists catalog items (supports ?type=video|book).
+ GET /api/media/{id} : Retrieves single item metadata.
+ GET /api/progress/{id} : Retrieves saved playback position for active user.
+ POST /api/progress : Saves current playback position or completion status.
+ POST /api/upload : Admin-only streaming multipart upload endpoint.
+ GET /api/libraries : Lists configured library paths.
+ POST /api/libraries : Adds a new library directory and initiates scan.
+ POST /api/scan : Triggers an immediate re-scan of all configured libraries.
