# Scraping Engine Specification

This document details the architecture, parsing pipeline, external metadata sources, and local storage strategy for Flan Media Server's scraping engine.

---

## 1. Overview and Objectives

The scraping engine is responsible for enriching raw video and document files with titles, summaries, ratings, genres, and cover artwork. To respect the strict 15 to 20mb operational memory ceiling of single-board computers, the scraper operates under the following constraints:

+ **Local-First Priority:** The scanner always checks the local folder for existing artwork (poster.jpg, cover.jpg, folder.jpg) or embedded EPUB covers before making any network requests. If local art exists, it is indexed immediately without calling external APIs.
+ **Offline Resilience:** All cover artwork, poster images, and metadata are saved locally on disk and in SQLite. The server never hotlinks remote image URLs in client templates, allowing the library to function seamlessly without internet access.
+ **Streamed Image Storage:** Downloaded covers are streamed directly from the remote HTTPS connection to disk in 32kb chunks via io.Copy. Image bytes never linger in Go application memory.
+ **Sequential Background Processing:** Scraping runs as a throttled, sequential background worker (2 to 3 requests per second) to prevent memory spikes and avoid triggering upstream HTTP 429 rate limits.
+ **Manual Override ("Fix Match"):** Administrators can manually search for matches, upload custom covers, or edit tags directly from the web interface.

---

## 2. The Four-Stage Scraping Pipeline

``` text
[Raw File on Disk]
       ↓
Stage 1: Local Art Check (poster.jpg / cover.jpg / embedded EPUB)
       ↓ (If not found locally)
Stage 2: Clean Filename Parsing (Regex Sanitizer)
       ↓
Stage 3: External API Query (TMDB for Video)
       ↓
Stage 4: Streamed Cover Download & Metadata Insertion (SQLite)
```

### Stage 1: Local Art Check (Zero-Network Path)

Before calling any external service, the scanner inspects the local directory:

+ For **Movies & TV:** Looks for poster.jpg, cover.jpg, or folder.jpg in the same directory as the media file.
+ For **EPUB Books:** Reads the internal EPUB zip directory to extract the embedded cover image (typically cover.jpeg or OEBPS/images/cover.jpg).
+ If local art is found, it is copied or linked to the local covers directory, and the file is indexed immediately with zero network latency.

### Stage 2: Clean Filename Parsing

If no local metadata exists, raw filenames are sanitized using regular expressions to remove release artifacts, resolution tags, and encoding information:

1. **Resolution & Source Stripping:** Removes 1080p, 720p, 4K, 2160p, WEBRip, BluRay, HDTV, and x264/x265 tags.
2. **Audio & Group Stripping:** Removes AAC, DTS, DDP5.1, and bracketed release group names.
3. **Token Extraction:**
   + **Movie Pattern:** Extracts Title and Year (e.g. "Dune Part Two" and "2024").
   + **TV Show Pattern:** Detects Show Title, Season, and Episode from standard patterns like S01E05 or 1x05.

### Stage 3: External API Matchers

When external metadata is needed, the scraper queries lightweight REST APIs over HTTPS using Go's standard library net/http:

+ **The Movie Database (TMDB):**
  + Endpoint: Search queries via TMDB REST API.
  + Retrieved Fields: Canonical title, release year, community rating (e.g. 8.2), overview synopsis, poster image path, and genre array (e.g. ["Animation", "Sci-Fi"]).
+ **Books:** For books lacking embedded covers, OpenLibrary is queried using the cleaned title.

### Stage 4: Streamed Cover Download & Metadata Storage

When a remote cover image URL is identified:

1. The server opens a destination file on local disk under the configured data directory: `data/covers/{type}/{id}.jpg`.
2. The remote image is fetched via an HTTPS GET request.
3. The response body is copied directly to disk using `io.Copy`.
4. The local file path is recorded in SQLite.
5. Genres are inserted into the normalized `genres` and `item_genres` tables, enabling fast index lookups and genre filtering without full-table string scanning.
6. Client browsers fetch covers from the local endpoint `/covers/{type}/{id}`, which serves files with long-lived browser caching headers (`Cache-Control: public, max-age=31536000`).

---

## 3. Manual Override and "Fix Match"

Automatic scrapers occasionally make mistakes, such as confusing a 1984 film with a 2020 remake sharing the same name. Flan provides an administrator-only modal on every media card:

+ **Manual Search:** Allows typing an alternative search title or directly pasting a TMDB ID (e.g. tmdb:12345) to force an exact re-fetch.
+ **Custom Image Upload:** Allows dragging and dropping a local JPG, PNG, or WEBP file to immediately replace the cover art.
+ **Direct Field Editing:** Allows manually adjusting the display title, release year, overview, and custom genres.

---

## 4. Supported File Naming Conventions

To ensure high match accuracy, media files should adhere to standard naming conventions within their respective libraries:

### Movies

``` text
Movies/
├── The Matrix (1999)/
│   ├── The Matrix (1999).mp4
│   └── poster.jpg                  # Optional local cover
└── Spirited Away (2001).mp4
```

### TV Shows

``` text
TV Shows/
└── Breaking Bad/
    ├── poster.jpg                  # Optional show cover
    ├── Season 01/
    │   ├── Breaking Bad - S01E01 - Pilot.mp4
    │   └── Breaking Bad - S01E02 - Cat's in the Bag.mp4
    └── Season 02/
        └── Breaking Bad - S02E01 - Seven Thirty-Seven.mp4
```

### Books

``` text
Books/
├── Frank Herbert/
│   └── Dune.epub                   # Embedded cover auto-extracted
└── Documentation/
    └── Go Programming Language.pdf
```

---

## 5. Resource Guardrails & Error Handling

To protect the server's 15 to 20mb memory target and prevent network lockups:

+ **Concurrency Limits:** Only one scraping task runs at any time in a single background goroutine.
+ **Request Throttling:** A ticker enforces a 350ms delay between consecutive API calls to stay within free API rate limits, governed by [Zone E of the Rate Limiting Architecture](docs/rate-limiting.md).
+ **Timeout Protection:** Every outbound HTTP request uses a strict 10-second timeout via context.WithTimeout.
+ **Graceful Fallbacks:** If an item cannot be matched online or the server is running without an internet connection:
  + The display title defaults to the cleaned filename.
  + The cover defaults to an embedded SVG placeholder or an embedded hamster mascot graphic.
  + The item remains fully playable and can be manually edited later.

---

### Related Documentation

+ [Master System Specifications](docs/design.md)
+ [Five-Zone Rate Limiting Architecture](docs/rate-limiting.md)
+ [Security Threat Model & SSRF Defense](docs/threat-model.md)
+ [Database Schema & Catalog Queries](docs/database.md)
+ [Ingestion & Scraping Sequence Diagram](docs/diagrams/data-flow.md#4-ingestion--scraping-data-flow-uml-sequence)
