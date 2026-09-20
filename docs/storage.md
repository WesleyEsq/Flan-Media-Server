# Storage Architecture & Drive Resiliency

This document details the multi-drive storage architecture, mount failure defenses, database protection, and drive disconnection handling for Flan Media Server.

---

## 1. Real-World SBC Hardware Landscape

Single-board computers in homelabs rarely run off a single storage device. Users commonly operate with a tiered, multi-drive setup:

```
[SBC System]
   ├── Boot Tier: Micro-SD Card (32gb to 64gb) or eMMC
   │     └── Operating system, systemd, base system (read-heavy, low writes)
   │
   ├── App Tier: NVMe M.2 SSD or SATA HAT (128gb to 512gb)
   │     ├── flan executable
   │     ├── data/flan.db (SQLite database in WAL mode)
   │     └── data/covers/ (Locally cached image artwork)
   │     → Stays awake 24/7, high random I/O, sub-millisecond UI responses
   │
   └── Bulk Media Tier: External USB 3.0 HDD / SATA HAT / Storage Pool
         ├── /mnt/hdd1/movies (Standalone films)
         └── /mnt/hdd2/tv (TV seasons and episodes)
         → Spins down when idle, large capacity, sequential streaming reads
```

Because drives can be scattered across different buses (micro-SD, PCIe NVMe, USB 3.0), the server must handle drive spin-down, mount ordering, disconnections, and power cuts without corrupting data or breaking user workflows.

---

## 2. Preventing the Ghost Database (The Empty Mount Trap)

### The Failure Mode
If the administrator configures DB_PATH=/mnt/nvme/flan/flan.db, and the NVMe drive fails to mount at boot (due to systemd mount order, dirty filesystem, or a loose connection), Linux leaves an empty directory at /mnt/nvme on the root micro-SD card.

A naive server would see that the database does not exist, create a fresh empty database on the micro-SD card, detect zero users, and launch the /setup onboarding wizard. The user panics, thinking their library was wiped, and later mounting the NVMe drive causes split-brain state.

### The Flan Defenses:
1. **The Marker File (.flan-keep):**
   When Flan initializes a database directory, it writes a small hidden marker file next to the database:
   /mnt/nvme/flan/.flan-keep

2. **Mount Verification on Boot:**
   Before initializing a database, the server inspects DB_PATH:
   + If DB_PATH is pointing to an external mount (outside the current executable directory) and the directory or .flan-keep file is missing, the server **REFUSES to initialize a new database**.
   + The server halts with a clear fatal log message:
     ```
     FATAL: Database storage path '/mnt/nvme/flan' is missing marker file.
     The server will NOT initialize a new database to prevent data loss.
     Please verify that your external storage drive is mounted properly.
     ```
   + This completely prevents creating ghost databases on the root boot drive.

---

## 3. Isolating Mechanical Drive Spin-Up

External 2.5" and 3.5" USB mechanical hard drives spin down after 5 to 15 minutes of inactivity to save power and reduce heat. Spinning up a mechanical drive takes 3 to 7 seconds.

### The Golden Decoupling Rule:
+ **App Data & Covers on Fast Storage:** The SQLite database (flan.db) and image thumbnails (data/covers/) must live on the fast boot drive, eMMC, or NVMe SSD.
+ **Media on Mechanical Drives:** Only video and document files live on the spinning mechanical hard drive.
+ **The Result:** Browsing the catalog, searching genres, viewing posters, and updating profile settings never touches the mechanical drive. The web interface responds in milliseconds. The spinning drive only wakes up when a user clicks Play on a video stream (/stream/{id}).

---

## 4. Drive Disconnections & Safe Library Scanning

USB cables can be bumped, USB hubs can experience power dips, and external drives can temporarily unmount while the server is running.

### The Catalog Wipe Hazard:
In naive media servers, if a background library scan runs while a USB drive is disconnected, the scanner notices that 500 movie files are missing from disk, assumes the user deleted them, and deletes all 500 rows from SQLite. When the drive is plugged back in, the server has to re-scrape the entire library from scratch.

### The Flan Safeguards:
1. **Directory Liveness Check:**
   Before scanning any configured media directory (e.g. /mnt/hdd1/movies), the scanner verifies:
   + Does the directory exist?
   + Does the directory contain files?
   + If a folder is completely empty or missing, the scanner skips it entirely, logs a warning (Drive /mnt/hdd1 is offline), and **leaves all database records untouched**.
2. **Catalog Integrity & Soft-Offline State:**
   + Catalog cards, metadata, and watch progress remain visible in the web interface even when the physical drive is unplugged.
   + If a user attempts to stream a file whose drive is currently missing, the server responds with HTTP 503 Service Unavailable:
     "Media storage drive is offline. Please check drive connection."
   + As soon as the drive is reconnected, streaming resumes immediately with zero rescanning required.

---

## 5. Multi-Drive Upload Routing & Disk Space Checks

For administrators uploading files through the web interface:

1. **Explicit Destination Selection:**
   When opening the upload modal, the administrator selects which configured directory to place the files into (for example, "Movies on USB Drive (/mnt/hdd1/movies)").
2. **Mount-Specific Free Space Verification:**
   Before accepting file bytes, the server runs statfs on the exact directory chosen.
   + If free space on that specific drive drops below 2gb, the upload is rejected immediately with HTTP 507 Insufficient Storage.
   + This guarantees that a large upload to an external USB drive will never fill up or crash the micro-SD boot card.

---

## 6. Power-Cut Insurance (SQLite Snapshots)

SBCs in homelabs often suffer unexpected power loss when cords are bumped or power blips occur.

### Durability Defenses:
1. **WAL Mode + Synchronous Normal:**
   SQLite in Write-Ahead Logging mode with synchronous normal is crash-resilient and minimizes flash memory wear by batching fsync operations.
2. **Startup Quick Check:**
   Every time the server boots, it runs:
   ```sql
   PRAGMA quick_check;
   ```
   If errors are detected, the server halts safely rather than writing to a corrupted file.
3. **Automated Zero-Lock Daily Snapshots:**
   Once every 24 hours, a lightweight background goroutine runs:
   ```sql
   VACUUM INTO 'data/flan.db.backup';
   ```
   This command creates an atomic, consistent hot copy of the database without locking active readers or writers. If a catastrophic corruption ever occurs, the administrator can restore from the backup file.

---

### Related Documentation

+ [Master System Specifications](docs/design.md)
+ [Database Schema & Wear-Leveling Pragmas](docs/database.md)
+ [Cross-Compilation & SBC Deployment Guide](docs/compilation.md)
+ [Five-Zone Rate Limiting Architecture](docs/rate-limiting.md)
+ [Security Threat Model & Mount Hijacking Defense](docs/threat-model.md)
