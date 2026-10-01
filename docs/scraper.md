# Local-First Metadata Engine

Flan Media Server adopts a **strictly 100% offline, zero-overhead metadata strategy**. External APIs (such as TMDB) have been completely removed, eliminating API keys, outbound HTTP requests, network timeouts, and background worker queues.

---

## 1. Core Principles

* **100% Local-First by Default:** Titles and covers come straight from your folder structure:
  * **Container Title:** Inferred from the folder name (e.g., `Breaking Bad` or `Blade Runner (1982)`).
  * **Playable Items:** Inferred from media filenames inside the folder (e.g., `S01E01 - Pilot.mp4` or `Theatrical Cut.mp4`), sorted naturally.
  * **Cover Artwork:** Looked up locally in the container folder as `poster.jpg` or `cover.jpg`, or extracted from embedded EPUB metadata.
* **Offline Autonomy:** Operates completely offline with zero network latency. If no local image exists, Flan displays an embedded high-contrast SVG mascot.
* **Zero Background Scraping Queues:** No background workers, retry loops, or upstream rate-limit timers. Crawling a library is a simple, synchronous filesystem walk that inserts records directly into SQLite.

---

## 2. Ingestion Pipeline

```text
[Filesystem Directory under ./media/video or ./media/books]
       │
       ▼
1. Walk Directory Tree ──────► Discovers container folder and media files
       │
       ▼
2. Local Art Check ──────────► Looks for poster.jpg / cover.jpg / embedded EPUB art
       │
       ▼
3. Clean Title & Insert ─────► Container title = folder name; files = sorted items
       │
       ▼
4. Database Persistence ─────► Inserts into `videos` & `video_files` (or `books` & `book_files`)
```

---

## 3. Supported Folder Conventions

| Media Type | Recommended Folder Structure | Example |
| :--- | :--- | :--- |
| **Video (Series or Movie)** | `./media/video/<Title>/<file1>.<ext>`<br>`./media/video/<Title>/poster.jpg`<br>or flat `./media/video/<Title>.<ext>` | `video/Breaking Bad/S01E01.mp4`<br>`video/Breaking Bad/poster.jpg`<br>`video/Spirited Away (2001).mp4` |
| **Books** | `./media/books/<Title>/<file1>.<ext>`<br>or flat `./media/books/<Title>.<ext>` | `books/Dune/Book 1.epub`<br>`books/Linux Kernel.pdf` |

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Storage Architecture](storage.md)
* [Simplified Database Schema (6 Tables)](database.md)
