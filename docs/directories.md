# Directories in the Project

The project is organized according to standard Go conventions. This keeps application entry points, internal packages, web assets, and documentation cleanly separated.

## Directory Structure

```
Flan-Media-Server/
├── cmd/
│   └── flan/               # Server entry point (main.go)
├── internal/               # Private application packages
│   ├── config/             # .env and environment variable loading
│   ├── database/           # SQLite connection, schema, and queries
│   ├── handler/            # HTTP handlers for auth, pages, streaming, and API
│   └── scraper/            # Media directory scanner and cover scraper
├── web/                    # Embedded web assets (via embed.FS)
│   ├── templates/          # Server-rendered Go HTML templates
│   └── static/             # Static web files served at /static/
│       ├── css/            # Stylesheets (soft dark theme with lavender accents)
│       ├── js/             # Modular vanilla JavaScript (API, player, reader)
│       └── assets/         # Mascots, icons, and default cover art
├── docs/                   # Specifications, architecture, and diagrams
│   ├── design.md           # Full system specification and database schema
│   ├── threat-model.md     # Security posture, attack vectors, and mitigations
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
Handles reading configuration settings from the .env file and system environment variables. Provides typed configuration values with sensible defaults for server port, host, media directory path, database file location, and memory tuning parameters.

### internal/database
Manages the SQLite3 database connection. Applies performance pragmas such as WAL mode, normal synchronous disk writes, and a 2mb page cache limit. Handles creating database tables (libraries, users, sessions, media_items, playback_progress) and indexes, and provides query functions for authentication, media indexing, and watch progress.

### internal/handler
Implements the HTTP request handlers and routing logic:
+ **Auth & Setup Handlers:** Manages first-time setup (/setup), user profile PIN authentication (/login), session cookies, rate limiting, and emergency recovery.
+ **Page Handlers:** Renders the Go templates (catalog, watch, read, and settings).
+ **Streaming Handler:** Handles video and document streaming using http.ServeContent. On Linux, this uses the sendfile system call to transfer byte ranges directly from the kernel page cache to the client socket without copying data into Go heap memory.
+ **Upload Handler:** Implements zero-memory streaming file and folder uploads via r.MultipartReader and io.Copy with pre-upload disk space verification.
+ **REST API Handlers:** Provides JSON endpoints for querying catalog items, retrieving playback positions, saving watch progress, and triggering library scans.

### internal/scraper
Scans configured media folders on disk for supported video formats (MP4, WebM, AV1) and book formats (EPUB, PDF). It extracts file details and can fetch cover artwork and descriptions from online sources over HTTPS.

### web/templates
Contains the Go HTML templates used to render web pages:
+ **base.html:** Common layout template containing the HTML shell, navigation bar, logo, and footer.
+ **setup.html:** First-time onboarding wizard to create the admin account and register the initial library folder.
+ **login.html:** Profile selector ("Who is watching?") and numeric PIN keypad.
+ **catalog.html:** Media grid view showing videos and books with filter tabs.
+ **watch.html:** Video player view with the HTML5 video element and resume playback controls.
+ **read.html:** Document reader view for EPUB and PDF books.
+ **settings.html:** Server status, library management, and user profiles.

### web/static
Holds static client assets that are served directly to browsers under /static/:
+ **css:** Responsive stylesheet styled with a soft dark slate background and lavender accents.
+ **js:** Modular vanilla JavaScript modules. Includes api.js for backend communication, player.js for video tracking and auto-saving progress, and reader.js for document navigation.
+ **assets:** Mascots, icons, and fallback cover images.

All templates and static assets are embedded into the Go binary using embed.FS, meaning the server can be deployed as a single standalone executable.

### docs
Stores design specifications, security threat models, architectural documentation, and diagrams for the project.

### tests
Houses unit and integration tests, including tests for HTTP range handling, database queries, and media scanner behavior.
