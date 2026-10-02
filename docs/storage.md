# Storage Architecture & Drive Resiliency

Flan Media Server uses a fixed, predictable directory layout designed for single-board computers where storage is decoupled between fast flash storage and bulk mechanical drives.

---

## 1. Storage Tiers & Decoupling

```text
[SBC System]
 ├── Boot & App Tier (Micro-SD / eMMC / NVMe SSD)
 │     ├── flan executable
 │     ├── data/flan.db (SQLite WAL with 2 MB cache)
 │     ├── data/covers/ (Locally cached artwork)
 │     └── data/avatars/ (Uploaded custom user avatars)
 │     → Fast random I/O; catalog browsing never touches mechanical drives.
 │
 └── Bulk Media Tier (External USB 3.0 / SATA HDD)
       ├── ./media/video/ (All video containers: movies and series)
       └── ./media/books/ (All book containers: EPUB and PDF)
       → Spinning drives sleep when idle; wake only when video streams start.
```

---

## 2. Fixed Directory Layout (No `libraries` Table)

Rather than maintaining a dynamic `libraries` database table, Flan uses fixed, deterministic subdirectories under `MEDIA_DIR` (defaults to `./media`):

* `./media/video/`: Root directory for all video containers.
* `./media/books/`: Root directory for all book containers.

### Deterministic Folder Conventions

| Media Type | Folder Structure | Example |
| :--- | :--- | :--- |
| **Video** | `./media/video/<Container Title>/<file>.<ext>`<br>`./media/video/<Container Title>/poster.jpg`<br>or flat `./media/video/<Title>.<ext>` | `video/Breaking Bad/S01E01.mp4`<br>`video/Breaking Bad/poster.jpg`<br>`video/Blade Runner/Final Cut.mp4` |
| **Books** | `./media/books/<Container Title>/<file>.<ext>`<br>or flat `./media/books/<Title>.<ext>` | `books/Dune/Book 1.epub`<br>`books/Linux Kernel.pdf` |

---

## 3. Host Media Population & Storage Tiers

Flan eliminates the fragility of browser-based multi-gigabyte uploads by relying on proven, standard homelab media transfer methods. Administrators manage files directly on host storage:

1. **Local Network File Shares (SMB / NFS):** Mount `./media` as a Samba share on local desktop workstations (Windows Explorer / macOS Finder) and drag-and-drop video folders with full network resilience.
2. **Secure Copy / Terminal (SCP / rsync):** Sync media libraries directly over SSH (`rsync -avP /local/movies/ pi@flan:./media/video/`).
3. **Direct USB Drive Attachment:** Mount external USB 3.0 drives containing media directly into `./media/video` or `./media/books`.
4. **Immediate Catalog Indexing:** Once files are copied, clicking `[ ⟳ Rescan All Media ]` in the `/manage` console synchronously indexes all new titles into SQLite in sub-second times.

---

## 4. Drive Disconnection & Mount Failure Defenses

| Failure Scenario | Risk | Flan Defensive Mitigation |
| :--- | :--- | :--- |
| **Unmounted Drive on Boot** | An empty directory on root flash is mistaken for a new database, creating a split-brain "ghost" database. | **Marker File (`.flan-keep`):** Before initializing an external DB path, Flan checks for `.flan-keep`. If absent, startup halts with a fatal log message. |
| **Drive Disconnect During Scan** | Scanner assumes all files were deleted and purges database catalog rows. | **Mount Liveness Check:** Before scanning, Flan verifies that the directory exists and contains files. If absent or empty, the scan aborts safely, leaving database rows untouched. |
| **Streaming an Offline File** | Playback freezes or crashes the server handler. | **Soft-Offline Response:** The server returns `HTTP 503 Service Unavailable` (*"Media drive is offline"*). When reconnected, streaming resumes with zero re-indexing. |
| **Mechanical Drive Thrashing** | Multiple concurrent video streams cause mechanical head seeking, dropping throughput from 100 MB/s to 3 MB/s. | **Stream Governor Semaphore:** Caps active streams at 3, ensuring smooth sequential disk reads. |

---

## 5. Related Documentation

* [Master System Architecture](design.md)
* [Simplified Database Schema (6 Tables)](database.md)
* [Rate Limiting Architecture](rate-limiting.md)
* [Security Threat Model](threat-model.md)
