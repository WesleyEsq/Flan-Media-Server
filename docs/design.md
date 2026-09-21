# Software Specifications

Explains the software design of Flan Media Server in one summarized file.

## 1. Overview and Goals

Flan Media Server is a lightweight media streaming server written in Go, specifically built for underpowered Linux devices like single-board computers (SBCs). The primary engineering constraint is keeping the operational footprint around 15 to 20mb RAM while serving 1 to 5 users and 2 to 3 active streams.

The server operates through direct streaming without on-the-fly transcoding, kernel-level zero-copy data transfer, an idiomatic **Model-View-Controller (MVC)** software architecture, an embedded SQLite database with a normalized relational schema, server-rendered Go templates styled with a comfortable soft dark theme, and custom players for video and EPUB books.

---

## 2. Technology and Dependencies

+ **Language:** Go (1.22 or higher).
+ **Architecture:** Model-View-Controller (MVC) with decentralized route mounting:
  + **Model (`internal/model`):** Domain entities, sentinel errors, and database query implementations.
  + **View (`web/templates`):** Server-rendered Go HTML templates styled with vanilla CSS.
  + **Controller (`internal/controller`):** HTTP transport, status codes, ViewModel assembly, and route registration.
  + **Middleware (`internal/middleware`):** Composable HTTP filters for auth, rate limiting, and stream capacity.
+ **Network & HTTP:** Go standard library `net/http` using native Go 1.22+ method routing. No third-party web frameworks are used.
+ **Web Client:** Go standard library `html/template` for multi-page server rendering, paired with modular vanilla JavaScript and CSS.
+ **Players:**
  + Video: Vendored Plyr (lightweight HTML5 media player styled with custom CSS).
  + Books: Native browser PDF rendering via iframe/embed, and vendored ePub.js with JSZip for client-side EPUB reading with an immediate download option.
+ **Asset Packaging:** Go standard library embed.FS to package templates, styles, scripts, and vendored player assets directly into the single binary executable with zero external CDN dependencies.
+ **Database:** SQLite3 managed through `database/sql` using the pure-Go `modernc.org/sqlite` driver (zero CGo, allowing direct cross-compilation to ARMv6, ARMv7, and ARM64). Dedicated configuration and initialization module (`internal/database/database.go`) manages connection pooling (`SetMaxOpenConns(1)`), wear-leveling pragmas (`WAL`, `cache_size = -2000`), `.flan-keep` verification, and schema DDL migrations.
+ **Authentication:** Password/PIN hashing using bcrypt and HMAC-SHA256-signed session cookies backed by an automatically persisted 32-byte secret key (avoiding database reads on every page load).
+ **External Dependencies:** Kept to an absolute minimum, adhering to Apache 2.0, MIT, or BSD licensing.

---

## 3. User Experience & Design System

The visual design is aimed at regular people who want an approachable, comfortable interface rather than a clinical dashboard or a harsh pitch-black screen. Complete design tokens, components, and template wireframes are documented in [docs/client/design-system.md](docs/client/design-system.md), [docs/client/pages.md](docs/client/pages.md), and [docs/client/components.md](docs/client/components.md).

### Color Palette: Soft Dark with Lavender

+ **Background:** Deep charcoal and slate (#16161e and #1a1b26) to reduce eye strain in low-light environments without the harsh contrast of pure black.
+ **Accents & Highlights:** Soft muted lavender (#bb9af7 and #7aa2f7) for buttons, active tabs, and playback indicators.
+ **Cards & Surfaces:** Slightly lighter slate (#24283b) with subtle rounded corners and clean borders.
+ **Typography:** Clean system sans-serif fonts (Inter, system-ui) with high-contrast, off-white text (#c0caf5) for easy reading on TVs, tablets, and phones.

### Interface Division: Videos and Books

Mixing movies, multi-season TV shows, and books on the same screen creates confusion. The interface separates media intuitively:

1. **Top Navigation Tabs:**
   + **Home:** Displays a "Continue Watching" shelf for in-progress videos, a "Continue Reading" shelf for in-progress books, and recently added media.
   + **Videos:** Contains sub-tabs for **Movies** and **TV Shows**, with genre filter pills (Animation, Comedy, Sci-Fi, Drama).
   + **Books:** Grid of books with author and genre tags.
2. **Distinct Card Aspect Ratios:**
   + **Videos:** Standard 2:3 movie poster aspect ratio (e.g. 200px by 300px) or 16:9 landscape cards for episodes.
   + **Books:** Standard 1:1.4 book cover aspect ratio with a subtle faux book spine shadow on the left edge.
3. **Dedicated Progress Bars:**
   + Video cards show remaining duration (e.g. "35m left").
   + Book cards show reading progress (e.g. "Page 120 of 340" or "45%").

### Profile Avatars & Customization

Users can personalize their profile tiles:

+ **Curated Built-In Icons:** A collection of lightweight, embedded SVG icons included directly in the binary (Flan the hamster, popcorn bowl, retro TV, cat, dog, book, robot, and cassette tape).
+ **Custom Accent Colors:** Profiles pair their icon with a customizable background color (defaulting to soft lavender #bb9af7).
+ **Custom Uploads:** Users can also upload a personal photo or image to serve as their avatar.

---

## 4. Media Players: Video and Books

Default browser players look inconsistent across operating systems and lack critical media features. Flan provides tailored players while avoiding heavy library bloat:

### Custom Video Player (Plyr)

The video viewing page (/watch/{type}/{id}) embeds a tailored instance of Plyr styled with the soft dark and lavender theme:

+ **Custom Accent:** CSS variable overrides (--plyr-color-main: #bb9af7) match the server theme.
+ **Controls:** Touch-friendly scrub bar, 10-second skip forward/backward buttons, playback speed selection (0.5x to 2x), and fullscreen toggle.
+ **Auto-Resume Prompt:** If a user previously watched part of the video, a prominent prompt appears on start: "Resume from 24:12?".
+ **Progress Syncing:** player.js hooks into timeupdate events and sends a throttled update to POST /api/progress every 5 seconds.

### Document & Book Reading (Native PDF and Web ePub.js)

+ **PDF Documents:** Rendered directly using the browser's native PDF viewing engine via an iframe or embed element. This completely avoids bundling heavy PDF rendering engines (like the 8MB PDF.js distribution) into the server binary.
+ **EPUB Books:** The /read/{id} page embeds vendored ePub.js and JSZip to unpack and render chapters in the browser with dark mode styling, font size adjustments, and reading progress tracking. A prominent "Download EPUB" button allows users to open the book in their favorite native reading app on tablets or phones.
+ **Offline Asset Vendoring:** All third-party libraries (Plyr, ePub.js, JSZip) are stored directly within `web/static/vendor/` and embedded into the binary using Go's `embed.FS`. The web client makes zero runtime network calls to external CDNs, maintaining full functionality in air-gapped or offline homelab environments.

---

## 5. Multi-Directory Storage Architecture (Libraries)

Instead of relying on rigid, single-path environment configurations, Flan treats storage directories as first-class dynamic entities in SQLite via the `libraries` table:

1. **Multi-Mount Flexibility:** Users can map separate storage pools or physical drives to designated library records (e.g. `/mnt/hdd1/movies` for movies, `/mnt/hdd2/tv` for shows, and `/mnt/nvme/books` for books).
2. **Deterministic Typing:** Each library specifies its `media_type` (`movies`, `tv`, `books`). This eliminates guessing during library scanning and ensures accurate metadata scraping.
3. **Graceful Defaults:** If `.env` leaves `MEDIA_DIR` blank or omitted, Flan boots without crashing and initializes default directories at `./media/movies`, `./media/tv`, and `./media/books` during the `/setup` wizard.
4. **Dynamic Management:** Administrators can register or remove libraries dynamically via the `/settings` page and JSON API without needing to restart the daemon.

---

## 6. Media Hierarchy (Movies vs. TV Shows)

Treating every video file as an individual catalog item causes TV shows with dozens of episodes to flood the home grid. Flan establishes a clean hierarchy:

+ **Movies:** Standalone items displayed directly on the catalog grid.
+ **TV Shows (Series → Seasons → Episodes):**
  + The main catalog displays **one poster card** for the entire television series (grouped by series_title).
  + Clicking the series card opens the series view (/show/{id}) showing the series synopsis, overall rating, and a season selector tab (Season 1, Season 2).
  + Selecting a season displays a clean list of episode cards with episode titles, overview, and individual watch progress bars.

---

## 7. Web Scraping & Local-First Covers

The server enriches media with covers, ratings, and genre tags using a local-first approach. Detailed specifications are documented in [docs/scraper.md](file:///home/wess/Documents/MechanicalSpeak/Flan-Media-Server/docs/scraper.md).

Key features:

+ **Local-First Covers:** The scanner first checks if a local poster image (poster.jpg, cover.jpg, or folder.jpg) exists in the media folder, or if an EPUB contains an embedded cover image. If found, it uses it immediately without any network calls.
+ **TMDB Fallback:** If no local cover exists, the server queries The Movie Database (TMDB) for movies and TV shows.
+ **Streamed Local Storage:** Downloaded covers are streamed directly from remote HTTPS connections to local disk storage (data/covers/{type}/{id}.jpg) and served locally with long-lived browser caching.
+ **Genre Tags:** Genres are saved in normalized relational tables (`genres` and `item_genres`), allowing instant indexed filtering via SQL without full-table string scans.
+ **Admin Fix Match:** Administrators can search manually by title or TMDB ID, upload custom poster images, or edit metadata by hand.

---

## 8. First-Time Setup Wizard

When the server boots with an empty database:

1. Any request redirects to /setup.
2. The setup wizard prompts for:
   + Admin username (e.g. Wesley).
   + 4 to 6-digit numeric PIN.
   + Primary media directory path on the host.
3. The server validates the directory, creates the admin user with a bcrypt-hashed PIN, starts the background scan, sets an HMAC-signed session cookie, and opens the catalog.
4. Once an admin user exists, /setup is permanently disabled.

---

## 9. PIN Authentication and Account Recovery

+ **Profile Selection ("Who is watching?"):** Users choose their profile tile and enter their 4 to 6-digit PIN.
+ **Brute-Force Lockout:** After 5 failed attempts, the profile is locked for 5 minutes with exponential backoff on further failures.
+ **HMAC-Signed Session Cookies:** Authenticated sessions use signed cookies containing the payload `userID:role:issuedAt:signature` generated with HMAC-SHA256. Signatures are verified in constant time (`hmac.Equal`) and checked against a 30-day expiration window, eliminating database lookups on page views.
+ **Persistent Secret Management:** The HMAC secret is loaded from `SESSION_SECRET` or read from a persistent `0600`-permission `.session_secret` file in the database directory (auto-generated on first boot via `crypto/rand`). This ensures user sessions remain valid across server restarts without manual intervention.
+ **Two-Tier Account Recovery:**
  1. **Standard Users:** Admin resets any user's PIN via the settings page.
  2. **Admin Terminal CLI Failsafe:** Running `./flan --reset-admin` from the host shell resets the admin PIN directly in SQLite. This eliminates the need for emergency recovery keys and extra recovery web endpoints.

---

## 10. File Intake Pipeline & Multi-Drive Resiliency

Flan Media Server supports scattered storage across multiple directories and drives (SD, NVMe, USB HDD):

+ **First-Class Libraries:** Storage folders are managed through the `libraries` table in SQLite, each assigned a specific `media_type` (`movies`, `tv`, `books`).
+ **Local In-Place Scanning:** The server crawls configured library paths recursively. Files remain in place on disk, recorded with paths relative to their library root.
+ **Deterministic Folder Structure:**
  + Movies: `<library_path>/Movie Title (Year).mp4` or `<library_path>/Movie Title (Year)/Movie Title (Year).mp4`.
  + TV Series: `<library_path>/<Series Title>/Season <NN>/<Series Title> - S<NN>E<NN> - <Title>.<ext>`.
  + Books: `<library_path>/<Author>/<Book Title>.<ext>` or `<library_path>/<Book Title>.<ext>`.
+ **Mount Liveness Safeguard:** If an external drive disconnects or is unmounted, the scanner detects that the directory is empty or absent and skips it entirely, preserving the catalog in SQLite without wiping records. When a user streams an offline item, the server returns HTTP 503 Service Unavailable ("Media drive is offline").
+ **Admin Web Uploads & Routing:** Administrators can upload files or folders via the web client. The modal allows selecting the target library. For TV shows, the modal captures Series Title and Season Number (or preserves folder structures via folder uploads), placing files directly into the correct season folder. The server validates free space via statfs before streaming incoming files directly to disk via `r.MultipartReader` in 32kb chunks.
+ **Ghost Database Prevention:** Uses marker files (.flan-keep) to prevent accidentally creating empty databases on root boot drives when external mounts fail. Detailed multi-drive guidelines are in [docs/storage.md](docs/storage.md).

---

## 11. Memory, Concurrency & Stream Governor (~15 to 20mb RAM)

+ **Zero-Copy Streaming:** Video delivery uses Go's http.ServeContent, delegating byte transfers directly to Linux sendfile. Media data moves straight from kernel page cache to socket without touching the Go application heap.
+ **Stream Governor (Disk Thrashing Defense):** Limits active concurrent streams via semaphore (MAX_CONCURRENT_STREAMS=3). This protects mechanical USB hard drive read heads from seeking thrashing, guaranteeing stutter-free streaming.
+ **Traffic Rate Limiting:** Enforces five rate-limiting zones covering stream capacity, PIN brute-forcing, API token buckets, scan cooldowns, and scraper pacing. Complete details are in [docs/rate-limiting.md](docs/rate-limiting.md).
+ **Runtime Memory Ceilings:** GOMEMLIMIT=16MiB and GOGC=30 enforce disciplined garbage collection.
+ **Read-Only Binary Assets:** HTML templates, CSS, JS, and player libraries are embedded into the binary via embed.FS, residing in read-only memory rather than the application heap.
+ **Goroutines:** Each connection consumes ~2kb. 10 idle connections and 2 to 3 active streams consume less than 50kb of memory.

---

## 12. Database Schema and Storage Strategy

The database uses a clean, normalized relational schema evaluated up to 4NF:

+ **users:** Profiles with bcrypt-hashed PINs, avatars, roles, and lockout tracking.
+ **libraries:** Configured storage roots with explicit media types (`movies`, `tv`, `books`).
+ **movies:** Standalone films with release year, duration, rating, overview, and cover path.
+ **series & episodes:** Dedicated tables for TV series hierarchy and episode files.
+ **books:** Books and documents with author, overview, format, and cover path.
+ **genres & item_genres:** Normalized genres eliminating full-table string matching.
+ **video_progress & book_progress:** Dedicated progress tracking for videos (seconds) and books (EPUB CFI / PDF page counts).

Complete schema declarations, storage location advice for single-board computers, and core query procedures are documented in [docs/database.md](docs/database.md).

---

## 13. Server Endpoints

### Onboarding & Authentication

+ GET /setup : First-time setup wizard (disabled once users exist).
+ POST /api/setup : Initializes admin account and initial library.
+ GET /login : Renders profile selector and PIN entry screen.
+ POST /api/login : Validates PIN and issues signed session cookie.
+ POST /api/logout : Clears session cookie.
+ GET /api/users : Lists profile tiles for selector and settings.
+ POST /api/users : Admin creates new user profile.
+ PUT /api/users/{id} : Updates username, avatar icon, and avatar color.
+ PUT /api/users/{id}/pin : Changes or resets profile PIN.
+ DELETE /api/users/{id} : Admin deletes a user profile.

### Web Pages (Go Templates)

+ GET / : Main dashboard (Continue Watching, Continue Reading, Recent).
+ GET /videos : Videos catalog with Movies and TV tabs, and genre filter pills.
+ GET /show/{id} : Series detail view with season tabs and episode lists.
+ GET /books : Books catalog with genre and author filters.
+ GET /watch/{type}/{id} : Video player page with custom Plyr interface (`type`: `movie` or `episode`).
+ GET /read/{id} : Document viewer with native PDF iframe or ePub.js.
+ GET /settings : Server settings, library management, storage health, and user profiles.

### Media & Covers Streaming

+ GET /stream/{type}/{id} : Streams video/documents using HTTP 206 Partial Content and sendfile.
+ GET /covers/{type}/{id} : Serves locally cached cover images with long-lived browser caching.
+ GET /static/* : Serves embedded CSS, JS, player scripts, and icons.

### Management & Catalog APIs (JSON)

+ GET /api/libraries : Lists configured libraries with disk space stats.
+ POST /api/libraries : Admin registers a new library directory and media type.
+ DELETE /api/libraries/{id} : Admin removes a library directory.
+ POST /api/libraries/{id}/scan : Triggers an immediate re-scan of a specific library.
+ POST /api/scan : Triggers an immediate re-scan across all configured libraries.
+ POST /api/upload : Admin streaming multipart upload with library and series routing.
+ GET /api/movies : Lists standalone movies with optional genre filters.
+ GET /api/series : Lists TV series catalog cards.
+ GET /api/series/{id} : Retrieves single series with season list and episodes.
+ GET /api/books : Lists books with format and genre filters.
+ GET /api/progress/{type}/{id} : Retrieves saved playback or reading position.
+ POST /api/progress : Saves current playback position or completion status.
+ PUT /api/media/{type}/{id} : Admin edits title, year, genres, overview.
+ POST /api/media/{type}/{id}/match : Admin manual metadata override ("Fix Match").
+ DELETE /api/media/{type}/{id} : Admin deletes item from catalog.

---

### Related Documentation

+ [System Directories & Package Anatomy](docs/directories.md)
+ [Database Schema & Storage Procedures](docs/database.md)
+ [Cross-Compilation & SBC Deployment Guide](docs/compilation.md)
+ [Storage Architecture & Mount Resiliency](docs/storage.md)
+ [Five-Zone Rate Limiting Architecture](docs/rate-limiting.md)
+ [Security Threat Model & Defensive Posture](docs/threat-model.md)
+ [Testing Strategy & TDD Guidelines](docs/testing.md)
+ [Data Flow & Sequence Diagrams](docs/diagrams/data-flow.md)
+ [User Flows & Navigation Journeys](docs/diagrams/user-flows.md)
