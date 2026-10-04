package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
)

type SourceRepository struct {
	db *database.DB
}

func NewSourceRepository(db *database.DB) *SourceRepository {
	return &SourceRepository{db: db}
}

func (r *SourceRepository) ListAll(ctx context.Context) ([]*model.StorageSource, error) {
	query := `SELECT source_id, name, media_type, folder_path, is_active, created_at FROM storage_sources ORDER BY name ASC`
	rows, err := r.db.Reader.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []*model.StorageSource
	for rows.Next() {
		var s model.StorageSource
		var isActiveInt int
		if err := rows.Scan(&s.ID, &s.Name, &s.MediaType, &s.FolderPath, &isActiveInt, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.IsActive = isActiveInt == 1
		sources = append(sources, &s)
	}
	return sources, rows.Err()
}

func (r *SourceRepository) GetByID(ctx context.Context, id int64) (*model.StorageSource, error) {
	query := `SELECT source_id, name, media_type, folder_path, is_active, created_at FROM storage_sources WHERE source_id = ?`
	row := r.db.Reader.QueryRowContext(ctx, query, id)

	var s model.StorageSource
	var isActiveInt int
	err := row.Scan(&s.ID, &s.Name, &s.MediaType, &s.FolderPath, &isActiveInt, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	s.IsActive = isActiveInt == 1
	return &s, nil
}

func (r *SourceRepository) GetByPath(ctx context.Context, folderPath string) (*model.StorageSource, error) {
	query := `SELECT source_id, name, media_type, folder_path, is_active, created_at FROM storage_sources WHERE folder_path = ?`
	row := r.db.Reader.QueryRowContext(ctx, query, folderPath)

	var s model.StorageSource
	var isActiveInt int
	err := row.Scan(&s.ID, &s.Name, &s.MediaType, &s.FolderPath, &isActiveInt, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	s.IsActive = isActiveInt == 1
	return &s, nil
}

func (r *SourceRepository) Create(ctx context.Context, s *model.StorageSource) (int64, error) {
	query := `INSERT INTO storage_sources (name, media_type, folder_path, is_active) VALUES (?, ?, ?, ?)`
	isActiveInt := 0
	if s.IsActive {
		isActiveInt = 1
	}
	res, err := r.db.Writer.ExecContext(ctx, query, s.Name, s.MediaType, s.FolderPath, isActiveInt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SourceRepository) UpdateName(ctx context.Context, id int64, name string) error {
	query := `UPDATE storage_sources SET name = ? WHERE source_id = ?`
	res, err := r.db.Writer.ExecContext(ctx, query, name, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *SourceRepository) SetActive(ctx context.Context, id int64, active bool) error {
	query := `UPDATE storage_sources SET is_active = ? WHERE source_id = ?`
	val := 0
	if active {
		val = 1
	}
	_, err := r.db.Writer.ExecContext(ctx, query, val, id)
	return err
}

func (r *SourceRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM storage_sources WHERE source_id = ?`
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
