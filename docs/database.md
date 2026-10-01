# Database Schema & Query Procedures

Flan Media Server uses embedded SQLite3 via Go's `database/sql` and the pure-Go driver `modernc.org/sqlite` (zero CGo, instant cross-compilation).

The schema is reduced to **6 straightforward tables** centered on containers, file lists, and unified progress tracking.

---

## 1. Storage Location & Wear-Leveling Pragmas

* **Database File (`DB_PATH`):** Defaults to `./data/flan.db` on fast flash storage (eMMC, NVMe, or root micro-SD).
* **Ghost Mount Defense (`.flan-keep`):** Before initializing an external DB path, Flan verifies that `.flan-keep` exists in the folder. If missing, startup halts immediately to prevent writing a blank database to an unmounted root card.

### Connection Pooling & Pragmas

```go
// internal/database/database.go
db.SetMaxOpenConns(1) // Single serialized connection avoids write contention & saves RAM
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(0)
```

| Pragma | Value | Purpose |
| :--- | :--- | :--- |
| `journal_mode` | `WAL` | Enables non-blocking concurrent reads while writes occur. |
| `synchronous` | `NORMAL` | Reduces `fsync` calls; safe in WAL mode; extends micro-SD card lifespan. |
| `cache_size` | `-2000` | Limits SQLite page cache strictly to ~2 MB of RAM. |
| `busy_timeout` | `5000` | Automatically waits up to 5 seconds for write locks to clear. |
| `foreign_keys` | `ON` | Enforces relational integrity and cascading deletes. |

---

## 2. Ultra-Streamlined Relational Schema (6 Tables)

```sql
-- 1. User Profiles & Lockouts
CREATE TABLE IF NOT EXISTS users (
    user_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username        TEXT NOT NULL UNIQUE,
    pin_hash        TEXT NOT NULL,
    role            TEXT NOT NULL CHECK(role IN ('admin', 'user')),
    token_version   INTEGER NOT NULL DEFAULT 1,
    avatar_icon     TEXT DEFAULT 'default', -- bundled SVG icon name
    avatar_path     TEXT,                   -- optional uploaded custom avatar image
    failed_attempts INTEGER DEFAULT 0,
    locked_until    DATETIME,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Video Containers (Represents a Movie or TV Series)
CREATE TABLE IF NOT EXISTS videos (
    video_id        INTEGER PRIMARY KEY AUTOINCREMENT,
    title           TEXT NOT NULL UNIQUE,
    overview        TEXT,
    cover_path      TEXT,
    folder_path     TEXT NOT NULL UNIQUE, -- Relative path under ./media/video/
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Playable Video Files (Episodes of a Series OR Versions/Cuts of a Movie)
CREATE TABLE IF NOT EXISTS video_files (
    file_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    video_id         INTEGER NOT NULL REFERENCES videos(video_id) ON DELETE CASCADE,
    title            TEXT NOT NULL, -- e.g. "S01E01 - Pilot" or "Director's Cut" or "Movie"
    relative_path    TEXT NOT NULL UNIQUE, -- Path under ./media/video/
    file_size        INTEGER NOT NULL,
    format           TEXT NOT NULL, -- 'mp4', 'webm', 'mkv'
    duration_seconds INTEGER DEFAULT 0,
    order_index      INTEGER DEFAULT 0
);

-- 4. Book Containers (Represents a Book or Multi-Volume Series)
CREATE TABLE IF NOT EXISTS books (
    book_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title           TEXT NOT NULL UNIQUE,
    author          TEXT,
    overview        TEXT,
    cover_path      TEXT,
    folder_path     TEXT NOT NULL UNIQUE, -- Relative path under ./media/books/
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 5. Book Files (Volumes, Editions, or Single Files)
CREATE TABLE IF NOT EXISTS book_files (
    file_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id          INTEGER NOT NULL REFERENCES books(book_id) ON DELETE CASCADE,
    title            TEXT NOT NULL,
    relative_path    TEXT NOT NULL UNIQUE, -- Path under ./media/books/
    file_size        INTEGER NOT NULL,
    format           TEXT NOT NULL CHECK(format IN ('epub', 'pdf')),
    order_index      INTEGER DEFAULT 0
);

-- 6. Unified Progress Tracking (Video & Books)
CREATE TABLE IF NOT EXISTS progress (
    user_id          INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    media_type       TEXT NOT NULL CHECK(media_type IN ('video', 'book')),
    file_id          INTEGER NOT NULL,
    position_data    TEXT NOT NULL, -- Seconds string (e.g. "1450") or EPUB CFI / Page string
    percentage       REAL DEFAULT 0.0,
    is_finished      INTEGER NOT NULL DEFAULT 0,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, media_type, file_id)
);

-- Fast Index Lookups
CREATE INDEX IF NOT EXISTS idx_video_files_video ON video_files(video_id, order_index);
CREATE INDEX IF NOT EXISTS idx_book_files_book ON book_files(book_id, order_index);
CREATE INDEX IF NOT EXISTS idx_progress_lookup ON progress(user_id, media_type, updated_at DESC);
```

---

## 3. Core Query Patterns

### Video Catalog List
```sql
SELECT video_id, title, cover_path FROM videos ORDER BY title ASC;
```

### Video Detail with Playable File List
```sql
SELECT 
    f.file_id,
    f.title,
    f.duration_seconds,
    f.order_index,
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

### Zero-Lock Backups & Health Checks
* **Integrity Check (Boot):** `PRAGMA quick_check;` validates file consistency on startup.
* **Daily Hot Backup:** `VACUUM INTO 'data/flan.db.backup';` takes an atomic non-blocking snapshot every 24 hours.

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Storage Architecture](storage.md)
* [Security Threat Model](threat-model.md)
