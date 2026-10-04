# Storage Architecture & Drive Resiliency

Flan Media Server uses a fixed, predictable directory layout designed for homelab environments where storage is separated between fast primary storage and bulk media drives.

---

## 1. Storage Tiers & Drive Decoupling

```text
Host System (Repurposed PC / SBC)
 ├── Boot & App Tier (Fast Flash / NVMe SSD / eMMC)
 │     ├── flan binary
 │     ├── data/flan.db (SQLite WAL with 2 MB cache)
 │     ├── data/covers/ (Cached artwork)
 │     └── data/avatars/ (Uploaded custom user avatars)
 │     → Fast random I/O; catalog browsing does not access bulk drives.
 │
 └── Bulk Media Tier (External USB 3.0 / SATA HDD)
       ├── ./media/video/ (Video containers: movies and series)
       └── ./media/books/ (Book containers: EPUB and PDF)
       → Spinning drives sleep when idle; wake only when playback begins.
```

---

## 2. Configurable Storage Sources & Directory Layout

Rather than restricting the entire server to a single hardcoded path, Flan supports **multiple configurable storage sources** managed directly from the `/manage` web console. Each source maps to a specific media type:

* **Video Sources:** Can span multiple physical disks (e.g. `./media/video`, `/mnt/usb-movies`, `/mnt/usb-tv`).
* **Book Sources:** Can map to local or mounted directories (e.g. `./media/books`, `/mnt/storage/ebooks`).

### Video Conventions Inside a Storage Source

Videos are organized as containers containing files. A container represents a movie or television series:

* **Series with Episodes (with nested season support):**
  `Breaking Bad/S01E01.mp4` OR `Breaking Bad/Season 1/S01E01.mp4`
  `Breaking Bad/poster.jpg` (optional container cover)
* **Movie with Cuts / Versions:**
  `Blade Runner (1982)/Theatrical Cut.mp4`
  `Blade Runner (1982)/Final Cut.mp4`
  `Blade Runner (1982)/poster.jpg`
* **Flat File (Single Video):**
  `Short Film.mp4` (in the database, `videos.folder_path` is set to `"Short Film.mp4"` and `video_files.relative_path = ""` to guarantee uniqueness across multiple flat files under `UNIQUE(source_id, folder_path)`)

### Book Conventions Inside a Storage Source

* **Multi-Volume Series:**
  `Dune/Book 1.epub`
  `Dune/Book 2.epub`
  `Dune/poster.jpg`
* **Single Document:**
  `Computer Architecture.pdf` (in the database, `books.folder_path` is set to `"Computer Architecture.pdf"` and `book_files.relative_path = ""`)

---

## 3. Ingestion Methods: Direct Web Upload vs. Disk Sync

Flan provides two complementary ingestion workflows tailored to file size and user preference:

### 1. Direct Web Upload (In-Browser Convenience)
Designed for digital books (EPUB/PDF), single movies, and custom cover images without requiring SSH or network shares:
* **Direct-to-Disk Streaming via `r.MultipartReader()`:** Web uploads parse parts incrementally using `r.MultipartReader()` and stream directly to destination storage using `io.Copy()`. This explicitly bypasses standard `r.ParseMultipartForm()`, which dumps parts >10 MB to root flash `/tmp`.
* **Zero Root Flash Wear & Atomic Staging:** Uploaded bytes stream directly to `<target_folder>/<filename>.part` on the bulk media drive. Upon successful complete transfer, `os.Rename()` atomically commits the file to its final name. If the transfer is cancelled or disconnected (`<-r.Context().Done()`), the partial `.part` file is immediately deleted via `os.Remove()`, ensuring incomplete uploads never pollute the catalog.
* **Disk Space & Traversal Protection:**
  * Before accepting bytes, Flan inspects free space on the destination drive using `syscall.Statfs()` (using explicit `uint64` casting to avoid 32-bit integer overflow on drives $\ge 2\text{TB}$ on ARMv7/ARMv6 SBCs). If free space is below 1 GB, the upload is rejected with `HTTP 507 Insufficient Storage`.
  * Filenames are strictly sanitized with `filepath.Base()`, stripped of control characters, and checked against allowed media extensions (`.mp4`, `.mkv`, `.webm`, `.epub`, `.pdf`, `.jpg`, `.png`).
* **HTML5 Progress Reporting:** The web interface provides real-time upload progress (percentage, speed, and time remaining).
* **Automatic Metadata Extraction:**
  * Books (EPUB): Embedded title, author, and cover art are extracted in pure Go upon upload.
  * Videos: Prompts the user for container title and release year during upload.

### 2. Disk Synchronization (Bulk Homelab Loading)
For large video libraries and multi-season television shows, administrators manage files directly using standard host tools:
* **Samba / NFS File Shares:** Mount media drives as network shares on desktop workstations.
* **Terminal Synchronization (rsync / scp):** Sync media directly over SSH.
* **Direct USB Drives:** Plug external USB drives directly into the server and register their mount paths as storage sources in `/manage`.
* **Supervised Ingestion:** After copying files, click `[ Rescan All Media ]` in `/manage`. Discovered items are cleaned and staged in the Ingestion Queue for review.

---

## 4. Drive Disconnection & Failure Mitigations

### 1. Two-Tier Marker Protection (`.flan-keep`)
* **Primary App Tier (`./data/.flan-keep`):** Verifies the critical application storage partition. If missing at startup, Flan halts immediately to prevent writing a blank SQLite database or session secret to an unmounted root mount point.
* **Storage Sources Tier (`<source_path>/.flan-keep`):** Mandated inside each configured storage source directory to prevent an unmounted media partition from being mistaken for an emptied library.
* **Device Boundary Safety on Initialization:** When initializing `.flan-keep` via the web console for a newly added source, Flan compares the filesystem device identifier (`syscall.Stat_t.Dev`) of the folder against the root device (`/`). If a mount directory (e.g. `/mnt/usb-movies`) shares the root device ID, Flan flags a warning confirming the external partition is truly mounted before creating the marker.

### 2. Isolated Source Failures & Startup Resiliency
* If an external media drive fails to mount at boot or disconnects:
  * Server startup **proceeds normally** without halting.
  * The disconnected storage source is flagged as offline in memory and the `/manage` console.
  * Catalog items belonging to that `source_id` are temporarily treated as `is_missing = 1` rather than deleted.
  * Media on other active storage drives remains fully accessible.
  * Once the drive is remounted and verified, items return to normal active status.

### 3. Drive Disconnection During Active Scan
* Before and during crawling of any storage source, Flan checks for `.flan-keep`. If the marker disappears mid-scan, crawling of that source aborts immediately without altering database records or purging historical watch progress.

---

## 5. Related Documentation

* [Master System Architecture](design.md)
* [Database Schema & Queries](database.md)
* [Authentication Security & Rate Limiting](rate-limiting.md)
* [Security Threat Model](threat-model.md)
