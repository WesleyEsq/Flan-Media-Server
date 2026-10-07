# Database Schema & Query Procedures

Flan Media Server uses an embedded SQLite3 via Go's `database/sql` and the pure-Go driver `modernc.org/sqlite` (no CGo, standalone cross-compilation).

The database contains 7 relational tables representing users, storage sources, media containers, media files, and user progress.

---

## 1. Storage Location & Connection Pragmas

* **Database File (`DB_PATH`):** Defaults to `./data/flan.db` on fast flash storage (eMMC, NVMe, or root micro-SD).
* **Mount Marker Verification:** Before initializing SQLite, Flan verifies that `.flan-keep` exists in the database directory. If missing, startup halts immediately to prevent writing a blank database to an unmounted mount point.

### Connection Configuration

In SQLite WAL mode, concurrent readers do not block writers, and writers do not block readers. To support parallel page rendering without write contention, Flan separates database handles:

```go
// Dedicated Writer DB handle (serialized writes)
writerDB.SetMaxOpenConns(1)
writerDB.SetMaxIdleConns(1)

// Dedicated Reader DB handle (parallel HTTP reads)
readerDB.SetMaxOpenConns(3)
readerDB.SetMaxIdleConns(3)
readerDB.SetConnMaxLifetime(0)
```

### Essential Pragmas

The following pragmas are executed on every connection initialization:

* `PRAGMA journal_mode = WAL;`
  Enables Write-Ahead Logging for non-blocking concurrent reads during active writes.
* `PRAGMA synchronous = NORMAL;`
  Reduces `fsync` system calls; safe in WAL mode and extends flash memory (micro-SD / eMMC) lifespan.
* `PRAGMA cache_size = -2000;`
  Constrains SQLite memory cache strictly to ~2 MB of RAM per connection.
* `PRAGMA busy_timeout = 5000;`
  Instructs queries to wait up to 5 seconds for write locks to clear before failing.
* `PRAGMA foreign_keys = ON;`
  Enforces relational foreign key constraints and cascading deletes.
* `PRAGMA data_version;`
  Queried dynamically for instant, zero-overhead change detection between writer and reader pools or external CLI processes (`--reset-admin`). Returns an in-memory integer that increments on every committed write transaction without disk I/O.
* `PRAGMA quick_check;`
  Executed at server startup across database handles to verify B-tree and WAL structural integrity without full table scans.

---

## 2. Relational Schema (7 Tables)

```sql
-- 1. User Profiles (Auth & Display)
CREATE TABLE IF NOT EXISTS users (
    user_id         INTEGER PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    display_name    TEXT,                   -- optional friendly household nickname
    pin_hash        TEXT NOT NULL,          -- bcrypt hash of 4-6 digit numeric PIN
    role            TEXT NOT NULL CHECK(role IN ('admin', 'user')),
    token_version   INTEGER NOT NULL DEFAULT 1,
    avatar_icon     TEXT DEFAULT 'default', -- bundled SVG icon name
    avatar_path     TEXT,                   -- optional uploaded custom avatar image
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Configurable Storage Sources (Multi-Drive Roots)
CREATE TABLE IF NOT EXISTS storage_sources (
    source_id       INTEGER PRIMARY KEY,
    name            TEXT NOT NULL,          -- e.g. "Main USB Movies", "TV Drive", "Books"
    media_type      TEXT NOT NULL CHECK(media_type IN ('video', 'book')),
    folder_path     TEXT NOT NULL UNIQUE,   -- Absolute or root-relative directory path
    is_active       INTEGER NOT NULL DEFAULT 1,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Video Containers (Represents a Movie or TV Series)
CREATE TABLE IF NOT EXISTS videos (
    video_id        INTEGER PRIMARY KEY,
    source_id       INTEGER NOT NULL REFERENCES storage_sources(source_id) ON DELETE CASCADE,
    title           TEXT NOT NULL,          -- Display title (freely editable, protected from rescan)
    video_type      TEXT NOT NULL DEFAULT 'series' CHECK(video_type IN ('movie', 'series')),
    release_year    INTEGER,                -- e.g. 1982 or 2008
    overview        TEXT,
    cover_path      TEXT,                   -- Path to cached or uploaded poster art
    folder_path     TEXT NOT NULL,          -- Directory name or flat filename relative to storage source
    is_hidden       INTEGER NOT NULL DEFAULT 0, -- 1 = hidden from catalog by administrator
    metadata_locked INTEGER NOT NULL DEFAULT 1, -- 1 = preserve user/admin edits against rescan overwrite
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_id, folder_path)
);

-- 4. Playable Video Files (Episodes of a Series OR Versions/Cuts of a Movie)
CREATE TABLE IF NOT EXISTS video_files (
    file_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    video_id         INTEGER NOT NULL REFERENCES videos(video_id) ON DELETE CASCADE,
    title            TEXT NOT NULL,         -- Auto-detected or fallback title
    custom_title     TEXT,                  -- User-specified title override (e.g. "S01E01 - Pilot")
    relative_path    TEXT NOT NULL,         -- Relative path within container (empty string for flat files)
    file_size        INTEGER NOT NULL,
    mtime            INTEGER DEFAULT 0,     -- Unix timestamp for fast change detection
    format           TEXT NOT NULL,         -- 'mp4', 'webm', 'mkv'
    duration_seconds INTEGER DEFAULT 0,
    season_number    INTEGER DEFAULT 1,
    episode_number   INTEGER DEFAULT 0,
    order_index      INTEGER DEFAULT 0,     -- Calculated: (season_number * 1000) + episode_number
    is_hidden        INTEGER NOT NULL DEFAULT 0, -- 1 = hide extras/sample clips from public listing
    is_missing       INTEGER NOT NULL DEFAULT 0, -- 1 = temporarily missing from disk
    UNIQUE(video_id, relative_path)
);

-- 5. Book Containers (Represents a Book or Multi-Volume Series)
CREATE TABLE IF NOT EXISTS books (
    book_id         INTEGER PRIMARY KEY,
    source_id       INTEGER NOT NULL REFERENCES storage_sources(source_id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    author          TEXT,
    overview        TEXT,
    cover_path      TEXT,
    folder_path     TEXT NOT NULL,          -- Directory name or flat filename relative to storage source
    is_hidden       INTEGER NOT NULL DEFAULT 0, -- 1 = hidden from catalog by administrator
    metadata_locked INTEGER NOT NULL DEFAULT 1, -- 1 = preserve admin edits against rescan overwrite
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_id, folder_path)
);

-- 6. Book Files (Volumes, Editions, or Single Files)
CREATE TABLE IF NOT EXISTS book_files (
    file_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id          INTEGER NOT NULL REFERENCES books(book_id) ON DELETE CASCADE,
    title            TEXT NOT NULL,
    custom_title     TEXT,                  -- User-specified volume title override
    relative_path    TEXT NOT NULL,         -- Relative path within container (empty string for single-file books)
    file_size        INTEGER NOT NULL,
    mtime            INTEGER DEFAULT 0,     -- Unix timestamp for fast change detection
    format           TEXT NOT NULL CHECK(format IN ('epub', 'pdf')),
    order_index      INTEGER DEFAULT 0,
    is_hidden        INTEGER NOT NULL DEFAULT 0,
    is_missing       INTEGER NOT NULL DEFAULT 0, -- 1 = temporarily missing from disk
    UNIQUE(book_id, relative_path)
);

-- 7. Unified Progress Tracking (Video & Books)
CREATE TABLE IF NOT EXISTS progress (
    user_id          INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    media_type       TEXT NOT NULL CHECK(media_type IN ('video', 'book')),
    file_id          INTEGER NOT NULL,
    position_data    TEXT NOT NULL,         -- Video: seconds string (e.g. "1450.5"); Book: "unread", "reading", "finished"
    percentage       REAL DEFAULT 0.0,
    is_finished      INTEGER NOT NULL DEFAULT 0,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, media_type, file_id)
);

-- Indices for Catalog Lookups & Scans
CREATE INDEX IF NOT EXISTS idx_videos_source ON videos(source_id, video_type, is_hidden);
CREATE INDEX IF NOT EXISTS idx_books_source ON books(source_id, is_hidden);
CREATE INDEX IF NOT EXISTS idx_video_files_video ON video_files(video_id, order_index);
CREATE INDEX IF NOT EXISTS idx_book_files_book ON book_files(book_id, order_index);
CREATE INDEX IF NOT EXISTS idx_progress_lookup ON progress(user_id, media_type, updated_at DESC);

-- Automated Progress Cleanup Triggers (Eliminates Orphaned Progress on File Deletions)
CREATE TRIGGER IF NOT EXISTS trg_cleanup_video_progress
AFTER DELETE ON video_files
BEGIN
    DELETE FROM progress WHERE media_type = 'video' AND file_id = OLD.file_id;
END;

CREATE TRIGGER IF NOT EXISTS trg_cleanup_book_progress
AFTER DELETE ON book_files
BEGIN
    DELETE FROM progress WHERE media_type = 'book' AND file_id = OLD.file_id;
END;
```

---

## 3. Common Query Procedures

### Video Catalog List (Fast & Filter-Aware)

```sql
SELECT v.video_id, v.title, v.release_year, v.video_type, v.cover_path
FROM videos v
JOIN storage_sources s ON v.source_id = s.source_id
WHERE v.is_hidden = 0 
  AND s.is_active = 1
  AND EXISTS (SELECT 1 FROM video_files vf WHERE vf.video_id = v.video_id AND vf.is_missing = 0)
ORDER BY v.title ASC;
```

### Video Detail with Playable File List

```sql
SELECT 
    f.file_id,
    f.title,
    f.duration_seconds,
    f.order_index,
    f.is_missing,
    COALESCE(p.position_data, '0') AS position_data,
    COALESCE(p.is_finished, 0) AS is_finished
FROM video_files f
LEFT JOIN progress p ON f.file_id = p.file_id AND p.media_type = 'video' AND p.user_id = ?
WHERE f.video_id = ?
ORDER BY f.order_index ASC, f.title ASC;
```

### Unified Progress Upsert

```sql
INSERT INTO progress (user_id, media_type, file_id, position_data, percentage, is_finished, updated_at)
VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, media_type, file_id) DO UPDATE SET
    position_data = excluded.position_data,
    percentage = excluded.percentage,
    is_finished = excluded.is_finished,
    updated_at = CURRENT_TIMESTAMP;
```

### Integrity Checks and Backups

* **Startup Health Check:** `PRAGMA quick_check;` runs during startup to verify database consistency.
* **Safe Snapshot Backup:**
  1. Run `VACUUM INTO 'data/flan.db.tmp';` via a background worker.
  2. Atomically rename `data/flan.db.tmp` to `data/flan.db.backup` via `os.Rename`.
  This guarantees a consistent point-in-time copy without blocking active database readers.

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Storage Architecture](storage.md)
* [Security Threat Model](threat-model.md)
