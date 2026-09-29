# Flan Media Server

A lightweight personal media server for streaming videos and reading books, built specifically for resource-constrained single-board computers (Raspberry Pi Zero, 1–5, Orange Pi, Rock Pi) and low-power Linux systems.

Flan delivers direct streaming without on-the-fly transcoding, keeping its operational footprint strictly around **15 to 20 MB of RAM** for 1 to 5 household users.

---

## At a Glance

| Feature | Specification |
| :--- | :--- |
| **Target Memory** | ~15–20 MB RAM (`GOMEMLIMIT=16MiB`, `GOGC=30`) |
| **Concurrent Streams**| Max 3 active video streams (governed by semaphore) |
| **Video Delivery** | Linux kernel `sendfile` zero-copy transfer (HTTP 206 Range requests) |
| **Supported Video** | Direct-play web formats: MP4 (H.264/AAC), WebM (VP9/Opus/AV1), web-safe MKV |
| **Reading Material**| Books & documents: EPUB (in-browser ePub.js) and PDF (native browser viewer) |
| **Database** | Pure-Go SQLite3 (`modernc.org/sqlite`) running in WAL mode (~2 MB page cache) |
| **Web Client** | Server-rendered Go `html/template` + vanilla JS & CSS embedded via `embed.FS` |
| **Theme** | Soft dark slate (`#12131a` / `#1a1b24`) with lavender accents (`#bb9af7`) |
| **Metadata** | Local-first (`poster.jpg` / embedded EPUB art) with optional TMDB fallback |

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

* **First Run:** Automatically opens the onboarding wizard at `/setup` to create an admin profile and register your first media library.
* **Returning Users:** Presents a "Who is watching?" profile selector with numeric PIN authentication.
* **CLI Account Recovery:** Reset a forgotten admin PIN directly from the terminal with `./flan --reset-admin`.

For cross-compiling to Raspberry Pi boards (ARMv6, ARMv7, ARM64) and systemd deployment, see the [Deployment Guide](docs/compilation.md).

---

## Core Principles

1. **Zero-Copy Streaming:** Video delivery uses Go's `http.ServeContent` and Linux `sendfile`. Bytes travel straight from the filesystem cache to the network socket, bypassing the Go garbage-collected heap.
2. **No Transcoding:** Avoids CPU saturation on low-power ARM cores. Files must be pre-encoded in web-compatible formats.
3. **Multi-Directory Libraries:** Dynamic libraries table in SQLite lets you map separate storage drives (e.g. fast NVMe for database/covers, spinning USB HDDs for bulk movies and TV shows).
4. **Resilient to Disconnections:** External drives can spin down when idle or unplug safely without wiping catalog records in the database.
5. **Completely Self-Contained:** Static assets, icons, and player libraries (Plyr, ePub.js) are baked into the single binary. Operates 100% offline with zero CDN dependencies.

---

## Documentation

* **Architecture & System Design:**
  * [Master System Architecture](docs/design.md)
  * [Directory Structure & Package Anatomy](docs/directories.md)
  * [Database Schema & Wear-Leveling](docs/database.md)
  * [Storage Architecture & Multi-Drive Resilience](docs/storage.md)
  * [Local-First Scraping Engine](docs/scraper.md)
* **Security & Traffic Control:**
  * [Five-Zone Rate Limiting](docs/rate-limiting.md)
  * [Security Threat Model & Mitigations](docs/threat-model.md)
* **Operations & Engineering:**
  * [Cross-Compilation & SBC Deployment Guide](docs/compilation.md)
  * [Testing Strategy & TDD Guidelines](docs/testing.md)
* **Web Client & UX:**
  * [Design System & Foundations](docs/client/design-system.md)
  * [Component Specifications](docs/client/components.md)
  * [Page Templates & Wireframes](docs/client/pages.md)
* **Diagrams:**
  * [Data Flow & Sequence Diagrams](docs/diagrams/data-flow.md)
  * [User Journeys & Technical Flows](docs/diagrams/user-flows.md)

---

## License

This project is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.
