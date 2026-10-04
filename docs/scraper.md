# Local-First Metadata Engine

Flan Media Server uses an offline, directory-based metadata strategy. External metadata APIs (such as TMDB) are completely omitted, eliminating external network dependencies, API keys, and outbound network traffic.

---

## 1. Operating Principles

* **Decoupled Architecture:** Physical filesystem paths are treated strictly as immutable byte locators (`folder_path` and `relative_path`). Display metadata (`title`, `release_year`, `overview`, episode names, `order_index`) is stored independently in SQLite.
* **Intelligent Filename Cleaning:** When discovering unformatted torrent or scene releases, the crawler automatically strips release noise (`1080p`, `x264`, `BluRay`, `WEBRip`, `[eztv]`) and detects season/episode numbers (`S01E02` $\to$ Season 1, Episode 2) to propose human-readable defaults.
* **Natural Alphanumeric Sorting:** Files within a container are sorted using human natural ordering (`Episode 1`, `Episode 2`, ..., `Episode 10`) rather than standard ASCII lexicographical sort.
* **On-Demand Ingestion Pipeline:** Ingestion is triggered explicitly when adding a storage source or clicking `[ Scan ]`. Discovered items are presented in an on-demand review wizard where administrators can check/uncheck candidates, fix typos inline, expand episode lists, and commit only chosen items to SQLite.
* **Clean Database Guarantee:** The database is never polluted with unapproved or rejected candidate rows. Only items explicitly selected and ingested are written to SQLite.
* **Preservation of User Edits:** Once an item is ingested or manually edited, `metadata_locked = 1` is enforced. Future maintenance rescans will **never** overwrite user-edited titles, descriptions, episode titles, or order indices.

---

## 2. Ingestion Pipeline & Synchronization Algorithm

```text
User Triggers [ Scan ] on Storage Source
       │
       ▼
1. Marker Check ───────────► Verifies .flan-keep in source directory
       │
       ▼
2. Ephemeral Crawler Walk ─► Discovers container folders and media files (handles seasons)
       │
       ▼
3. Local Artwork Check ────► Identifies poster.jpg / cover.jpg / embedded EPUB art
       │
       ▼
4. Name Cleaner & Parse ───► Strips scene noise, detects SxxExx, applies natural sorting
       │
       ▼
5. Interactive Review ─────► Presents gathered items in Ingestion Pipeline Wizard (modal):
       │                     • Checkbox per item (uncheck sample clips or unwanted folders)
       │                     • Inline title & year editing to fix typos
       │                     • Expandable file & episode list preview
       │
       ▼
6. Ingest Selected ────────► POST /api/sources/{id}/ingest commits ONLY checked items to SQLite
       │                     with metadata_locked = 1
```

### Synchronization & Presentation Rules

1. **Identity Resolution:**
   * Containers are uniquely identified by `(source_id, folder_path)`.
   * For folder containers (series or movies in dedicated directories), `folder_path` is the container directory name (e.g. `Breaking Bad`), and `relative_path` is the filename within it (e.g. `S01E01.mp4`).
   * For flat files directly in the storage source root, `folder_path` is the filename itself (e.g. `Documentary.mp4`), and `video_files.relative_path = ""`. This guarantees that `UNIQUE(source_id, folder_path)` is never violated when multiple flat files exist in the same source.
   * Display titles in SQLite can differ completely from disk paths without affecting file discovery.

2. **Fixing Typos in the UI (Zero Disk Pain):**
   * If a folder on disk is named `spirted.away.2001.1080p`, the administrator can fix the title to `Spirited Away (2001)` in the web UI.
   * This updates the database record directly and sets `metadata_locked = 1`.
   * The physical folder on disk is **never** renamed, preventing file lock errors and preserving active torrent seeds or backup sync tools.
   * During subsequent rescans, the crawler verifies that the files still exist via `folder_path`, but leaves the custom title untouched.

3. **Intelligent Filename Cleaning Regex:**
   * **Container Cleaners:** Removes bracketed groups (`[...]`, `(-...)`), release years, resolution tags (`720p`, `1080p`, `4k`), audio tags (`DTS`, `AAC`, `5.1`), and replaces dots/underscores with clean spaces.
   * **Episode Cleaners:** Extracts `(?i)s(\d+)e(\d+)` or `(\d+)x(\d+)` to assign `order_index` and format clean episode labels (`S01E01 - Episode 1`).

4. **Missing vs. Deleted Files:**
   * If a file disappears from disk while `.flan-keep` is verified, it is marked `is_missing = 1` rather than deleted immediately.
   * If the file reappears in a subsequent scan, `is_missing` is reset to `0`.
   * Administrators can permanently purge missing records from the `/manage` console.

5. **New Files in Existing Containers & Reconciliation:**
   * **Container Lock Scope:** `videos.metadata_locked = 1` locks container-level fields (`title`, `release_year`, `overview`, `cover_path`, `video_type`). Maintenance rescans (`POST /api/scan`) and storage source scans will **never** overwrite these fields with raw disk directory names.
   * **Incremental Episode Discovery:** When a user drops a new file (e.g. `S01E03.mp4`) into an existing container directory, rescanning detects the new file and inserts a new `video_files` record with clean defaults, associating it with the existing `video_id`.
   * **Preservation of Existing Files:** Previously indexed files retain their `custom_title`, `order_index`, and user playback progress.

6. **Stateless Ingestion Pipeline:**
   * When `POST /api/sources/{id}/scan` runs, candidate items are returned directly as JSON in the HTTP response to the browser modal.
   * The browser modal holds candidate state while the administrator reviews, checks, or edits titles.
   * Ingestion (`POST /api/sources/{id}/ingest`) accepts the approved items array directly in the POST body: `{ "approved": [{ "folder_path": "the.wire.s01", "title": "The Wire", "year": 2002, "files": [...] }] }`.
   * The server asserts each `folder_path` resides strictly within the verified source directory (`filepath.Clean`), committing approved items with `metadata_locked = 1`. This keeps the backend 100% stateless without ephemeral server-side UUID tokens or memory cache TTLs.

7. **Video Type Classification & Season Ordering:**
   * **Movie vs. Series Classification:** A container is classified as `movie` if it contains a single video file without season/episode tags; it is classified as `series` if it contains multiple files or detects `SxxExx` patterns. This is stored in `videos.video_type` and can be manually adjusted in the video edit modal.
   * **Season & Episode Extraction:** Regex extracts `(?i)s(\d+)e(\d+)` to populate `season_number` and `episode_number`. The file's `order_index` is deterministically computed as:
     $$\text{order\_index} = (\text{season\_number} \times 1000) + \text{episode\_number}$$
     This prevents episode number collisions across multi-season series (e.g. S01E01 has index 1001, S02E01 has index 2001).

8. **Subtitles & Duration Extraction (Zero External Tools):**
   * **Sidecar Subtitle Discovery:** Scans container directories for subtitle sidecars sharing the base filename (e.g. `S01E01.en.srt`, `S01E01.vtt`). Served dynamically via `GET /stream/subtitles/{file_id}/{track_id}` with pure-Go SRT-to-WebVTT conversion.
   * **Zero-ffprobe Duration Discovery:**
     * **MP4 Files:** Duration is parsed directly from the MP4 `mvhd` box in pure Go using `io.ReadSeeker` to skip box payloads without buffering media bytes.
     * **MKV / WebM Files:** Duration is parsed directly from the Matroska EBML header (`Segment -> Info -> Duration`) in pure Go in under 1 millisecond. This ensures 100% offline, immediate duration discovery without relying on browser HTML5 `<video>` events (which fail on browsers lacking native MKV support).

---

## 3. Directory Conventions

### Video Media

* **Series with Episodes:**
  `./media/video/Breaking Bad/S01E01.mp4`
  `./media/video/Breaking Bad/poster.jpg`
* **Single Movie / Feature:**
  `./media/video/Blade Runner (1982)/Final Cut.mp4`
  `./media/video/Blade Runner (1982)/poster.jpg`
* **Flat File:**
  `./media/video/Documentary.mp4`

### Book Media

* **Multi-Volume Series:**
  `./media/books/Dune/Book 1.epub`
  `./media/books/Dune/poster.jpg`
* **Single Document:**
  `./media/books/Linux Kernel.pdf`

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Storage Architecture](storage.md)
* [Database Schema & Queries](database.md)
