package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
)

type VideoRepository struct {
	db *database.DB
}

func NewVideoRepository(db *database.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) List(ctx context.Context, filterType string, showHidden bool) ([]*model.Video, error) {
	query := `
	SELECT v.video_id, v.source_id, v.title, v.video_type, v.release_year, v.overview, v.cover_path, v.folder_path, v.is_hidden, v.metadata_locked, v.created_at, s.name,
	       (SELECT COUNT(*) FROM video_files vf WHERE vf.video_id = v.video_id AND vf.is_missing = 0) AS active_files
	FROM videos v
	JOIN storage_sources s ON v.source_id = s.source_id
	WHERE (v.is_hidden = 0 OR ?)
	  AND s.is_active = 1
	  AND (? = '' OR ? = 'all' OR v.video_type = ?)
	  AND EXISTS (SELECT 1 FROM video_files vf WHERE vf.video_id = v.video_id AND vf.is_missing = 0)
	ORDER BY v.title ASC;`

	rows, err := r.db.Reader.QueryContext(ctx, query, showHidden, filterType, filterType, filterType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		var v model.Video
		var releaseYear sql.NullInt64
		var overview, coverPath sql.NullString
		var isHiddenInt, metaLockedInt int

		err := rows.Scan(
			&v.ID, &v.SourceID, &v.Title, &v.VideoType, &releaseYear, &overview, &coverPath, &v.FolderPath,
			&isHiddenInt, &metaLockedInt, &v.CreatedAt, &v.SourceName, &v.ActiveFiles,
		)
		if err != nil {
			return nil, err
		}
		if releaseYear.Valid {
			v.ReleaseYear = int(releaseYear.Int64)
		}
		v.Overview = overview.String
		v.CoverPath = coverPath.String
		v.IsHidden = isHiddenInt == 1
		v.MetadataLocked = metaLockedInt == 1
		videos = append(videos, &v)
	}
	return videos, rows.Err()
}

func (r *VideoRepository) Search(ctx context.Context, queryStr string, showHidden bool) ([]*model.Video, error) {
	likePattern := "%" + queryStr + "%"
	query := `
	SELECT v.video_id, v.source_id, v.title, v.video_type, v.release_year, v.overview, v.cover_path, v.folder_path, v.is_hidden, v.metadata_locked, v.created_at, s.name,
	       (SELECT COUNT(*) FROM video_files vf WHERE vf.video_id = v.video_id AND vf.is_missing = 0) AS active_files
	FROM videos v
	JOIN storage_sources s ON v.source_id = s.source_id
	WHERE (v.is_hidden = 0 OR ?)
	  AND s.is_active = 1
	  AND (v.title LIKE ? OR v.overview LIKE ?)
	ORDER BY v.title ASC;`

	rows, err := r.db.Reader.QueryContext(ctx, query, showHidden, likePattern, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		var v model.Video
		var releaseYear sql.NullInt64
		var overview, coverPath sql.NullString
		var isHiddenInt, metaLockedInt int

		err := rows.Scan(
			&v.ID, &v.SourceID, &v.Title, &v.VideoType, &releaseYear, &overview, &coverPath, &v.FolderPath,
			&isHiddenInt, &metaLockedInt, &v.CreatedAt, &v.SourceName, &v.ActiveFiles,
		)
		if err != nil {
			return nil, err
		}
		if releaseYear.Valid {
			v.ReleaseYear = int(releaseYear.Int64)
		}
		v.Overview = overview.String
		v.CoverPath = coverPath.String
		v.IsHidden = isHiddenInt == 1
		v.MetadataLocked = metaLockedInt == 1
		videos = append(videos, &v)
	}
	return videos, rows.Err()
}

func (r *VideoRepository) GetByID(ctx context.Context, id int64) (*model.Video, error) {
	query := `
	SELECT v.video_id, v.source_id, v.title, v.video_type, v.release_year, v.overview, v.cover_path, v.folder_path, v.is_hidden, v.metadata_locked, v.created_at, s.name
	FROM videos v
	JOIN storage_sources s ON v.source_id = s.source_id
	WHERE v.video_id = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, id)
	var v model.Video
	var releaseYear sql.NullInt64
	var overview, coverPath sql.NullString
	var isHiddenInt, metaLockedInt int

	err := row.Scan(
		&v.ID, &v.SourceID, &v.Title, &v.VideoType, &releaseYear, &overview, &coverPath, &v.FolderPath,
		&isHiddenInt, &metaLockedInt, &v.CreatedAt, &v.SourceName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	if releaseYear.Valid {
		v.ReleaseYear = int(releaseYear.Int64)
	}
	v.Overview = overview.String
	v.CoverPath = coverPath.String
	v.IsHidden = isHiddenInt == 1
	v.MetadataLocked = metaLockedInt == 1
	return &v, nil
}

func (r *VideoRepository) GetBySourceAndFolder(ctx context.Context, sourceID int64, folderPath string) (*model.Video, error) {
	query := `
	SELECT v.video_id, v.source_id, v.title, v.video_type, v.release_year, v.overview, v.cover_path, v.folder_path, v.is_hidden, v.metadata_locked, v.created_at
	FROM videos v
	WHERE v.source_id = ? AND v.folder_path = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, sourceID, folderPath)
	var v model.Video
	var releaseYear sql.NullInt64
	var overview, coverPath sql.NullString
	var isHiddenInt, metaLockedInt int

	err := row.Scan(
		&v.ID, &v.SourceID, &v.Title, &v.VideoType, &releaseYear, &overview, &coverPath, &v.FolderPath,
		&isHiddenInt, &metaLockedInt, &v.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	if releaseYear.Valid {
		v.ReleaseYear = int(releaseYear.Int64)
	}
	v.Overview = overview.String
	v.CoverPath = coverPath.String
	v.IsHidden = isHiddenInt == 1
	v.MetadataLocked = metaLockedInt == 1
	return &v, nil
}

func (r *VideoRepository) Create(ctx context.Context, v *model.Video) (int64, error) {
	query := `
	INSERT INTO videos (source_id, title, video_type, release_year, overview, cover_path, folder_path, is_hidden, metadata_locked)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	isHiddenInt := 0
	if v.IsHidden {
		isHiddenInt = 1
	}
	metaLockedInt := 0
	if v.MetadataLocked {
		metaLockedInt = 1
	}

	res, err := r.db.Writer.ExecContext(ctx, query, v.SourceID, v.Title, v.VideoType, v.ReleaseYear, v.Overview, v.CoverPath, v.FolderPath, isHiddenInt, metaLockedInt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *VideoRepository) UpdateMetadata(ctx context.Context, v *model.Video) error {
	query := `
	UPDATE videos 
	SET title = ?, release_year = ?, video_type = ?, overview = ?, cover_path = ?, metadata_locked = 1
	WHERE video_id = ?;`

	_, err := r.db.Writer.ExecContext(ctx, query, v.Title, v.ReleaseYear, v.VideoType, v.Overview, v.CoverPath, v.ID)
	return err
}

func (r *VideoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM videos WHERE video_id = ?;`
	res, err := r.db.Writer.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *VideoRepository) ListFilesByVideoID(ctx context.Context, videoID int64, userID int64, showHidden bool) ([]*model.VideoFile, error) {
	query := `
	SELECT 
	    f.file_id, f.video_id, f.title, f.custom_title, f.relative_path, f.file_size, f.mtime, f.format,
	    f.duration_seconds, f.season_number, f.episode_number, f.order_index, f.is_hidden, f.is_missing,
	    COALESCE(p.position_data, '0') AS position_data,
	    COALESCE(p.percentage, 0.0) AS percentage,
	    COALESCE(p.is_finished, 0) AS is_finished
	FROM video_files f
	LEFT JOIN progress p ON f.file_id = p.file_id AND p.media_type = 'video' AND p.user_id = ?
	WHERE f.video_id = ? AND (f.is_hidden = 0 OR ?)
	ORDER BY f.order_index ASC, f.title ASC;`

	rows, err := r.db.Reader.QueryContext(ctx, query, userID, videoID, showHidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*model.VideoFile
	for rows.Next() {
		var f model.VideoFile
		var customTitle sql.NullString
		var isHiddenInt, isMissingInt, isFinishedInt int
		var posDataStr string

		err := rows.Scan(
			&f.ID, &f.VideoID, &f.Title, &customTitle, &f.RelativePath, &f.FileSize, &f.MTime, &f.Format,
			&f.DurationSeconds, &f.SeasonNumber, &f.EpisodeNumber, &f.OrderIndex, &isHiddenInt, &isMissingInt,
			&posDataStr, &f.Percentage, &isFinishedInt,
		)
		if err != nil {
			return nil, err
		}
		f.CustomTitle = customTitle.String
		f.IsHidden = isHiddenInt == 1
		f.IsMissing = isMissingInt == 1
		f.IsFinished = isFinishedInt == 1
		if posSec, err := strconv.ParseFloat(posDataStr, 64); err == nil {
			f.PositionSeconds = posSec
		}
		files = append(files, &f)
	}
	return files, rows.Err()
}

func (r *VideoRepository) GetFileByID(ctx context.Context, fileID int64) (*model.VideoFile, error) {
	query := `
	SELECT file_id, video_id, title, custom_title, relative_path, file_size, mtime, format,
	       duration_seconds, season_number, episode_number, order_index, is_hidden, is_missing
	FROM video_files
	WHERE file_id = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, fileID)
	var f model.VideoFile
	var customTitle sql.NullString
	var isHiddenInt, isMissingInt int

	err := row.Scan(
		&f.ID, &f.VideoID, &f.Title, &customTitle, &f.RelativePath, &f.FileSize, &f.MTime, &f.Format,
		&f.DurationSeconds, &f.SeasonNumber, &f.EpisodeNumber, &f.OrderIndex, &isHiddenInt, &isMissingInt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	f.CustomTitle = customTitle.String
	f.IsHidden = isHiddenInt == 1
	f.IsMissing = isMissingInt == 1
	return &f, nil
}

func (r *VideoRepository) GetFileWithContainer(ctx context.Context, fileID int64) (*model.VideoFile, *model.Video, *model.StorageSource, error) {
	query := `
	SELECT 
	    f.file_id, f.video_id, f.title, f.custom_title, f.relative_path, f.file_size, f.mtime, f.format,
	    f.duration_seconds, f.season_number, f.episode_number, f.order_index, f.is_hidden, f.is_missing,
	    v.video_id, v.source_id, v.title, v.video_type, v.folder_path,
	    s.source_id, s.name, s.media_type, s.folder_path, s.is_active
	FROM video_files f
	JOIN videos v ON f.video_id = v.video_id
	JOIN storage_sources s ON v.source_id = s.source_id
	WHERE f.file_id = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, fileID)
	var f model.VideoFile
	var v model.Video
	var s model.StorageSource
	var customTitle sql.NullString
	var fHidden, fMissing, sActive int

	err := row.Scan(
		&f.ID, &f.VideoID, &f.Title, &customTitle, &f.RelativePath, &f.FileSize, &f.MTime, &f.Format,
		&f.DurationSeconds, &f.SeasonNumber, &f.EpisodeNumber, &f.OrderIndex, &fHidden, &fMissing,
		&v.ID, &v.SourceID, &v.Title, &v.VideoType, &v.FolderPath,
		&s.ID, &s.Name, &s.MediaType, &s.FolderPath, &sActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil, model.ErrNotFound
		}
		return nil, nil, nil, err
	}
	f.CustomTitle = customTitle.String
	f.IsHidden = fHidden == 1
	f.IsMissing = fMissing == 1
	s.IsActive = sActive == 1
	return &f, &v, &s, nil
}

func (r *VideoRepository) CreateFile(ctx context.Context, f *model.VideoFile) (int64, error) {
	query := `
	INSERT INTO video_files (video_id, title, custom_title, relative_path, file_size, mtime, format, duration_seconds, season_number, episode_number, order_index, is_hidden, is_missing)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	res, err := r.db.Writer.ExecContext(ctx, query,
		f.VideoID, f.Title, f.CustomTitle, f.RelativePath, f.FileSize, f.MTime, f.Format,
		f.DurationSeconds, f.SeasonNumber, f.EpisodeNumber, f.OrderIndex, 0, 0,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *VideoRepository) UpdateFile(ctx context.Context, f *model.VideoFile) error {
	query := `
	UPDATE video_files
	SET custom_title = ?, order_index = ?, is_hidden = ?
	WHERE file_id = ?;`

	isHiddenInt := 0
	if f.IsHidden {
		isHiddenInt = 1
	}
	_, err := r.db.Writer.ExecContext(ctx, query, f.CustomTitle, f.OrderIndex, isHiddenInt, f.ID)
	return err
}

func (r *VideoRepository) BatchUpdateFiles(ctx context.Context, files []*model.VideoFile) error {
	tx, err := r.db.Writer.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `UPDATE video_files SET custom_title = ?, order_index = ?, is_hidden = ? WHERE file_id = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, f := range files {
		isHiddenInt := 0
		if f.IsHidden {
			isHiddenInt = 1
		}
		if _, err := stmt.ExecContext(ctx, f.CustomTitle, f.OrderIndex, isHiddenInt, f.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *VideoRepository) SetFileMissing(ctx context.Context, fileID int64, missing bool) error {
	query := `UPDATE video_files SET is_missing = ? WHERE file_id = ?;`
	val := 0
	if missing {
		val = 1
	}
	_, err := r.db.Writer.ExecContext(ctx, query, val, fileID)
	return err
}

func (r *VideoRepository) MarkAllFilesMissingForSource(ctx context.Context, sourceID int64, missing bool) error {
	query := `
	UPDATE video_files 
	SET is_missing = ? 
	WHERE video_id IN (SELECT video_id FROM videos WHERE source_id = ?);`
	val := 0
	if missing {
		val = 1
	}
	_, err := r.db.Writer.ExecContext(ctx, query, val, sourceID)
	return err
}
