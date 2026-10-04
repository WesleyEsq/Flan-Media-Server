# Flan Media Server

Flan is a self-hosted personal media server written in Go for streaming video files and reading digital books. It is designed to run efficiently on low-resource hardware—such as single-board computers (Raspberry Pi, Orange Pi) and repurposed personal computers—with a low memory footprint and direct file delivery.

---

## Key Characteristics

* **Direct Streaming (No Transcoding):** Video delivery uses standard HTTP 206 Range requests delegating to the Linux `sendfile` system call for zero-copy file transfers. Transcoding is omitted to avoid heavy CPU and disk load on low-power devices.
* **Local-First Metadata:** Operates completely offline with zero external metadata APIs (no TMDB or third-party service dependencies). Directory names define container titles, and local `poster.jpg` files provide cover art.
* **Embedded Storage:** Uses an embedded SQLite database (`modernc.org/sqlite`, CGO-free) operating in WAL mode with a bounded memory cache (~2 MB).
* **Single Static Binary:** The server compiles into a standalone binary with all web templates, stylesheets, scripts, and default avatars embedded via Go's `embed.FS`.
* **Layered Architecture:** Organized into a clean Controller-Service-Repository-Model structure with constructor dependency injection.
* **Authentication & Recovery:** User accounts are protected by 4 to 6-digit numeric PINs with brute-force lockout safeguards. Forgotten administrator credentials can be reset from the host terminal via a CLI flag.

---

## Quickstart

### Prerequisites

* Go 1.24 or newer (for building from source)
* Linux host (x86_64, ARM64, or ARMv7)

### Build and Run

```bash
# 1. Clone the repository
git clone https://github.com/WesleyEsq/Flan-Media-Server.git
cd Flan-Media-Server

# 2. Build standalone binary
go build -o flan ./cmd/flan

# 3. Start the server
./flan
```

By default, the server listens on `http://localhost:4907`. Configuration can be overridden using a `.env` file or environment variables.

### Initial Setup and Account Recovery

* **First Boot:** On initial startup with an empty database, the server generates a one-time 6-character bootstrap setup token and prints it to the terminal/journal. Visit `http://localhost:4907/setup` and enter the token to configure the initial administrator account.
* **Admin PIN Reset:** If the administrator PIN is lost, run `./flan --reset-admin` directly on the host to interactively set a new administrator PIN.

For cross-compiling to ARM targets or setting up a systemd service, see [docs/compilation.md](docs/compilation.md).

---

## Media Storage & Directory Layout

Flan supports configurable storage sources across multiple drives, initialized with default folders under `./media` (configurable via `MEDIA_DIR`):

```text
media/
├── video/
│   ├── Breaking Bad/
│   │   ├── S01E01.mp4
│   │   ├── S01E02.mp4
│   │   └── poster.jpg
│   └── Blade Runner (1982)/
│       ├── Final Cut.mp4
│       └── poster.jpg
└── books/
    ├── Dune/
    │   ├── Book 1.epub
    │   └── poster.jpg
    └── Operating Systems.pdf
```

To protect against unmounted drives or silent filesystem disconnections, Flan relies on `.flan-keep` marker files:
* **Primary App Tier (`./data/.flan-keep`):** The server halts startup immediately if missing, preventing writes to an unmounted root flash partition.
* **Storage Sources Tier (`<source_path>/.flan-keep`):** If an external drive disconnects, Flan safely flags that source's files as missing and keeps the server running without dropping catalog entries.

---

## Documentation

Comprehensive technical documentation is maintained in the `docs/` directory. See the [Documentation Index](docs/README.md) for full details:

* **Architecture & Backend:**
  * [Master System Architecture](docs/design.md): System constraints, MVC layered design, and HTTP routes.
  * [Directory Structure & Architecture](docs/directories.md): Package layout, responsibilities, and Java-to-Go concept mapping.
  * [Database Schema](docs/database.md): 7-table SQLite schema, WAL mode pragmas, and migrations.
  * [Storage Architecture](docs/storage.md): Drive decoupling, mount safety, and media directory conventions.
  * [Local Metadata Engine](docs/scraper.md): Filesystem scanner and database reconciliation logic.
* **Security & Traffic:**
  * [Authentication Security & Rate Limiting](docs/rate-limiting.md): PIN lockout safeguards, bcrypt throttling, and brute-force defenses.
  * [Threat Model & Mitigations](docs/threat-model.md): Security analysis, cookie authentication, and CSRF protection.
* **Operations:**
  * [Compilation & Deployment](docs/compilation.md): Cross-compilation commands and systemd unit configuration.
  * [Testing Strategy](docs/testing.md): Unit testing guidelines, in-memory SQLite, and virtual filesystems.
* **Web Client & Interface:**
  * [Design System](docs/client/design-system.md): Layout grid, color palette, and CSS foundations.
  * [Component Specifications](docs/client/components.md): Modal dialogs, card components, and form controls.
  * [Page Templates](docs/client/pages.md): Structure and wireframes for all 9 application views.
  * [Accessibility Guide](docs/client/accessibility.md): WCAG 2.1 AA requirements and focus management.
  * [Mobile Responsiveness](docs/client/responsiveness.md): Breakpoint specifications and mobile bottom navigation.
* **Diagrams:**
  * [Data Flow & Architecture](docs/diagrams/data-flow.md): Sequence diagrams for requests, streaming, and scanning.
  * [User Journeys](docs/diagrams/user-flows.md): Interaction flows for setup, login, and playback.

---

## License

This project is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.
