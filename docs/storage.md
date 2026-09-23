# Storage Architecture & Drive Resiliency

This document details the multi-drive storage architecture, dynamic library management, deterministic folder structures, upload routing, mount failure defenses, and drive disconnection handling for Flan Media Server.

---

## 1. Real-World SBC Hardware Landscape

Single-board computers in homelabs rarely run off a single storage device. Users commonly operate with a tiered, multi-drive setup:

```text
[SBC System]
   ├── Boot Tier: Micro-SD Card (32gb to 64gb) or eMMC
   │     └── Operating system, systemd, base system (read-heavy, low writes)
   │
   ├── App Tier: Fast Flash / NVMe M.2 SSD / SATA HAT (128gb to 512gb)
   │     ├── flan executable
   │     ├── data/flan.db (SQLite database in WAL mode)
   │     └── data/covers/ (Locally cached image artwork)
   │     → Stays awake 24/7, high random I/O, sub-millisecond UI responses
   │
   └── Bulk Media Tier: External USB 3.0 HDD / Multi-Drive Storage Pool
         ├── /mnt/hdd1/movies (Standalone films)
         ├── /mnt/hdd2/tv (TV series and episodes)
         └── /mnt/nvme/books (Books and documents)
         → Mechanical drives spin down when idle; sequential streaming reads
```

Because drives can be scattered across different buses (micro-SD, PCIe NVMe, USB 3.0), the server handles drive spin-down, mount ordering, disconnections, and power cuts without corrupting data or breaking user workflows.

---

## 2. Dynamic Library Management (`libraries` Table)

Instead of relying on rigid, single-path `.env` definitions, Flan manages media roots as dynamic records in the SQLite database:

```sql
CREATE TABLE IF NOT EXISTS libraries (
    library_id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT NOT NULL,
    path         TEXT NOT NULL UNIQUE,
    media_type   TEXT NOT NULL CHECK(media_type IN ('movies', 'tv', 'books')),
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### What Happens If `.env` is Blank or Missing?
1. **Graceful Startup:** If `.env` has no `MEDIA_DIR` or if the variable is left blank, Flan starts normally without error.
2. **Setup Wizard Default:** If the `libraries` table is empty on first boot, the `/setup` wizard creates the admin account and prompts to create default libraries under `./media`:
   - `./media/movies` (Type: `movies`)
   - `./media/tv` (Type: `tv`)
   - `./media/books` (Type: `books`)
3. **Environment Import:** If `MEDIA_DIR` is set in `.env`, Flan imports it into the `libraries` table on first boot as the default library.
4. **Dynamic Administration:** Administrators can add, inspect, or delete libraries anytime through the `/settings` web page or via `POST /api/libraries` without stopping or restarting the server.

---

## 3. Deterministic Folder Structure for Media Intake

Because every library has an explicit `media_type`, the file intake pipeline operates deterministically without guessing whether a directory contains movies, TV series, or books:

### A. Movies Library (`media_type = 'movies'`)
Files can be organized in either of two standard layouts:
```text
/mnt/storage/movies/
├── The Matrix (1999).mp4
└── Spirited Away (2001)/
    ├── Spirited Away (2001).mp4
    └── poster.jpg                 # Optional local cover
```

### B. TV Series Library (`media_type = 'tv'`)
Strict hierarchical structure mapping to `series` and `episodes` tables:
```text
/mnt/storage/tv/
└── Breaking Bad/
    ├── poster.jpg                 # Show-level cover art
    ├── Season 01/
    │   ├── Breaking Bad - S01E01 - Pilot.mp4
    │   └── Breaking Bad - S01E02 - Cat's in the Bag.mp4
    └── Season 02/
        └── Breaking Bad - S02E01 - Seven Thirty-Seven.mp4
```

### C. Books Library (`media_type = 'books'`)
```text
/mnt/storage/books/
├── Frank Herbert/
│   └── Dune.epub                  # Embedded cover auto-extracted
└── Documentation/
    └── Linux Kernel.pdf
```

---

## 4. Admin Web Upload Routing & Intake

Administrators can upload media directly through the web client without SSH or terminal access.

### Upload Workflow:
1. **Target Library Selection:** The admin selects the destination library from a dropdown populated by `GET /api/libraries`.
2. **Type-Specific Routing:**
   + **Movies Library:** Files stream directly into `<library_path>/<filename>` or `<library_path>/<Folder>/<filename>`.
   + **TV Series Library:**
     + The upload modal prompts for:
       - **Series Title** (text input or dropdown of existing series in this library).
       - **Season Number** (numeric input, defaults to 1).
     - Alternatively, dropping a folder via `<input webkitdirectory>` preserves relative directory paths (`<Series>/Season <NN>/<file>`).
     - Files stream directly to `<library_path>/<Series Title>/Season <NN>/<filename>`.
   + **Books Library:**
     - The modal allows an optional **Author** field (defaults to "Unknown" or extracted from EPUB metadata).
     - Files stream to `<library_path>/<Author>/<filename>`.
3. **Mount-Specific Free Space Verification:**
   Before accepting bytes, the server queries available disk space on the selected library's mount point. If free space is below **2gb**, the upload is immediately rejected with HTTP 507 Insufficient Storage. Filesystem free-space queries are decoupled behind a portable helper in `internal/storage` using Go build tags:
   + `disk_linux.go` (`//go:build linux`): Implements `CheckFreeSpace(path string)` using `unix.Statfs` or `syscall.Statfs` to calculate available blocks on production Linux single-board computers.
   + `disk_other.go` (`//go:build !linux`): Provides a portable fallback for developer workstations (macOS, Windows, BSD), preventing cross-compilation and local testing failures.
4. **Direct Play Format Verification:**
   The intake pipeline verifies that incoming video files match web-compatible standards (`.mp4`, `.webm`). Files with `.mkv` extensions are accepted only if they contain web-safe codecs (H.264/VP9 and AAC/Opus). Uploading legacy non-web containers (like `.avi`) or incompatible audio codecs (DTS/AC3) triggers a client-side warning advising external playback, as on-the-fly transcoding is strictly excluded.
5. **Zero-Memory Streaming:**
   File bytes stream straight from `r.MultipartReader` to disk using `io.Copy` in 32kb buffers. Memory usage stays under 1MB even when uploading a 4GB video file.
6. **Immediate Indexing:**
   Once written to disk with `0644` permissions, the file is immediately indexed into the SQLite database and metadata is fetched in the background.

---

## 5. Preventing the Ghost Database (The Empty Mount Trap)

### The Failure Mode
If the administrator configures `DB_PATH=/mnt/nvme/flan/flan.db`, and the NVMe drive fails to mount at boot (due to systemd mount order, dirty filesystem, or a loose connection), Linux leaves an empty directory at `/mnt/nvme` on the root micro-SD card. A naive server would initialize a fresh empty database on the boot card, causing split-brain confusion.

### The Flan Defenses:
1. **The Marker File (`.flan-keep`):**
   When Flan initializes a database directory, it writes a hidden marker file next to the database:
   `/mnt/nvme/flan/.flan-keep`
2. **Mount Verification on Boot:**
   Before initializing a database, the server inspects `DB_PATH`:
   + If `DB_PATH` is pointing to an external mount and the directory or `.flan-keep` file is missing, the server **REFUSES to initialize a new database**.
   + The server halts with a clear fatal log message:
     ```text
     FATAL: Database storage path '/mnt/nvme/flan' is missing marker file.
     The server will NOT initialize a new database to prevent data loss.
     Please verify that your external storage drive is mounted properly.
     ```

---

## 6. Isolating Mechanical Drive Spin-Up

External USB mechanical hard drives spin down after 5 to 15 minutes of inactivity to save power and reduce heat. Waking up a spinning drive takes 3 to 7 seconds.

### The Decoupling Rule:
+ **App Data & Covers on Fast Storage:** The SQLite database (`flan.db`) and image thumbnails (`data/covers/`) live on the fast boot drive, eMMC, or NVMe SSD.
+ **Media on Mechanical Drives:** Only video and document files live on the spinning drive.
+ **The Result:** Browsing catalogs, filtering genres, viewing posters, and updating progress never touches the mechanical drive. The drive only wakes up when a user clicks Play on a video stream (`/stream/{type}/{id}`).

---

## 7. Drive Disconnections & Safe Library Scanning

USB cables can be bumped and drives can temporarily unmount while the server is running.

### The Catalog Wipe Hazard:
In naive media servers, if a scan runs while a drive is disconnected, the scanner assumes files were deleted and purges hundreds of database rows.

### The Flan Safeguards:
1. **Directory Liveness Check:**
   Before scanning any library path, the scanner verifies:
   + Does the directory exist?
   + Does the directory contain files?
   + If a folder is empty or absent, the scanner skips it, logs a warning (`Library [Movies] is offline`), and **leaves all database records untouched**.
2. **Soft-Offline State:**
   + Catalog items remain visible in the client UI even when the physical drive is unplugged.
   + If a user attempts to stream an offline item, the server returns HTTP 503 Service Unavailable:
     `"Media storage drive is offline. Please check drive connection."`
   + As soon as the drive is reconnected, streaming resumes immediately with zero re-indexing required.
3. **Database Lock Yielding During Scans:**
   + To prevent scan jobs from monopolizing SQLite's single connection (`db.SetMaxOpenConns(1)`), newly discovered items are persisted using small transactions or individual upserts.
   + The scanner yields execution (`runtime.Gosched()`) between files, ensuring active video streams, playback progress updates, and user web navigation are never blocked by long-running background scans.

---

## 8. Power-Cut Insurance (SQLite Snapshots)

1. **WAL Mode + Synchronous Normal:**
   SQLite in Write-Ahead Logging mode with synchronous normal is crash-resilient and minimizes flash memory wear by batching `fsync` operations.
2. **Startup Quick Check:**
   Every time the server boots, it runs `PRAGMA quick_check;`. If errors are detected, the server halts safely rather than writing to a corrupted file.
3. **Automated Zero-Lock Daily Snapshots:**
   Once every 24 hours, a lightweight background goroutine runs `VACUUM INTO 'data/flan.db.backup';` to create an atomic, hot snapshot without locking active readers.

---

### Related Documentation

+ [Master System Specifications](design.md)
+ [Database Schema & Wear-Leveling Pragmas](database.md)
+ [Cross-Compilation & SBC Deployment Guide](compilation.md)
+ [Five-Zone Rate Limiting Architecture](rate-limiting.md)
+ [Security Threat Model & Mount Hijacking Defense](threat-model.md)
