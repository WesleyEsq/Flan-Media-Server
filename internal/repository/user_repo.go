package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.Reader.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT user_id, username, display_name, pin_hash, role, token_version, avatar_icon, avatar_path, created_at FROM users WHERE user_id = ?`
	row := r.db.Reader.QueryRowContext(ctx, query, id)

	var u model.User
	var displayName, avatarIcon, avatarPath sql.NullString
	err := row.Scan(&u.ID, &u.Username, &displayName, &u.PINHash, &u.Role, &u.TokenVersion, &avatarIcon, &avatarPath, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	u.DisplayName = displayName.String
	u.AvatarIcon = avatarIcon.String
	u.AvatarPath = avatarPath.String
	return &u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT user_id, username, display_name, pin_hash, role, token_version, avatar_icon, avatar_path, created_at FROM users WHERE username = ?`
	row := r.db.Reader.QueryRowContext(ctx, query, username)

	var u model.User
	var displayName, avatarIcon, avatarPath sql.NullString
	err := row.Scan(&u.ID, &u.Username, &displayName, &u.PINHash, &u.Role, &u.TokenVersion, &avatarIcon, &avatarPath, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	u.DisplayName = displayName.String
	u.AvatarIcon = avatarIcon.String
	u.AvatarPath = avatarPath.String
	return &u, nil
}

func (r *UserRepository) ListAll(ctx context.Context) ([]*model.User, error) {
	query := `SELECT user_id, username, display_name, pin_hash, role, token_version, avatar_icon, avatar_path, created_at FROM users ORDER BY role DESC, username ASC`
	rows, err := r.db.Reader.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		var displayName, avatarIcon, avatarPath sql.NullString
		if err := rows.Scan(&u.ID, &u.Username, &displayName, &u.PINHash, &u.Role, &u.TokenVersion, &avatarIcon, &avatarPath, &u.CreatedAt); err != nil {
			return nil, err
		}
		u.DisplayName = displayName.String
		u.AvatarIcon = avatarIcon.String
		u.AvatarPath = avatarPath.String
		users = append(users, &u)
	}
	return users, rows.Err()
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) (int64, error) {
	query := `INSERT INTO users (username, display_name, pin_hash, role, token_version, avatar_icon, avatar_path) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.Writer.ExecContext(ctx, query, u.Username, u.DisplayName, u.PINHash, u.Role, 1, u.AvatarIcon, u.AvatarPath)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	query := `UPDATE users SET display_name = ?, avatar_icon = ?, avatar_path = ? WHERE user_id = ?`
	_, err := r.db.Writer.ExecContext(ctx, query, u.DisplayName, u.AvatarIcon, u.AvatarPath, u.ID)
	return err
}

func (r *UserRepository) UpdatePIN(ctx context.Context, userID int64, pinHash string) error {
	query := `UPDATE users SET pin_hash = ?, token_version = token_version + 1 WHERE user_id = ?`
	_, err := r.db.Writer.ExecContext(ctx, query, pinHash, userID)
	return err
}

func (r *UserRepository) IncrementTokenVersion(ctx context.Context, userID int64) error {
	query := `UPDATE users SET token_version = token_version + 1 WHERE user_id = ?`
	_, err := r.db.Writer.ExecContext(ctx, query, userID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE user_id = ?`
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
