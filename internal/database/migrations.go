package database

import (
	"database/sql"
	"fmt"
)

const TargetSchemaVersion = 1

const migrationV1 = `
-- 1. User Profiles
CREATE TABLE IF NOT EXISTS users (
    user_id         INTEGER PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    display_name    TEXT,
    pin_hash        TEXT NOT NULL,
    role            TEXT NOT NULL CHECK(role IN ('admin', 'user')),
    token_version   INTEGER NOT NULL DEFAULT 1,
    avatar_icon     TEXT DEFAULT 'default',
    avatar_path     TEXT,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Configurable Storage Sources
CREATE TABLE IF NOT EXISTS storage_sources (
    source_id       INTEGER PRIMARY KEY,
    name            TEXT NOT NULL,
    media_type      TEXT NOT NULL CHECK(media_type IN ('video', 'book')),
    folder_path     TEXT NOT NULL UNIQUE,
    is_active       INTEGER NOT NULL DEFAULT 1,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Video Containers
CREATE TABLE IF NOT EXISTS videos (
    video_id        INTEGER PRIMARY KEY,
    source_id       INTEGER NOT NULL REFERENCES storage_sources(source_id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    video_type      TEXT NOT NULL DEFAULT 'series' CHECK(video_type IN ('movie', 'series')),
    release_year    INTEGER,
    overview        TEXT,
    cover_path      TEXT,
    folder_path     TEXT NOT NULL,
    is_hidden       INTEGER NOT NULL DEFAULT 0,
    metadata_locked INTEGER NOT NULL DEFAULT 1,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_id, folder_path)
);

-- 4. Playable Video Files
CREATE TABLE IF NOT EXISTS video_files (
    file_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    video_id         INTEGER NOT NULL REFERENCES videos(video_id) ON DELETE CASCADE,
    title            TEXT NOT NULL,
    custom_title     TEXT,
    relative_path    TEXT NOT NULL,
    file_size        INTEGER NOT NULL,
    mtime            INTEGER DEFAULT 0,
    format           TEXT NOT NULL,
    duration_seconds INTEGER DEFAULT 0,
    season_number    INTEGER DEFAULT 1,
    episode_number   INTEGER DEFAULT 0,
    order_index      INTEGER DEFAULT 0,
    is_hidden        INTEGER NOT NULL DEFAULT 0,
    is_missing       INTEGER NOT NULL DEFAULT 0,
    UNIQUE(video_id, relative_path)
);

-- 5. Book Containers
CREATE TABLE IF NOT EXISTS books (
    book_id         INTEGER PRIMARY KEY,
    source_id       INTEGER NOT NULL REFERENCES storage_sources(source_id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    author          TEXT,
    overview        TEXT,
    cover_path      TEXT,
    folder_path     TEXT NOT NULL,
    is_hidden       INTEGER NOT NULL DEFAULT 0,
    metadata_locked INTEGER NOT NULL DEFAULT 1,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_id, folder_path)
);

-- 6. Book Files
CREATE TABLE IF NOT EXISTS book_files (
    file_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id          INTEGER NOT NULL REFERENCES books(book_id) ON DELETE CASCADE,
    title            TEXT NOT NULL,
    custom_title     TEXT,
    relative_path    TEXT NOT NULL,
    file_size        INTEGER NOT NULL,
    mtime            INTEGER DEFAULT 0,
    format           TEXT NOT NULL CHECK(format IN ('epub', 'pdf')),
    order_index      INTEGER DEFAULT 0,
    is_hidden        INTEGER NOT NULL DEFAULT 0,
    is_missing       INTEGER NOT NULL DEFAULT 0,
    UNIQUE(book_id, relative_path)
);

-- 7. Unified Progress Tracking
CREATE TABLE IF NOT EXISTS progress (
    user_id          INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    media_type       TEXT NOT NULL CHECK(media_type IN ('video', 'book')),
    file_id          INTEGER NOT NULL,
    position_data    TEXT NOT NULL,
    percentage       REAL DEFAULT 0.0,
    is_finished      INTEGER NOT NULL DEFAULT 0,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, media_type, file_id)
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_videos_source ON videos(source_id, video_type, is_hidden);
CREATE INDEX IF NOT EXISTS idx_books_source ON books(source_id, is_hidden);
CREATE INDEX IF NOT EXISTS idx_video_files_video ON video_files(video_id, order_index);
CREATE INDEX IF NOT EXISTS idx_book_files_book ON book_files(book_id, order_index);
CREATE INDEX IF NOT EXISTS idx_progress_lookup ON progress(user_id, media_type, updated_at DESC);

-- Automated Progress Cleanup Triggers
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
`

func Migrate(db *sql.DB) error {
	var currentVersion int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&currentVersion); err != nil {
		return fmt.Errorf("read user_version: %w", err)
	}

	if currentVersion < 1 {
		if _, err := db.Exec(migrationV1); err != nil {
			return fmt.Errorf("apply migration v1: %w", err)
		}
		if _, err := db.Exec("PRAGMA user_version = 1;"); err != nil {
			return fmt.Errorf("set user_version = 1: %w", err)
		}
	}

	return nil
}
