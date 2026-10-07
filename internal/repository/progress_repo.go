package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
)

type ProgressRepository struct {
	db *database.DB
}

func NewProgressRepository(db *database.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) Get(ctx context.Context, userID int64, mediaType model.MediaType, fileID int64) (*model.Progress, error) {
	query := `
	SELECT user_id, media_type, file_id, position_data, percentage, is_finished, updated_at
	FROM progress
	WHERE user_id = ? AND media_type = ? AND file_id = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, userID, mediaType, fileID)
	var p model.Progress
	var isFinishedInt int
	err := row.Scan(&p.UserID, &p.MediaType, &p.FileID, &p.PositionData, &p.Percentage, &isFinishedInt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	p.IsFinished = isFinishedInt == 1
	return &p, nil
}

func (r *ProgressRepository) Upsert(ctx context.Context, p *model.Progress) error {
	query := `
	INSERT INTO progress (user_id, media_type, file_id, position_data, percentage, is_finished, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(user_id, media_type, file_id) DO UPDATE SET
	    position_data = excluded.position_data,
	    percentage = excluded.percentage,
	    is_finished = excluded.is_finished,
	    updated_at = CURRENT_TIMESTAMP;`

	isFinishedInt := 0
	if p.IsFinished {
		isFinishedInt = 1
	}

	_, err := r.db.Writer.ExecContext(ctx, query, p.UserID, p.MediaType, p.FileID, p.PositionData, p.Percentage, isFinishedInt)
	return err
}

func (r *ProgressRepository) ListContinueWatching(ctx context.Context, userID int64, limit int) ([]*model.ContinueWatchingItem, error) {
	if limit <= 0 {
		limit = 10
	}
	query := `
	SELECT 
	    v.video_id, f.file_id, v.title, 
	    COALESCE(NULLIF(f.custom_title, ''), f.title) AS subtitle,
	    v.cover_path, p.position_data, p.percentage, p.updated_at
	FROM progress p
	JOIN video_files f ON p.file_id = f.file_id
	JOIN videos v ON f.video_id = v.video_id
	JOIN storage_sources s ON v.source_id = s.source_id
	WHERE p.user_id = ? 
	  AND p.media_type = 'video' 
	  AND p.is_finished = 0 
	  AND p.percentage > 1.0 
	  AND p.percentage < 95.0
	  AND f.is_missing = 0
	  AND v.is_hidden = 0
	  AND s.is_active = 1
	ORDER BY p.updated_at DESC
	LIMIT ?;`

	rows, err := r.db.Reader.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.ContinueWatchingItem
	for rows.Next() {
		var item model.ContinueWatchingItem
		var coverPath sql.NullString
		var posDataStr string

		err := rows.Scan(
			&item.MediaID, &item.FileID, &item.Title, &item.Subtitle,
			&coverPath, &posDataStr, &item.Percentage, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		item.MediaType = model.MediaTypeVideo
		item.CoverPath = coverPath.String
		item.StatusTag = "In Progress"

		if posSec, err := strconv.ParseFloat(posDataStr, 64); err == nil {
			item.PositionSeconds = posSec
			// Format pos as MM:SS or HH:MM:SS
			secs := int(posSec)
			h := secs / 3600
			m := (secs % 3600) / 60
			s := secs % 60
			if h > 0 {
				item.FormattedPos = fmt.Sprintf("%d:%02d:%02d", h, m, s)
			} else {
				item.FormattedPos = fmt.Sprintf("%02d:%02d", m, s)
			}
		}

		items = append(items, &item)
	}
	return items, rows.Err()
}

func (r *ProgressRepository) ListJumpBackInBooks(ctx context.Context, userID int64, limit int) ([]*model.ContinueWatchingItem, error) {
	if limit <= 0 {
		limit = 10
	}
	query := `
	SELECT 
	    b.book_id, f.file_id, b.title, 
	    COALESCE(NULLIF(b.author, ''), COALESCE(NULLIF(f.custom_title, ''), f.title)) AS subtitle,
	    b.cover_path, p.position_data, p.percentage, p.updated_at
	FROM progress p
	JOIN book_files f ON p.file_id = f.file_id
	JOIN books b ON f.book_id = b.book_id
	JOIN storage_sources s ON b.source_id = s.source_id
	WHERE p.user_id = ? 
	  AND p.media_type = 'book' 
	  AND p.position_data = 'reading'
	  AND p.is_finished = 0 
	  AND f.is_missing = 0
	  AND b.is_hidden = 0
	  AND s.is_active = 1
	ORDER BY p.updated_at DESC
	LIMIT ?;`

	rows, err := r.db.Reader.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.ContinueWatchingItem
	for rows.Next() {
		var item model.ContinueWatchingItem
		var coverPath sql.NullString
		var posDataStr string

		err := rows.Scan(
			&item.MediaID, &item.FileID, &item.Title, &item.Subtitle,
			&coverPath, &posDataStr, &item.Percentage, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		item.MediaType = model.MediaTypeBook
		item.CoverPath = coverPath.String
		item.StatusTag = "Reading"
		item.FormattedPos = "Reading"

		items = append(items, &item)
	}
	return items, rows.Err()
}
