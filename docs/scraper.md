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
   * Files are uniquely identified by `(video_id, relative_path)` or `(book_id, relative_path)`.
   * Titles on disk can differ completely from the catalog display title without issue.

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
