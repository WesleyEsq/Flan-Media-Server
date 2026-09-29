# Storage Architecture & Drive Resiliency

Flan Media Server is designed for multi-tier single-board computer setups where storage is distributed across micro-SD cards, NVMe SSDs, and external USB spinning hard drives.

---

## 1. Storage Tiers & Decoupling

```text
[SBC System]
 ├── Boot & App Tier (Micro-SD / eMMC / NVMe SSD)
 │     ├── flan executable
 │     ├── data/flan.db (SQLite WAL with 2 MB cache)
 │     └── data/covers/ (Locally cached artwork)
 │     → Fast random I/O; catalog browsing never touches mechanical drives.
 │
 └── Bulk Media Tier (External USB 3.0 / SATA HDD)
       ├── /mnt/hdd1/movies (Standalone films)
       ├── /mnt/hdd2/tv     (TV series & seasons)
       └── /mnt/nvme/books  (EPUB & PDF books)
       → Spinning drives sleep when idle; wake only when video streams start.
```

---

## 2. Dynamic Library Management

Instead of hardcoding paths in `.env`, storage roots are managed dynamically in SQLite via the `libraries` table:

* **Graceful Defaults:** If `MEDIA_DIR` is blank, Flan boots cleanly and initializes `./media/movies`, `./media/tv`, and `./media/books` during the `/setup` wizard.
* **Runtime Management:** Admins can add or remove library roots dynamically via `/settings` or `POST /api/libraries` without restarting the server daemon.

---

## 3. Deterministic Folder Conventions

Because each library explicitly declares its `media_type` (`movies`, `tv`, `books`), the intake crawler works deterministically:

| Media Type | Accepted Folder Structure |
| :--- | :--- |
| **Movies** | `<library_path>/Movie Title (Year).mp4`<br>`<library_path>/Movie Title (Year)/Movie Title (Year).mp4` (with optional `poster.jpg`) |
| **TV Shows** | `<library_path>/<Series Title>/Season <NN>/<Series Title> - S<NN>E<NN> - <Title>.<ext>`<br>(Optional series poster: `<library_path>/<Series Title>/poster.jpg`) |
| **Books** | `<library_path>/<Author>/<Book Title>.<ext>` or `<library_path>/<Book Title>.<ext>` |

---

## 4. Admin Web Upload Pipeline

Administrators can upload media through the web interface without terminal access:

1. **Upload Routing:** The modal lets the admin pick the target library. For TV shows, it captures Series Title and Season Number, streaming directly into the appropriate season folder.
2. **Pre-Flight Disk Verification (`statfs`):** Verifies that the destination filesystem has **>2 GB free space**. If space is low, it rejects the upload immediately with `HTTP 507 Insufficient Storage`.
   * Handled by portable helpers: `internal/storage/disk_linux.go` (`unix.Statfs`) and `disk_other.go` (fallback for developer OSs).
3. **Direct-Play Codec Check:** Validates `.mp4`, `.webm`, and web-safe `.mkv`. Rejects legacy non-web containers like `.avi`.
4. **Zero-Memory Streaming:** Uses `r.MultipartReader` and `io.Copy` in 32 KB buffers. Memory usage stays <1 MB even when uploading multi-gigabyte video files.
5. **Immediate Indexing:** Writes files with `0644` permissions and indexes them immediately into SQLite.

---

## 5. Drive Disconnection & Mount Failure Defenses

| Failure Scenario | Risk | Flan Defensive Mitigation |
| :--- | :--- | :--- |
| **Unmounted Drive on Boot** | An empty directory on root flash is mistaken for a new database, creating a split-brain "ghost" database. | **Marker File (`.flan-keep`):** Before initializing an external DB path, Flan checks for `.flan-keep`. If absent, startup halts with a fatal log message. |
| **Drive Disconnect During Scan** | Scanner assumes all files were deleted and purges hundreds of database catalog rows. | **Mount Liveness Check:** Before scanning, Flan verifies that the directory exists and contains files. If absent or empty, the scan aborts safely, leaving database rows untouched. |
| **Streaming an Offline File** | Playback freezes or crashes the server handler. | **Soft-Offline Response:** The server returns `HTTP 503 Service Unavailable` (*"Media drive is offline"*). When reconnected, streaming resumes with zero re-indexing. |
| **Mechanical Drive Thrashing** | Multiple concurrent video streams cause mechanical head seeking, dropping throughput from 100 MB/s to 3 MB/s. | **Stream Governor Semaphore:** Caps active streams at 3 (Zone A), ensuring smooth sequential disk reads. |

---

## 6. Related Documentation

* [Master System Architecture](design.md)
* [Database Schema & Wear-Leveling](database.md)
* [Rate Limiting Architecture](rate-limiting.md)
* [Security Threat Model](threat-model.md)
