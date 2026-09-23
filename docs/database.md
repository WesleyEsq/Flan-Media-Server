# Database Schema & Query Procedures

This document details the SQLite database location, storage considerations for single-board computers, complete normalized table schemas evaluated up to the 4th Normal Form (4NF), and core query procedures for Flan Media Server.

---

## 1. Database Location & Storage Considerations

The database runs on embedded SQLite3 via Go's `database/sql` package. Because the server frequently runs on SBCs (like Raspberry Pis) using micro-SD cards, where the database file resides and how it writes to disk requires careful planning.

### Where the Database Lives

The database file location is configurable via the `DB_PATH` variable in the `.env` file:

+ **Local Development & Default:** `data/flan.db` (created automatically on first boot).
+ **Production Homelab / SBC:** `/var/lib/flan/data/flan.db` or `~/.local/share/flan/flan.db`.
+ **Attached Storage Recommendation:** If an external USB hard drive or SSD is attached to the SBC for media files, placing the database on fast flash storage (eMMC, NVMe, or root micro-SD) while media sits on the spinning drive is strongly recommended. Detailed multi-drive guidelines are documented in [docs/storage.md](storage.md).

### Ghost Database Prevention (`.flan-keep`)

To prevent accidentally creating a fresh empty database on an unmounted boot drive when an external drive fails to mount at startup, Flan writes a hidden marker file (`.flan-keep`) in the database folder. If `DB_PATH` points to an external path and `.flan-keep` is absent, the server refuses to initialize a new database and halts with a fatal warning.

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
+ **Synchronous Normal:** In WAL mode, synchronous normal is safe against application crashes and reduces the frequency of `fsync` disk calls, significantly extending the lifespan of micro-SD cards.
+ **Cache Size (-2000):** Strictly limits SQLite page cache to roughly 2mb of RAM, supporting the overall 15 to 20mb server memory budget.
+ **Busy Timeout (5000ms):** Prevents `SQLITE_BUSY` errors during simultaneous progress updates by having Go automatically wait up to 5 seconds for write locks to clear.

### Pure-Go Driver Selection (`modernc.org/sqlite`) & Connection Pooling

To satisfy Flan's single-binary deployment model and seamless cross-compilation across heterogeneous SBC hardware (ARMv6, ARMv7, ARM64), the server utilizes the pure-Go SQLite driver `modernc.org/sqlite` instead of CGo-dependent alternatives (`github.com/mattn/go-sqlite3`).

#### Connection Pool Configuration

Because `modernc.org/sqlite` is implemented in pure Go and runs within strict memory boundaries (`GOMEMLIMIT=16MiB`), unbounded connection pools can quickly exhaust RAM. In `internal/database`, the pool is strictly constrained:

```go
db.SetMaxOpenConns(1) // Single writer/reader serialization ensures zero lock contention and minimal memory
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(0)
```

##### Transaction Granularity & Lock Yielding

With `SetMaxOpenConns(1)`, all read and write queries share a single serialized connection handle. To ensure that background operations (such as initial library crawling or recursive directory scanning) never starve high-priority foreground operations (such as catalog browsing, video chunk delivery, or throttled playback progress syncs every 5 seconds):

1. **Short, Granular Transactions:** Heavy indexing routines must never wrap entire directory trees in a single monolithic transaction. Instead, files are indexed individually or in small batches of 5 to 10 items.
2. **Cooperative Yielding:** The background scanner explicitly yields execution (`runtime.Gosched()` and minimal inter-item delays) between file checks. This allows pending HTTP requests waiting on `busy_timeout = 5000` to acquire the database handle immediately without latency spikes.

### Dedicated Database Initialization & Configuration (`database.go`)

To maintain clean architectural separation, all SQLite driver initialization, configuration, and migration mechanics are isolated into a single dedicated file: `internal/database/database.go`.

```go
// internal/database/database.go
package database

// Open initializes the SQLite connection, enforces pragmas, and runs migrations
func Open(dbPath string) (*sql.DB, error) {
    // 1. Verify .flan-keep marker file if dbPath is on an external mount
    // 2. Open modernc.org/sqlite database handle
    // 3. Configure connection pool: SetMaxOpenConns(1), SetMaxIdleConns(1)
    // 4. Apply pragmas (WAL, synchronous=NORMAL, cache_size=-2000, busy_timeout=5000)
    // 5. Execute DDL migrations to create all normalized tables and indexes
    // 6. Run PRAGMA quick_check;
    // 7. Return configured *sql.DB
}

// Snapshot executes an atomic, zero-lock hot backup of the database
func Snapshot(db *sql.DB, backupPath string) error {
    _, err := db.Exec("VACUUM INTO ?", backupPath)
    return err
}
```

#### Separation of Concerns

+ **`internal/database/database.go`**: Manages the connection lifecycle, low-memory settings, `.flan-keep` verification, integrity checks, and DDL schema creation.
+ **`internal/model/`**: Contains domain-specific query logic (`user.go`, `movie.go`, `series.go`, etc.). Models receive the initialized `*sql.DB` handle and execute business queries, keeping driver/connection concerns cleanly decoupled from application business logic.

---

## 2. Normalized Relational Schema (Evaluated to 4NF)

To resolve attribute mismatches (such as EPUB CFI tracking vs video duration in seconds) and eliminate redundant series title duplication across hundreds of episodes, the database uses a normalized schema:

```sql
-- 1. User Profiles & Lockouts
CREATE TABLE IF NOT EXISTS users (
    user_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    username         TEXT NOT NULL UNIQUE,
    pin_hash         TEXT NOT NULL,
    role             TEXT NOT NULL CHECK(role IN ('admin', 'user')),
    token_version    INTEGER NOT NULL DEFAULT 1, -- Incremented on PIN reset or revocation
    avatar_icon      TEXT DEFAULT 'flan',        -- Curated SVG icon name
    avatar_color     TEXT DEFAULT '#bb9af7',     -- Hex color accent
    failed_attempts  INTEGER DEFAULT 0,
    locked_until     DATETIME,
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Storage Libraries (Multi-Directory Core)
CREATE TABLE IF NOT EXISTS libraries (
    library_id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name             TEXT NOT NULL,
    path             TEXT NOT NULL UNIQUE,
    media_type       TEXT NOT NULL CHECK(media_type IN ('movies', 'tv', 'books')),
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Standalone Movies
CREATE TABLE IF NOT EXISTS movies (
    movie_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id       INTEGER NOT NULL,
    title            TEXT NOT NULL,
    release_year     INTEGER,
    duration_seconds INTEGER DEFAULT 0,
    rating           REAL DEFAULT 0.0,
    overview         TEXT,
    cover_path       TEXT,
    relative_path    TEXT NOT NULL,              -- Path relative to library.path
    file_size        INTEGER NOT NULL,
    format           TEXT NOT NULL,              -- 'mp4', 'webm'
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (library_id) REFERENCES libraries(library_id) ON DELETE CASCADE,
    UNIQUE(library_id, relative_path)
);

-- 4. TV Series
CREATE TABLE IF NOT EXISTS series (
    series_id        INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id       INTEGER NOT NULL,
    title            TEXT NOT NULL,
    release_year     INTEGER,
    rating           REAL DEFAULT 0.0,
    overview         TEXT,
    cover_path       TEXT,                       -- Series poster
    folder_name      TEXT NOT NULL,              -- Series root folder name
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (library_id) REFERENCES libraries(library_id) ON DELETE CASCADE,
    UNIQUE(library_id, folder_name)
);

-- 5. TV Episodes
CREATE TABLE IF NOT EXISTS episodes (
    episode_id       INTEGER PRIMARY KEY AUTOINCREMENT,
    series_id        INTEGER NOT NULL,
    season_number    INTEGER NOT NULL DEFAULT 1,
    episode_number   INTEGER NOT NULL,
    title            TEXT NOT NULL,
    overview         TEXT,
    duration_seconds INTEGER DEFAULT 0,
    relative_path    TEXT NOT NULL,              -- Relative to library.path
    file_size        INTEGER NOT NULL,
    format           TEXT NOT NULL,              -- 'mp4', 'webm'
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES series(series_id) ON DELETE CASCADE,
    UNIQUE(series_id, season_number, episode_number)
);

-- 6. Books & Documents
CREATE TABLE IF NOT EXISTS books (
    book_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id       INTEGER NOT NULL,
    title            TEXT NOT NULL,
    author           TEXT,
    overview         TEXT,
    cover_path       TEXT,
    format           TEXT NOT NULL CHECK(format IN ('epub', 'pdf')),
    relative_path    TEXT NOT NULL,              -- Relative to library.path
    file_size        INTEGER NOT NULL,
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (library_id) REFERENCES libraries(library_id) ON DELETE CASCADE,
    UNIQUE(library_id, relative_path)
);

-- 7. Normalized Genres (1NF / 4NF)
CREATE TABLE IF NOT EXISTS genres (
    genre_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name             TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS item_genres (
    item_type        TEXT NOT NULL CHECK(item_type IN ('movie', 'series', 'book')),
    item_id          INTEGER NOT NULL,
    genre_id         INTEGER NOT NULL,
    PRIMARY KEY (item_type, item_id, genre_id),
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id) ON DELETE CASCADE
);

-- 8. Video Playback Progress (Movies & Episodes)
CREATE TABLE IF NOT EXISTS video_progress (
    user_id          INTEGER NOT NULL,
    video_type       TEXT NOT NULL CHECK(video_type IN ('movie', 'episode')),
    video_id         INTEGER NOT NULL,
    position_seconds INTEGER NOT NULL DEFAULT 0,
    duration_seconds INTEGER NOT NULL DEFAULT 0,
    is_finished      INTEGER NOT NULL DEFAULT 0, -- 0 = in progress, 1 = completed
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, video_type, video_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- 9. Book Reading Progress (EPUB & PDF)
CREATE TABLE IF NOT EXISTS book_progress (
    user_id          INTEGER NOT NULL,
    book_id          INTEGER NOT NULL,
    position_cfi     TEXT,                       -- Canonical Fragment Identifier for EPUBs
    current_page     INTEGER DEFAULT 0,          -- For PDF files
    total_pages      INTEGER DEFAULT 0,          -- For PDF files
    percentage       REAL DEFAULT 0.0,           -- 0.0 to 100.0% completion
    is_finished      INTEGER NOT NULL DEFAULT 0,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, book_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(book_id) ON DELETE CASCADE
);

-- Indexes for Fast Catalog Querying
CREATE INDEX IF NOT EXISTS idx_movies_lib ON movies(library_id);
CREATE INDEX IF NOT EXISTS idx_episodes_series ON episodes(series_id, season_number, episode_number);
CREATE INDEX IF NOT EXISTS idx_books_lib ON books(library_id);
CREATE INDEX IF NOT EXISTS idx_video_progress_user ON video_progress(user_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_book_progress_user ON book_progress(user_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_item_genres_lookup ON item_genres(genre_id, item_type);
```

---

## 3. Core Query Procedures & Access Patterns

### A. Dashboard Shelves (Continue Watching & Continue Reading)

#### Continue Watching (Movies & Episodes Union)

```sql
-- In-progress movies
SELECT 
    'movie' AS video_type,
    m.movie_id AS video_id,
    m.title,
    m.cover_path,
    p.position_seconds,
    p.duration_seconds,
    NULL AS series_title,
    0 AS season_number,
    0 AS episode_number,
    p.updated_at
FROM video_progress p
JOIN movies m ON p.video_id = m.movie_id
WHERE p.user_id = ? AND p.video_type = 'movie' AND p.is_finished = 0 AND p.position_seconds > 10

UNION ALL

-- In-progress TV episodes
SELECT 
    'episode' AS video_type,
    e.episode_id AS video_id,
    e.title,
    s.cover_path,
    p.position_seconds,
    p.duration_seconds,
    s.title AS series_title,
    e.season_number,
    e.episode_number,
    p.updated_at
FROM video_progress p
JOIN episodes e ON p.video_id = e.episode_id
JOIN series s ON e.series_id = s.series_id
WHERE p.user_id = ? AND p.video_type = 'episode' AND p.is_finished = 0 AND p.position_seconds > 10

ORDER BY updated_at DESC
LIMIT 12;
```

#### Continue Reading (Books)

```sql
SELECT 
    b.book_id,
    b.title,
    b.author,
    b.format,
    b.cover_path,
    p.position_cfi,
    p.current_page,
    p.total_pages,
    p.percentage
FROM book_progress p
JOIN books b ON p.book_id = b.book_id
WHERE p.user_id = ? AND p.is_finished = 0 AND (p.percentage > 1.0 OR p.current_page > 1)
ORDER BY p.updated_at DESC
LIMIT 12;
```

---

### B. Catalog Filtering by Genre

#### Movies Filtered by Genre

```sql
SELECT m.movie_id, m.title, m.cover_path, m.release_year, m.rating, m.duration_seconds
FROM movies m
WHERE (? IS NULL OR EXISTS (
    SELECT 1 FROM item_genres ig
    JOIN genres g ON ig.genre_id = g.genre_id
    WHERE ig.item_type = 'movie' AND ig.item_id = m.movie_id AND g.name = ?
))
ORDER BY m.title ASC;
```

#### TV Series Filtered by Genre

```sql
SELECT s.series_id, s.title, s.cover_path, s.release_year, s.rating,
       COUNT(e.episode_id) AS episode_count
FROM series s
LEFT JOIN episodes e ON s.series_id = e.series_id
WHERE (? IS NULL OR EXISTS (
    SELECT 1 FROM item_genres ig
    JOIN genres g ON ig.genre_id = g.genre_id
    WHERE ig.item_type = 'series' AND ig.item_id = s.series_id AND g.name = ?
))
GROUP BY s.series_id
ORDER BY s.title ASC;
```

#### Books Filtered by Genre or Format

```sql
SELECT b.book_id, b.title, b.author, b.format, b.cover_path
FROM books b
WHERE (? IS NULL OR b.format = ?)
  AND (? IS NULL OR EXISTS (
    SELECT 1 FROM item_genres ig
    JOIN genres g ON ig.genre_id = g.genre_id
    WHERE ig.item_type = 'book' AND ig.item_id = b.book_id AND g.name = ?
))
ORDER BY b.title ASC;
```

---

### C. TV Series Detail View (Seasons & Episodes)

Retrieves all episodes for a specific series alongside watch progress for the current user:

```sql
SELECT 
    e.episode_id,
    e.season_number,
    e.episode_number,
    e.title,
    e.overview,
    e.duration_seconds,
    COALESCE(p.position_seconds, 0) AS position_seconds,
    COALESCE(p.is_finished, 0) AS is_finished
FROM episodes e
LEFT JOIN video_progress p ON e.episode_id = p.video_id AND p.video_type = 'episode' AND p.user_id = ?
WHERE e.series_id = ?
ORDER BY e.season_number ASC, e.episode_number ASC;
```

---

### D. Progress Upserts

#### Video Progress Upsert

```sql
INSERT INTO video_progress (user_id, video_type, video_id, position_seconds, duration_seconds, is_finished, updated_at)
VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, video_type, video_id) DO UPDATE SET
    position_seconds = excluded.position_seconds,
    duration_seconds = excluded.duration_seconds,
    is_finished = excluded.is_finished,
    updated_at = CURRENT_TIMESTAMP;
```

#### Book Progress Upsert

```sql
INSERT INTO book_progress (user_id, book_id, position_cfi, current_page, total_pages, percentage, is_finished, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, book_id) DO UPDATE SET
    position_cfi = excluded.position_cfi,
    current_page = excluded.current_page,
    total_pages = excluded.total_pages,
    percentage = excluded.percentage,
    is_finished = excluded.is_finished,
    updated_at = CURRENT_TIMESTAMP;
```

---

### E. PIN Authentication & Brute-Force Rate Limiting

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
SET pin_hash = ?, failed_attempts = 0, locked_until = NULL, token_version = token_version + 1
WHERE role = 'admin';

-- 5. Change or reset user PIN via Settings or Admin API (invalidates existing sessions)
UPDATE users
SET pin_hash = ?, failed_attempts = 0, locked_until = NULL, token_version = token_version + 1
WHERE user_id = ?;
```

---

### F. Library Scanning & Ingestion

```sql
-- Check if file already exists in library before scanning
SELECT movie_id FROM movies WHERE library_id = ? AND relative_path = ?;

-- Insert newly discovered movie
INSERT INTO movies (
    library_id, title, release_year, duration_seconds, rating,
    overview, cover_path, relative_path, file_size, format
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Clean up removed movie
DELETE FROM movies WHERE library_id = ? AND relative_path = ?;
```

---

### G. Database Integrity & Zero-Lock Snapshots

```sql
-- 1. Fast integrity check executed on startup
PRAGMA quick_check;

-- 2. Zero-lock hot snapshot executed daily in background
VACUUM INTO 'data/flan.db.backup';
```

---

### Related Documentation

+ [Master System Specifications](design.md)
+ [Storage Architecture & Mount Resiliency](storage.md)
+ [Testing Strategy & TDD Decoupling](testing.md)
+ [Security Threat Model & Session Hardening](threat-model.md)
+ [Cross-Compilation & SBC Deployment Guide](compilation.md)
