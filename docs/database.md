# Database Schema & Query Procedures

Flan Media Server uses embedded SQLite3 via Go's `database/sql` and the pure-Go driver `modernc.org/sqlite` (no CGo, instant cross-compilation).

---

## 1. Storage Location & Wear-Leveling Pragmas

* **Database File (`DB_PATH`):** Defaults to `./data/flan.db`. Recommended on fast flash storage (eMMC, NVMe, or root micro-SD) even when media sits on external spinning USB hard drives.
* **Ghost Mount Defense (`.flan-keep`):** Before initializing an external DB path, Flan verifies that `.flan-keep` exists in the folder. If missing, startup halts immediately to prevent writing a blank database to an unmounted root micro-SD card.

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

## 2. Relational Schema

```sql
-- 1. User Profiles & Lockouts
CREATE TABLE IF NOT EXISTS users (
    user_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username        TEXT NOT NULL UNIQUE,
    pin_hash        TEXT NOT NULL,
    role            TEXT NOT NULL CHECK(role IN ('admin', 'user')),
    token_version   INTEGER NOT NULL DEFAULT 1,
    avatar_icon     TEXT DEFAULT 'flan',
    avatar_color    TEXT DEFAULT '#bb9af7',
    failed_attempts INTEGER DEFAULT 0,
    locked_until    DATETIME,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Storage Libraries (Multi-Directory Roots)
CREATE TABLE IF NOT EXISTS libraries (
    library_id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    path            TEXT NOT NULL UNIQUE,
    media_type      TEXT NOT NULL CHECK(media_type IN ('movies', 'tv', 'books')),
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Standalone Movies
CREATE TABLE IF NOT EXISTS movies (
    movie_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id       INTEGER NOT NULL REFERENCES libraries(library_id) ON DELETE CASCADE,
    title            TEXT NOT NULL,
    release_year     INTEGER,
    duration_seconds INTEGER DEFAULT 0,
    rating           REAL DEFAULT 0.0,
    overview         TEXT,
    cover_path       TEXT,
    relative_path    TEXT NOT NULL,
    file_size        INTEGER NOT NULL,
    format           TEXT NOT NULL,
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(library_id, relative_path)
);

-- 4. TV Series & Episodes Hierarchy
CREATE TABLE IF NOT EXISTS series (
    series_id        INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id       INTEGER NOT NULL REFERENCES libraries(library_id) ON DELETE CASCADE,
    title            TEXT NOT NULL,
    release_year     INTEGER,
    rating           REAL DEFAULT 0.0,
    overview         TEXT,
    cover_path       TEXT,
    folder_name      TEXT NOT NULL,
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(library_id, folder_name)
);

CREATE TABLE IF NOT EXISTS episodes (
    episode_id       INTEGER PRIMARY KEY AUTOINCREMENT,
    series_id        INTEGER NOT NULL REFERENCES series(series_id) ON DELETE CASCADE,
    season_number    INTEGER NOT NULL DEFAULT 1,
    episode_number   INTEGER NOT NULL,
    title            TEXT NOT NULL,
    overview         TEXT,
    duration_seconds INTEGER DEFAULT 0,
    relative_path    TEXT NOT NULL,
    file_size        INTEGER NOT NULL,
    format           TEXT NOT NULL,
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(series_id, season_number, episode_number)
);

-- 5. Books & Documents
CREATE TABLE IF NOT EXISTS books (
    book_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id       INTEGER NOT NULL REFERENCES libraries(library_id) ON DELETE CASCADE,
    title            TEXT NOT NULL,
    author           TEXT,
    overview         TEXT,
    cover_path       TEXT,
    format           TEXT NOT NULL CHECK(format IN ('epub', 'pdf')),
    relative_path    TEXT NOT NULL,
    file_size        INTEGER NOT NULL,
    added_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(library_id, relative_path)
);

-- 6. Normalized Genres
CREATE TABLE IF NOT EXISTS genres (
    genre_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name             TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS item_genres (
    item_type        TEXT NOT NULL CHECK(item_type IN ('movie', 'series', 'book')),
    item_id          INTEGER NOT NULL,
    genre_id         INTEGER NOT NULL REFERENCES genres(genre_id) ON DELETE CASCADE,
    PRIMARY KEY (item_type, item_id, genre_id)
);

-- 7. Playback & Reading Progress Stores
CREATE TABLE IF NOT EXISTS video_progress (
    user_id          INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    video_type       TEXT NOT NULL CHECK(video_type IN ('movie', 'episode')),
    video_id         INTEGER NOT NULL,
    position_seconds INTEGER NOT NULL DEFAULT 0,
    duration_seconds INTEGER NOT NULL DEFAULT 0,
    is_finished      INTEGER NOT NULL DEFAULT 0,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, video_type, video_id)
);

CREATE TABLE IF NOT EXISTS book_progress (
    user_id          INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    book_id          INTEGER NOT NULL REFERENCES books(book_id) ON DELETE CASCADE,
    position_cfi     TEXT,
    current_page     INTEGER DEFAULT 0,
    total_pages      INTEGER DEFAULT 0,
    percentage       REAL DEFAULT 0.0,
    is_finished      INTEGER NOT NULL DEFAULT 0,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, book_id)
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

## 3. Core Query Patterns

### Dashboard "Continue Watching" Shelf
```sql
SELECT 'movie' AS video_type, m.movie_id AS video_id, m.title, m.cover_path,
       p.position_seconds, p.duration_seconds, NULL AS series_title, 0 AS season_number, 0 AS episode_number, p.updated_at
FROM video_progress p JOIN movies m ON p.video_id = m.movie_id
WHERE p.user_id = ? AND p.video_type = 'movie' AND p.is_finished = 0 AND p.position_seconds > 10
UNION ALL
SELECT 'episode' AS video_type, e.episode_id AS video_id, e.title, s.cover_path,
       p.position_seconds, p.duration_seconds, s.title AS series_title, e.season_number, e.episode_number, p.updated_at
FROM video_progress p JOIN episodes e ON p.video_id = e.episode_id JOIN series s ON e.series_id = s.series_id
WHERE p.user_id = ? AND p.video_type = 'episode' AND p.is_finished = 0 AND p.position_seconds > 10
ORDER BY updated_at DESC LIMIT 12;
```

### Atomic Progress Upsert
```sql
INSERT INTO video_progress (user_id, video_type, video_id, position_seconds, duration_seconds, is_finished, updated_at)
VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, video_type, video_id) DO UPDATE SET
    position_seconds = excluded.position_seconds,
    duration_seconds = excluded.duration_seconds,
    is_finished = excluded.is_finished,
    updated_at = CURRENT_TIMESTAMP;
```

### PIN Lockout Management
```sql
-- Lock account for 5 minutes after 5 consecutive failures
UPDATE users
SET failed_attempts = failed_attempts + 1,
    locked_until = CASE WHEN failed_attempts + 1 >= 5 THEN datetime(CURRENT_TIMESTAMP, '+5 minutes') ELSE locked_until END
WHERE user_id = ?;

-- Reset counter on successful login
UPDATE users SET failed_attempts = 0, locked_until = NULL WHERE user_id = ?;
```

### Zero-Lock Backups & Health Checks
* **Integrity Check (Boot):** `PRAGMA quick_check;` validates file consistency on startup.
* **Daily Hot Backup:** `VACUUM INTO 'data/flan.db.backup';` takes an atomic non-blocking snapshot every 24 hours.

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Storage Architecture & Multi-Drive Resilience](storage.md)
* [Security Threat Model](threat-model.md)
