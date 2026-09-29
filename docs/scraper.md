# Scraping Engine & Metadata Pipeline

The scraping engine enriches raw media files with metadata, ratings, genre tags, and cover artwork while respecting Flan's 15–20 MB operational memory ceiling.

---

## 1. Core Principles

* **Local-First Priority:** Scans local folders for existing artwork (`poster.jpg`, `cover.jpg`, folder art, or embedded EPUB covers) before making any network requests.
* **Offline Resilience:** All cover images and metadata are stored locally on disk (`data/covers/`) and in SQLite. The server never hotlinks remote URLs. If `TMDB_API_KEY` is absent or offline, Flan uses embedded SVG placeholders with zero network delays.
* **Streamed Artwork Caching:** Remote images stream directly from the HTTPS connection to disk in 32 KB chunks via `io.Copy`. Image bytes never linger in Go heap memory.
* **Sequential Paced Worker:** Scraping runs in a single background goroutine metered at ~2.8 requests/second (350 ms ticker, Zone E) to prevent memory spikes and upstream IP bans.

---

## 2. Four-Stage Scraping Pipeline

```text
[Raw File on Disk]
       │
       ▼
Stage 1: Local Art Check ──(Found)──► Index immediately with local art (Zero Network)
       │ (Not found)
       ▼
Stage 2: Regex Filename Cleaner ──► Extracts Title, Year, Season, Episode tokens
       │
       ▼
Stage 3: External API Match ────► Throttled TMDB / OpenLibrary HTTPS query (10s timeout)
       │
       ▼
Stage 4: Streamed Cover Storage ─► Streams to data/covers/{type}/{id}.jpg & inserts SQLite rows
```

### Stage Details

1. **Local Art Inspection:** Checks directory for `poster.jpg`, `cover.jpg`, or extracts embedded EPUB zip covers (`cover.jpeg` / `OEBPS/images/cover.jpg`).
2. **Filename Sanitization:** Removes resolution tags (`1080p`, `4K`, `2160p`), sources (`BluRay`, `WEBRip`), codecs (`x264`, `x265`), and release group tags to parse clean title and year.
3. **External API Queries:** Queries TMDB API v3 via HTTPS for canonical synopsis, rating, genres, and poster URLs. If no API key is set, it logs an informational note and uses the clean filename with an embedded SVG hamster mascot.
4. **Cover Download & Indexing:** Streams poster bytes to disk, serves them locally at `/covers/{type}/{id}` with long-lived caching (`Cache-Control: public, max-age=31536000`), inserts normalized genres into `genres` and `item_genres`, and yields SQLite locks (`runtime.Gosched()`).

---

## 3. Supported File Naming Conventions

| Media Type | Recommended Folder Structure | Example |
| :--- | :--- | :--- |
| **Movies** | `<Library>/<Title> (<Year>)/<Title> (<Year>).mp4` (with optional `poster.jpg`)<br>or flat `<Library>/<Title> (<Year>).mp4` | `Movies/Dune (2021)/Dune (2021).mp4`<br>`Movies/Spirited Away (2001).mp4` |
| **TV Shows** | `<Library>/<Series>/Season <NN>/<Series> - S<NN>E<NN> - <Title>.<ext>` | `TV/Breaking Bad/Season 01/Breaking Bad - S01E01 - Pilot.mp4` |
| **Books** | `<Library>/<Author>/<Title>.<ext>` or flat `<Library>/<Title>.<ext>` | `Books/Frank Herbert/Dune.epub`<br>`Books/Linux Kernel Development.pdf` |

---

## 4. Admin Manual Override ("Fix Match")

When automatic scraping matches the wrong release or when custom metadata is desired, administrators can use the "Fix Match" modal on any media card:

* **Manual Search:** Search by custom title or paste a specific TMDB ID (e.g. `tmdb:12345`).
* **Custom Artwork Upload:** Drag-and-drop a local JPG, PNG, or WEBP image to replace the cover art immediately.
* **Direct Field Editing:** Manually edit display title, release year, overview synopsis, and genre tags.

---

## 5. Related Documentation

* [Master System Architecture](design.md)
* [Rate Limiting Architecture (Zone E)](rate-limiting.md)
* [Database Schema & Genres](database.md)
* [Ingestion Data Flow Diagram](diagrams/data-flow.md#4-ingestion--scraping-data-flow-uml-sequence)
