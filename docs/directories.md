# Directories in the Project

The project is organized according to standard Go conventions. This keeps application entry points, internal packages, web assets, and documentation cleanly separated.

## Directory Structure

``` text
Flan-Media-Server/
├── cmd/
│   └── flan/               # Server entry point (main.go)
├── internal/               # Private application packages
│   ├── config/             # .env and environment variable loading
│   ├── database/           # SQLite connection, 3-table schema, and queries
│   ├── handler/            # HTTP handlers for auth, pages, streaming, upload, and API
│   ├── scraper/            # Local-first artwork finder, TMDB scraper, cover downloader
│   └── subtitle/           # Sidecar subtitle discovery and on-the-fly SRT to WebVTT converter
├── web/                    # Embedded web assets (via embed.FS)
│   ├── templates/          # Server-rendered Go HTML templates
│   └── static/             # Static web files served at /static/
│       ├── css/            # Stylesheets (soft dark theme with lavender accents)
│       ├── js/             # Modular vanilla JS (API caller, player and reader wrappers)
│       ├── vendor/         # Lightweight embedded player assets (Plyr, ePub.js)
│       └── assets/         # Mascots, icons, and default cover art
├── docs/                   # Specifications, architecture, and diagrams
│   ├── design.md           # Full system specification and architecture
│   ├── database.md         # Database schema, storage location, and SQL query procedures
│   ├── storage.md          # Multi-drive architecture, mount defenses, and backup snapshots
│   ├── scraper.md          # Web scraping pipeline, APIs, local cover storage, and tags
│   ├── threat-model.md     # Security posture, attack vectors, and mitigations
│   ├── testing.md          # TDD guidelines, in-memory fs, and database decoupling
│   ├── directories.md      # Codebase organization and package guide
│   └── diagrams/           # User flows and architectural diagrams
│       └── user-flows.md   # Mermaid diagrams for setup, auth, streaming, intake
├── tests/                  # Unit and integration tests
├── .env                    # Default environment configuration
├── LICENSE                 # Apache 2.0 License
└── README.md               # Project overview and instructions
```

## Directory Explanations

### cmd/flan

Contains main.go, the main executable entry point for the server. It parses command-line flags (including the --reset-admin CLI failsafe), reads the .env file, initializes the SQLite database connection, registers HTTP routes, and handles graceful shutdown when receiving interrupt signals.

### internal/config

Handles reading configuration settings from the .env file and system environment variables. Provides typed configuration values with sensible defaults for server port, host, media directory paths, database file location, and memory tuning parameters.

### internal/database

Manages the SQLite3 database connection. Applies performance pragmas such as WAL mode, normal synchronous disk writes, and a 2mb page cache limit. Handles creating the 3-table schema (users, media_items, playback_progress) and indexes, and provides query functions for authentication, media catalog indexing, and watch progress.

### internal/handler

Implements the HTTP request handlers and routing logic:

+ **Auth & Setup Handlers:** Manages first-time setup (/setup), user profile PIN authentication (/login), HMAC-signed session cookies, and brute-force rate limiting.
+ **Page Handlers:** Renders the Go templates (dashboard, videos, series details, books, watch, read, and settings).
+ **Streaming Handler:** Handles video and document streaming using http.ServeContent. On Linux, this uses the sendfile system call to transfer byte ranges directly from the kernel page cache to the client socket without copying data into Go heap memory.
+ **Upload Handler:** Implements zero-memory streaming file and folder uploads via r.MultipartReader and io.Copy with pre-upload disk space verification.
+ **Cover & Subtitle Handlers:** Serves locally cached poster images and delivers on-the-fly converted WebVTT subtitles.
+ **REST API Handlers:** Provides JSON endpoints for querying catalog items, retrieving playback positions, saving watch progress, triggering library scans, and fixing metadata matches.

### internal/scraper

Implements the local-first scraping pipeline:

+ Inspects local folders for existing poster.jpg, cover.jpg, or embedded EPUB covers.
+ Sanitizes filenames using clean regex patterns to extract titles, years, seasons, and episodes.
+ Queries TMDB for movie and TV metadata when local artwork is missing.
+ Streams cover artwork directly from remote HTTPS connections to local disk storage without memory buffering.

### internal/subtitle

Discovers sidecar subtitle files (.srt and .vtt) in the video directory on disk and provides on-the-fly conversion from SRT format to browser-compatible WebVTT format without database overhead.

### web/templates

Contains the Go HTML templates used to render web pages:

+ **base.html:** Common layout template containing the HTML shell, navigation bar, logo, and footer.
+ **setup.html:** First-time onboarding wizard to create the admin account and register the media folder.
+ **login.html:** Profile selector ("Who is watching?") and numeric PIN keypad.
+ **index.html:** Home dashboard featuring Continue Watching, Continue Reading, and Recently Added shelves.
+ **videos.html:** Videos catalog with Movies and TV tabs, and genre filter pills.
+ **show.html:** TV series detail view with season tabs and episode lists.
+ **books.html:** Document and book catalog with genre and author filter pills.
+ **watch.html:** Video player view embedding Plyr with custom lavender styling, subtitles, and resume prompt.
+ **read.html:** Document reader view supporting native browser PDF embedding and ePub.js with a direct download button.
+ **settings.html:** Server status, library paths, and user profiles.

### web/static

Holds static client assets that are served directly to browsers under /static/:

+ **css:** Minimal, responsive stylesheet styled with a soft dark slate background and lavender accents.
+ **js:** Modular vanilla JavaScript modules. Includes api.js for backend communication, player.js for Plyr bindings and progress syncing, and reader.js for EPUB navigation.
+ **vendor:** Minimal embedded vendor libraries: Plyr for video and ePub.js for EPUB books.
+ **assets:** Mascots, icons, and fallback cover images.

All templates and static assets are embedded into the Go binary using embed.FS, meaning the server can be deployed as a single standalone executable.

### docs

Stores design specifications, complete database schemas and query procedures, multi-drive storage architecture, scraping engine designs, security threat models, TDD testing guidelines, architectural documentation, and user flow diagrams for the project.

### tests

Houses unit and integration tests, including tests for HTTP range handling, subtitle conversion, database queries, and media scanner behavior.
