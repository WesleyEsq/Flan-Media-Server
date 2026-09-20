# Database Schema & Query Procedures

This document details the SQLite database location, storage considerations for single-board computers, complete table schemas, and core query procedures for Flan Media Server.

---

## 1. Database Location & Storage Considerations

The database runs on embedded SQLite3 via Go's database/sql package. Because the server frequently runs on SBCs (like Raspberry Pis) using micro-SD cards, where the database file resides and how it writes to disk requires careful planning.

### Where the Database Lives

The database file location is configurable via the DB_PATH variable in the .env file:

+ **Local Development & Default:** data/flan.db (created automatically on first boot).
+ **Production Homelab / SBC:** /var/lib/flan/flan.db or ~/.local/share/flan/flan.db.
+ **Attached Storage Recommendation:** If an external USB hard drive or SSD is attached to the SBC for media files, placing the database on that external drive (e.g. /mnt/storage/flan.db) is strongly recommended. External drives offer significantly higher write endurance and faster random I/O than micro-SD cards. Detailed multi-drive guidelines are documented in [docs/storage.md](file:///home/wess/Documents/MechanicalSpeak/Flan-Media-Server/docs/storage.md).

### Ghost Database Prevention (.flan-keep)

To prevent accidentally creating a fresh empty database on an unmounted boot drive when an external drive fails to mount at startup, Flan writes a hidden marker file (`.flan-keep`) in the database folder. If DB_PATH points to an external path and `.flan-keep` is absent, the server refuses to initialize a new database and halts with a fatal warning.

### SQLite Performance and Wear-Leveling Pragmas

When opening the database connection pool, the server applies these pragmas immediately:

```sql
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = -2000;
PRAGMA temp_store = MEMORY;
PRAGMA busy_timeout = 5000;
PRAGMA foreign_keys = ON;
PRAGMA wal_autocheckpoint = 1000;
```

#### Why These Pragmas Matter

+ **WAL Mode (Write-Ahead Logging):** Enables non-blocking concurrent readers during write transactions. Video streaming and catalog browsing never block when playback progress or scraper results are being written.
+ **Synchronous Normal:** In WAL mode, synchronous normal is safe against application crashes and reduces the frequency of fsync disk calls, significantly extending the lifespan of micro-SD cards.
+ **Cache Size (-2000):** Strictly limits SQLite page cache to roughly 2mb of RAM, supporting the overall 15 to 20mb server memory budget.
+ **Busy Timeout (5000ms):** Prevents SQLITE_BUSY errors during simultaneous progress updates by having Go automatically wait up to 5 seconds for write locks to clear.

### Pure-Go Driver Selection (modernc.org/sqlite)

To satisfy Flan's single-binary deployment model and seamless cross-compilation across heterogeneous SBC hardware (ARMv6, ARMv7, ARM64), the server utilizes the pure-Go SQLite driver `modernc.org/sqlite` instead of CGo-dependent alternatives (`github.com/mattn/go-sqlite3`).

#### Rationale:
+ **Zero-Friction Cross-Compilation:** Cross-compiling for Raspberry Pi Zero/1 (`GOARCH=arm GOARM=6`) or Raspberry Pi 4/5 (`GOARCH=arm64`) requires zero host C cross-compilers or system header dependencies. Setting `CGO_ENABLED=0` produces an immutable, self-contained binary.
+ **Standard Database Interface:** Registers cleanly as a standard `database/sql` driver (`sqlite`), preserving idiomatic Go query semantics.
+ **Pragma Compatibility:** Fully honors all low-memory pragmas (`cache_size = -2000`, `journal_mode = WAL`, and `synchronous = NORMAL`), operating comfortably within the 15 to 20mb RAM target.

---

## 2. Streamlined Three-Table Schema

To avoid relational complexity, excessive joins, and write contention on low-spec hardware, the database uses three core tables:

1. **users:** Manages profiles, PIN hashes, and brute-force lockout states.
2. **media_items:** Central catalog table storing movies, TV episodes, and books.
3. **playback_progress:** Tracks current playback position and completion status per user.

*Note on Sessions:* Session tokens are not stored in SQLite. Instead, the server uses stateless, HMAC-SHA256-signed session cookies containing the user ID, role, and expiration timestamp. This eliminates a database read on every HTTP request and removes the need for periodic session table cleanup. Key lifecycle:
+ **Configuration:** The server looks for a 32-byte hexadecimal `SESSION_SECRET` in the `.env` file or environment.
+ **Automatic Persistence:** If no secret is configured, the server inspects the database folder for a `.session_secret` file. If missing, it generates 32 cryptographically secure random bytes via `crypto/rand` and writes the file with restrictive `0600` permissions. This ensures active user sessions persist across daemon restarts without requiring manual administrator configuration.

*Note on Subtitles:* Subtitles are not stored in SQLite. The server discovers sidecar subtitle files (.srt and .vtt) directly on disk in the same directory as the video.

```sql
-- User Profiles
CREATE TABLE IF NOT EXISTS users (
    user_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    username         TEXT NOT NULL UNIQUE,
    pin_hash         TEXT NOT NULL,
    role             TEXT NOT NULL CHECK(role IN ('admin', 'user')),
    avatar_icon      TEXT DEFAULT 'flan',        -- Built-in icon name ('flan', 'popcorn', 'cat') or custom path
    avatar_color     TEXT DEFAULT '#bb9af7',     -- Profile accent background color
    failed_attempts  INTEGER DEFAULT 0,
    locked_until     DATETIME,
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Media Catalog Items (Movies, TV Episodes, Books)
CREATE TABLE IF NOT EXISTS media_items (
    media_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title            TEXT NOT NULL,
    type             TEXT NOT NULL CHECK(type IN ('movie', 'tv', 'book')),
    format           TEXT NOT NULL,              -- 'mp4', 'webm', 'epub', 'pdf'
    file_path        TEXT NOT NULL UNIQUE,
    file_size        INTEGER NOT NULL,
    duration         INTEGER DEFAULT 0,          -- In seconds (videos)
    cover_path       TEXT,                       -- Local path to cached cover image
    overview         TEXT,
    release_year     INTEGER,
    rating           REAL DEFAULT 0.0,
    genres           TEXT DEFAULT '',            -- Comma-separated list: 'Animation, Comedy'
    series_title     TEXT,                       -- For TV shows: e.g. 'Breaking Bad'
    season_number    INTEGER DEFAULT 0,          -- For TV shows: e.g. 1
    episode_number   INTEGER DEFAULT 0,          -- For TV shows: e.g. 5
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Playback and Reading Progress
CREATE TABLE IF NOT EXISTS playback_progress (
    user_id           INTEGER NOT NULL,
    media_id          INTEGER NOT NULL,
    position_seconds  INTEGER NOT NULL DEFAULT 0, -- Current second (video) or page (book)
    total_seconds     INTEGER NOT NULL DEFAULT 0, -- Total seconds (video) or pages (book)
    is_finished       INTEGER NOT NULL DEFAULT 0, -- 0 = in progress, 1 = completed
    updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, media_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE
);

-- Indexes for Fast Querying
CREATE INDEX IF NOT EXISTS idx_media_type ON media_items(type);
CREATE INDEX IF NOT EXISTS idx_media_tv ON media_items(series_title, season_number, episode_number);
CREATE INDEX IF NOT EXISTS idx_progress_updated ON playback_progress(user_id, updated_at DESC);
```

---

## 3. Core Query Procedures & Access Patterns

These are the primary database procedures needed by the server handlers:

### A. Dashboard Shelves (Continue Watching & Continue Reading)

Retrieves media that the active user has started but not finished, ordered by most recently updated:

```sql
SELECT 
    m.media_id,
    m.title,
    m.type,
    m.format,
    m.cover_path,
    m.duration,
    p.position_seconds,
    p.total_seconds,
    m.series_title,
    m.season_number,
    m.episode_number
FROM playback_progress p
JOIN media_items m ON p.media_id = m.media_id
WHERE p.user_id = ? 
  AND p.is_finished = 0 
  AND p.position_seconds > 10
  AND m.type = ? -- 'movie' or 'tv' for video shelf, 'book' for reading shelf
ORDER BY p.updated_at DESC
LIMIT 12;
```

---

### B. Catalog Filtering by Type and Genre

Loads movies, TV series, or books, optionally filtered by genre keyword:

```sql
-- When browsing standalone movies
SELECT media_id, title, cover_path, release_year, rating, duration, genres
FROM media_items
WHERE type = 'movie'
  AND (? IS NULL OR genres LIKE '%' || ? || '%')
ORDER BY title ASC;

-- When browsing TV Series (distinct series cards)
SELECT 
    series_title,
    cover_path,
    release_year,
    rating,
    genres,
    COUNT(media_id) AS episode_count
FROM media_items
WHERE type = 'tv'
  AND (? IS NULL OR genres LIKE '%' || ? || '%')
GROUP BY series_title
ORDER BY series_title ASC;

-- When browsing Books
SELECT media_id, title, cover_path, format, genres
FROM media_items
WHERE type = 'book'
  AND (? IS NULL OR genres LIKE '%' || ? || '%')
ORDER BY title ASC;
```

---

### C. TV Series Detail View (Seasons & Episodes)

Retrieves all episodes for a specific series alongside watch progress for the current user:

```sql
SELECT 
    m.media_id,
    m.season_number,
    m.episode_number,
    m.title,
    m.overview,
    m.duration,
    COALESCE(p.position_seconds, 0) AS position_seconds,
    COALESCE(p.is_finished, 0) AS is_finished
FROM media_items m
LEFT JOIN playback_progress p ON m.media_id = p.media_id AND p.user_id = ?
WHERE m.type = 'tv' AND m.series_title = ?
ORDER BY m.season_number ASC, m.episode_number ASC;
```

---

### D. Playback Progress Upsert

Updates playback progress atomically when the video player emits timeupdate pings:

```sql
INSERT INTO playback_progress (user_id, media_id, position_seconds, total_seconds, is_finished, updated_at)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, media_id) DO UPDATE SET
    position_seconds = excluded.position_seconds,
    total_seconds = excluded.total_seconds,
    is_finished = excluded.is_finished,
    updated_at = CURRENT_TIMESTAMP;
```

---

### E. PIN Authentication & Brute-Force Rate Limiting

Procedures to validate users, enforce rate limiting, and reset lockouts:

```sql
-- 1. Fetch user for PIN verification with lockout status
SELECT user_id, pin_hash, role, failed_attempts, locked_until
FROM users
WHERE user_id = ?;

-- 2. Increment failed attempts and lock account if threshold reached
UPDATE users
SET failed_attempts = failed_attempts + 1,
    locked_until = CASE 
        WHEN failed_attempts + 1 >= 5 THEN datetime(CURRENT_TIMESTAMP, '+5 minutes')
        ELSE locked_until 
    END
WHERE user_id = ?;

-- 3. Reset failed attempts on successful login
UPDATE users
SET failed_attempts = 0, locked_until = NULL
WHERE user_id = ?;

-- 4. Reset admin PIN from CLI command (./flan --reset-admin)
UPDATE users
SET pin_hash = ?, failed_attempts = 0, locked_until = NULL
WHERE role = 'admin';
```

---

### F. Library Scanning & Ingestion

Procedures used by the background scanner:

```sql
-- Check if file already exists in database before scraping
SELECT media_id, file_size FROM media_items WHERE file_path = ?;

-- Insert newly discovered media item
INSERT INTO media_items (
    title, type, format, file_path, file_size, duration,
    cover_path, overview, release_year, rating, genres,
    series_title, season_number, episode_number
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Clean up files removed from disk
DELETE FROM media_items WHERE file_path = ?;
```

---

### G. Database Integrity & Zero-Lock Snapshots

Procedures for power-cut resilience and hot backups:

```sql
-- 1. Fast integrity check executed on startup
PRAGMA quick_check;

-- 2. Zero-lock hot snapshot executed daily in background
VACUUM INTO 'data/flan.db.backup';
```
