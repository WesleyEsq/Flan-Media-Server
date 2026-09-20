# Software Specifications

Explains the software design of Flan Media Server in one summarized file.

## 1. Overview and Goals

Flan Media Server is a lightweight media streaming server written in Go, specifically built for underpowered Linux devices like single-board computers (SBCs). The primary engineering constraint is keeping the operational footprint around 15 to 20mb RAM while serving 1 to 5 users and 2 to 3 active streams.

The server operates through direct streaming without on-the-fly transcoding, kernel-level zero-copy data transfer, an embedded SQLite database with a streamlined three-table schema, server-rendered Go templates styled with a comfortable soft dark theme, and custom players for video and EPUB books.

---

## 2. Technology and Dependencies

+ **Language:** Go (1.22 or higher).
+ **Network & HTTP:** Go standard library net/http. No third-party web frameworks are used.
+ **Web Client:** Go standard library html/template for multi-page server rendering, paired with modular vanilla JavaScript and CSS.
+ **Players:**
  + Video: Plyr (lightweight HTML5 media player styled with custom CSS).
  + Books: Native browser PDF rendering via iframe/embed, and ePub.js for client-side EPUB reading with an immediate download option.
+ **Asset Packaging:** Go standard library embed.FS to package templates, styles, scripts, and player assets directly into the single binary executable.
+ **Database:** SQLite3 managed through database/sql. Uses Write-Ahead Logging (WAL) and limited page caching to keep memory low.
+ **Authentication:** Password/PIN hashing using bcrypt and HMAC-signed session cookies (avoiding database reads on every page load).
+ **External Dependencies:** Kept to an absolute minimum, adhering to Apache 2.0, MIT, or BSD licensing.

---

## 3. User Experience & Design System

The visual design is aimed at regular people who want an approachable, comfortable interface rather than a clinical dashboard or a harsh pitch-black screen.

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
   + Book cards show page progress (e.g. "Page 120 of 340").

---

## 4. Media Players: Video and Books

Default browser players look inconsistent across operating systems and lack critical media features. Flan provides tailored players while avoiding heavy library bloat:

### Custom Video Player (Plyr)
The video viewing page (/watch/{id}) embeds a tailored instance of Plyr styled with the soft dark and lavender theme:
+ **Custom Accent:** CSS variable overrides (--plyr-color-main: #bb9af7) match the server theme.
+ **Controls:** Touch-friendly scrub bar, 10-second skip forward/backward buttons, playback speed selection (0.5x to 2x), and fullscreen toggle.
+ **Subtitles:** Native support for WebVTT subtitle tracks with customizable text size and background styling.
+ **Auto-Resume Prompt:** If a user previously watched part of the video, a prominent prompt appears on start: "Resume from 24:12?".
+ **Progress Syncing:** player.js hooks into timeupdate events and sends a throttled update to POST /api/progress every 5 seconds.

### Document & Book Reading (Native PDF and Web ePub.js)
+ **PDF Documents:** Rendered directly using the browser's native PDF viewing engine via an iframe or embed element. This completely avoids bundling heavy PDF rendering engines (like the 8MB PDF.js distribution) into the server binary.
+ **EPUB Books:** The /read/{id} page embeds ePub.js to unpack and render chapters in the browser with dark mode styling, font size adjustments, and reading progress tracking. A prominent "Download EPUB" button allows users to open the book in their favorite native reading app on tablets or phones.

---

## 5. Subtitle Handling (SRT to WebVTT)

Because the server avoids on-the-fly video transcoding, subtitles cannot be burned into the video stream. Subtitles are delivered as separate text tracks via the HTML5 video player:

1. **Direct Disk Discovery:** When streaming a video, the server checks the host directory for sidecar subtitle files matching the video filename (for example, movie.mp4 and movie.en.srt or movie.es.vtt). No database entries are needed for subtitles.
2. **On-The-Fly WebVTT Conversion:** Browsers only support the WebVTT format (.vtt). If the subtitle file on disk is an .srt file, the endpoint /subtitles/{id} converts the SRT timestamps to WebVTT format on the fly. This string conversion is lightweight and operates with virtually zero memory overhead.

---

## 6. Media Hierarchy (Movies vs. TV Shows)

Treating every video file as an individual catalog item causes TV shows with dozens of episodes to flood the home grid. Flan establishes a clean hierarchy:

+ **Movies:** Standalone items displayed directly on the catalog grid.
+ **TV Shows (Series → Seasons → Episodes):**
  + The main catalog displays **one poster card** for the entire television series (grouped by series_title).
  + Clicking the series card opens the series view (/show/{title}) showing the series synopsis, overall rating, and a season selector tab (Season 1, Season 2).
  + Selecting a season displays a clean list of episode cards with episode titles, overview, and individual watch progress bars.

---

## 7. Web Scraping & Local-First Covers

The server enriches media with covers, ratings, and genre tags using a local-first approach. Detailed specifications are documented in [docs/scraper.md](file:///home/wess/Documents/MechanicalSpeak/Flan-Media-Server/docs/scraper.md).

Key features:
+ **Local-First Covers:** The scanner first checks if a local poster image (poster.jpg, cover.jpg, or folder.jpg) exists in the media folder, or if an EPUB contains an embedded cover image. If found, it uses it immediately without any network calls.
+ **TMDB Fallback:** If no local cover exists, the server queries The Movie Database (TMDB) for movies and TV shows.
+ **Streamed Local Storage:** Downloaded covers are streamed directly from remote HTTPS connections to local disk storage (data/covers/{media_id}.jpg) and served locally with long-lived browser caching.
+ **Genre Tags:** Genres are saved as a simple comma-separated string on the media item record (e.g. "Animation, Comedy"), allowing instant filtering via SQL.
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
+ **HMAC-Signed Session Cookies:** Authenticated sessions use signed cookies containing the user ID and timestamp. This avoids database lookups on every page request.
+ **Two-Tier Account Recovery:**
  1. **Standard Users:** Admin resets any user's PIN via the settings page.
  2. **Admin Terminal CLI Failsafe:** Running `./flan --reset-admin` from the host shell resets the admin PIN directly in SQLite. This eliminates the need for emergency recovery keys and extra recovery web endpoints.

---

## 10. File Intake Pipeline

Flan Media Server supports two intake methods:

+ **Local In-Place Scanning:** The server crawls configured media directories recursively. Files remain in place on disk. The server records the file size, format, and path into SQLite without moving or modifying files.
+ **Admin Web Uploads:** For remote headless servers, administrators can drag and drop folders or files into the browser. Incoming files are streamed directly from the network socket to disk in 32kb chunks via r.MultipartReader and io.Copy with pre-upload disk space verification. Memory usage remains constant even during multi-gigabyte transfers.

---

## 11. Memory and Concurrency Model (~15 to 20mb RAM)

+ **Zero-Copy Streaming:** Video delivery uses Go's http.ServeContent, delegating byte transfers directly to Linux sendfile. Media data moves straight from kernel page cache to socket without touching the Go application heap.
+ **Runtime Memory Ceilings:** GOMEMLIMIT=16MiB and GOGC=30 enforce disciplined garbage collection.
+ **Read-Only Binary Assets:** HTML templates, CSS, JS, and player libraries are embedded into the binary via embed.FS, residing in read-only memory rather than the application heap.
+ **Goroutines:** Each connection consumes ~2kb. 10 idle connections and 2 to 3 active streams consume less than 50kb of memory.

---

## 12. Database Schema and Storage Strategy

The database uses a clean, non-bloated three-table schema:
+ **users:** User profiles with bcrypt-hashed PINs, avatar colors, roles (admin/user), and lockout tracking.
+ **media_items:** Central catalog table storing movies, TV episodes, and books with metadata, genres, and series hierarchy.
+ **playback_progress:** Per-user playback positions and completion status.

Complete schema declarations, storage location advice for single-board computers, and core query procedures are documented in [docs/database.md](file:///home/wess/Documents/MechanicalSpeak/Flan-Media-Server/docs/database.md).

---

## 13. Server Endpoints

### Onboarding & Authentication
+ GET /setup : First-time setup wizard (disabled once users exist).
+ POST /api/setup : Initializes admin account and primary library.
+ GET /login : Renders profile selector and PIN entry screen.
+ POST /api/login : Validates PIN and issues signed session cookie.
+ POST /api/logout : Clears session cookie.

### Web Pages (Go Templates)
+ GET / : Main dashboard (Continue Watching, Continue Reading, Recent).
+ GET /videos : Videos catalog with Movies and TV tabs, and genre filter pills.
+ GET /show/{title} : Series detail view with season tabs and episode lists.
+ GET /books : Books catalog with genre and author filters.
+ GET /watch/{id} : Video player page with custom Plyr interface and subtitles.
+ GET /read/{id} : Document viewer with native PDF iframe or ePub.js.
+ GET /settings : Server settings, library paths, and user profiles.

### Media, Covers & Subtitle Streaming
+ GET /stream/{id} : Streams video/documents using HTTP 206 Partial Content and sendfile.
+ GET /covers/{id} : Serves locally cached cover images with long-lived browser caching.
+ GET /subtitles/{id} : Delivers WebVTT subtitle tracks (converting SRT on the fly).
+ GET /static/* : Serves embedded CSS, JS, player scripts, and icons.

### Management API (JSON)
+ GET /api/media : Lists catalog items (supports ?type=movie|tv|book&genre=...).
+ GET /api/media/{id} : Retrieves single item metadata.
+ GET /api/progress/{id} : Retrieves saved playback position for active user.
+ POST /api/progress : Saves current playback position or completion status.
+ POST /api/upload : Admin-only streaming multipart upload endpoint.
+ POST /api/scan : Triggers an immediate re-scan of configured media folders.
+ POST /api/media/{id}/match : Admin manual metadata override ("Fix Match").
