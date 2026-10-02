# Flan Media Server

An ultra-simple, lightweight personal media server for streaming videos and reading books. Designed specifically for low-power single-board computers (Raspberry Pi Zero, 1–5, Orange Pi, Rock Pi) and Linux homelabs.

Flan eliminates complexity: 100% offline, zero external APIs, fixed storage directories, a 6-table database, and a tactile high-contrast UI with a **sidebar-only navigation system**. It operates strictly within a **15 to 20 MB RAM** footprint for 1 to 5 household users.

---

## At a Glance

| Feature | Specification |
| :--- | :--- |
| **Target Memory** | ~15–20 MB RAM (`GOMEMLIMIT=16MiB`, `GOGC=30`) |
| **Navigation System** | **Sidebar-Only (Desktop) & Bottom Bar (Mobile):** All platform navigation is consolidated strictly into the left sidebar (`Video`, `Books`, `Manage Server`), shifting to a thumb-accessible fixed bottom bar on mobile screens. The top header contains only the brand title, `(?)` manual, and user profile avatar. |
| **Interface Style** | High-contrast neo-tactile layout with thick black borders, split Start screen, purple card footers, zero emojis, and pure CSS mobile responsiveness |
| **Media Architecture** | **Unified Containers (100% Local-First):**<br>• **Video (`./media/video`):** Series with episode lists, or movies with version lists<br>• **Books (`./media/books`):** Multi-volume series or single books (EPUB/PDF) |
| **Metadata Engine** | Zero-network, 100% offline (Folder name = title; `poster.jpg` = cover; no TMDB dependencies) |
| **Concurrent Streams**| Max 3 active video streams (governed by semaphore) |
| **Video Delivery** | Linux kernel `sendfile` zero-copy transfer (HTTP 206 Range requests) + 1-click VLC / Direct Download fallback for unsupported audio codecs (AC3/DTS) |
| **Supported Formats** | Direct-play: MP4 (H.264/AAC), WebM (VP9/Opus/AV1), web-safe MKV; Books: EPUB, PDF (direct download & browser view) |
| **Avatars** | Whimsical preset SVG companions (Flan mascot, cat, robot, ghost, etc.) or custom user photo uploads (max 2MB) |
| **Database** | Pure-Go SQLite3 (`modernc.org/sqlite`) running in WAL mode (~2 MB cache, 6 tables) |
| **Web Client** | Server-rendered Go `html/template` + vanilla JS/CSS embedded via `embed.FS` |

---

## Quickstart

### Build & Run Locally

```bash
# Compile single standalone binary
go build -o flan ./cmd/flan

# Run with tuned memory bounds
GOMEMLIMIT=16MiB GOGC=30 ./flan
```

Default access is at `http://localhost:4907` (configurable via `.env`).

* **Start / Login Screen:** Split screen with a "Welcome" graphic on the left and a dropdown user selector + numeric PIN field + `[ Access ]` button on the right.
* **First Run:** If no users exist, automatically prompts to create the initial admin account.
* **CLI Account Recovery:** Reset a forgotten admin PIN directly from the host terminal with `./flan --reset-admin`.

For cross-compiling to Raspberry Pi boards (ARMv6, ARMv7, ARM64) and systemd deployment, see the [Deployment Guide](docs/compilation.md).

---

## Core Principles

1. **Sidebar-Only Platform Navigation:** Every page inside the platform uses a unified left sidebar containing only **`Video`**, **`Books`**, and **`Manage Server`** (at the bottom).
2. **Unified Container + List Model:** Everything is either a **Video** (series with episodes, or movie with cut versions) or a **Book** (series with volumes, or single title).
3. **100% Offline & Zero-Network:** No external metadata APIs (no TMDB), no API keys, no network timeouts. Folder name is the title; `poster.jpg` is the cover.
4. **Fixed Storage Paths:** Fixed directories at `./media/video` and `./media/books`. No complex `libraries` database table or dynamic mount management.
5. **High-Contrast Tactile UI:** Thick 2px black borders, purple footer card bands, instant button clicks (`translateY(2px)`), zero emojis, and zero hover float delays.
6. **Zero-Copy Streaming:** Video delivery uses Go's `http.ServeContent` and Linux `sendfile`. Bytes travel directly from the filesystem cache to the network socket, bypassing the Go heap.

---

## Documentation

* **Architecture & System Design:**
  * [Master System Architecture](docs/design.md)
  * [Directory Structure & Package Anatomy](docs/directories.md)
  * [Streamlined 6-Table Database Schema](docs/database.md)
  * [Storage Architecture & Resilience](docs/storage.md)
  * [Local-First Metadata Engine](docs/scraper.md)
* **Security & Traffic Control:**
  * [Two-Safeguard Rate Limiting](docs/rate-limiting.md)
  * [Security Threat Model & Mitigations](docs/threat-model.md)
* **Operations & Engineering:**
  * [Cross-Compilation & SBC Deployment Guide](docs/compilation.md)
  * [Testing Strategy & TDD Guidelines](docs/testing.md)
* **Web Client & UX:**
  * [Design System & High-Contrast Foundations](docs/client/design-system.md)
  * [Tactile Component Specifications](docs/client/components.md)
  * [Page Templates & Wireframes (9 Templates)](docs/client/pages.md)
  * [Mobile Responsiveness & Adaptive Navigation](docs/responsiveness.md)
* **Diagrams:**
  * [Data Flow & Architecture](docs/diagrams/data-flow.md)
  * [User Journeys & Technical Flows](docs/diagrams/user-flows.md)

---

## License

This project is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.
