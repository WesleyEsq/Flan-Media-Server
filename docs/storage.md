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
  `Short Film.mp4` (container title is derived from filename)

### Book Conventions Inside a Storage Source

* **Multi-Volume Series:**
  `Dune/Book 1.epub`
  `Dune/Book 2.epub`
  `Dune/poster.jpg`
* **Single Document:**
  `Computer Architecture.pdf`

---

## 3. Ingestion Methods: Direct Web Upload vs. Disk Sync

Flan provides two complementary ingestion workflows tailored to file size and user preference:

### 1. Direct Web Upload (In-Browser Convenience)
Designed for digital books (EPUB/PDF), single movies, and custom cover images without requiring SSH or network shares:
* **Direct-to-Disk Streaming:** Web uploads stream directly from the HTTP request body (`r.Body`) straight to the destination file on the target drive using `io.Copy`. 
* **Zero Root Flash Wear:** Uploaded video bytes **never** buffer in Go heap RAM and never write to temporary folders on the root micro-SD/eMMC (`/tmp`). This eliminates root flash wear and prevents crashes caused by filled boot drives.
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

### 1. Unmounted Storage Protection (`.flan-keep`)
* **Risk:** If an external drive fails to mount at boot, an empty mount folder could be mistaken for an empty library, leading to unintended database purges.
* **Mitigation:** Every storage source mandates a marker file named `.flan-keep`.
* **Zero-Terminal Setup:** When an administrator adds a new storage source in the web UI, Flan verifies `.flan-keep`. If missing, the UI offers an **[ Initialize .flan-keep ]** button, creating the file instantly without requiring terminal commands.

### 2. Isolated Source Failures
* If a secondary external drive unmounts or disconnects, only items linked to that specific `source_id` are temporarily flagged `is_missing = 1`. 
* Media on other active drives remains fully accessible, and server startup proceeds normally.

### 3. Drive Disconnection During Active Scan
* Before and during crawling of any storage source, Flan checks `.flan-keep`. If the marker disappears mid-scan, crawling of that source aborts immediately without altering database records.

### 4. Drive Thrashing Prevention (3-Stream Governor)
* Mechanical USB hard drives experience severe seek penalties under concurrent random reads. The playback lease governor restricts concurrent active video streams to 3, returning `HTTP 503` with a `Retry-After: 30` header when capacity is reached.

---

## 5. Related Documentation

* [Master System Architecture](design.md)
* [Database Schema & Queries](database.md)
* [Rate Limiting & Throttling](rate-limiting.md)
* [Security Threat Model](threat-model.md)
