package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
)

type BookRepository struct {
	db *database.DB
}

func NewBookRepository(db *database.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) List(ctx context.Context, formatFilter string, showHidden bool) ([]*model.Book, error) {
	query := `
	SELECT b.book_id, b.source_id, b.title, b.author, b.overview, b.cover_path, b.folder_path, b.is_hidden, b.metadata_locked, b.created_at, s.name,
	       (SELECT COUNT(*) FROM book_files bf WHERE bf.book_id = b.book_id AND bf.is_missing = 0) AS active_files
	FROM books b
	JOIN storage_sources s ON b.source_id = s.source_id
	WHERE (b.is_hidden = 0 OR ?)
	  AND s.is_active = 1
	  AND (? = '' OR ? = 'all' OR EXISTS (SELECT 1 FROM book_files bf WHERE bf.book_id = b.book_id AND bf.format = ?))
	  AND EXISTS (SELECT 1 FROM book_files bf WHERE bf.book_id = b.book_id AND bf.is_missing = 0)
	ORDER BY b.title ASC;`

	rows, err := r.db.Reader.QueryContext(ctx, query, showHidden, formatFilter, formatFilter, formatFilter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*model.Book
	for rows.Next() {
		var b model.Book
		var author, overview, coverPath sql.NullString
		var isHiddenInt, metaLockedInt int

		err := rows.Scan(
			&b.ID, &b.SourceID, &b.Title, &author, &overview, &coverPath, &b.FolderPath,
			&isHiddenInt, &metaLockedInt, &b.CreatedAt, &b.SourceName, &b.ActiveFiles,
		)
		if err != nil {
			return nil, err
		}
		b.Author = author.String
		b.Overview = overview.String
		b.CoverPath = coverPath.String
		b.IsHidden = isHiddenInt == 1
		b.MetadataLocked = metaLockedInt == 1
		books = append(books, &b)
	}
	return books, rows.Err()
}

func (r *BookRepository) Search(ctx context.Context, queryStr string, showHidden bool) ([]*model.Book, error) {
	likePattern := "%" + queryStr + "%"
	query := `
	SELECT b.book_id, b.source_id, b.title, b.author, b.overview, b.cover_path, b.folder_path, b.is_hidden, b.metadata_locked, b.created_at, s.name,
	       (SELECT COUNT(*) FROM book_files bf WHERE bf.book_id = b.book_id AND bf.is_missing = 0) AS active_files
	FROM books b
	JOIN storage_sources s ON b.source_id = s.source_id
	WHERE (b.is_hidden = 0 OR ?)
	  AND s.is_active = 1
	  AND (b.title LIKE ? OR b.author LIKE ? OR b.overview LIKE ?)
	ORDER BY b.title ASC;`

	rows, err := r.db.Reader.QueryContext(ctx, query, showHidden, likePattern, likePattern, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*model.Book
	for rows.Next() {
		var b model.Book
		var author, overview, coverPath sql.NullString
		var isHiddenInt, metaLockedInt int

		err := rows.Scan(
			&b.ID, &b.SourceID, &b.Title, &author, &overview, &coverPath, &b.FolderPath,
			&isHiddenInt, &metaLockedInt, &b.CreatedAt, &b.SourceName, &b.ActiveFiles,
		)
		if err != nil {
			return nil, err
		}
		b.Author = author.String
		b.Overview = overview.String
		b.CoverPath = coverPath.String
		b.IsHidden = isHiddenInt == 1
		b.MetadataLocked = metaLockedInt == 1
		books = append(books, &b)
	}
	return books, rows.Err()
}

func (r *BookRepository) GetByID(ctx context.Context, id int64) (*model.Book, error) {
	query := `
	SELECT b.book_id, b.source_id, b.title, b.author, b.overview, b.cover_path, b.folder_path, b.is_hidden, b.metadata_locked, b.created_at, s.name
	FROM books b
	JOIN storage_sources s ON b.source_id = s.source_id
	WHERE b.book_id = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, id)
	var b model.Book
	var author, overview, coverPath sql.NullString
	var isHiddenInt, metaLockedInt int

	err := row.Scan(
		&b.ID, &b.SourceID, &b.Title, &author, &overview, &coverPath, &b.FolderPath,
		&isHiddenInt, &metaLockedInt, &b.CreatedAt, &b.SourceName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	b.Author = author.String
	b.Overview = overview.String
	b.CoverPath = coverPath.String
	b.IsHidden = isHiddenInt == 1
	b.MetadataLocked = metaLockedInt == 1
	return &b, nil
}

func (r *BookRepository) GetBySourceAndFolder(ctx context.Context, sourceID int64, folderPath string) (*model.Book, error) {
	query := `
	SELECT b.book_id, b.source_id, b.title, b.author, b.overview, b.cover_path, b.folder_path, b.is_hidden, b.metadata_locked, b.created_at
	FROM books b
	WHERE b.source_id = ? AND b.folder_path = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, sourceID, folderPath)
	var b model.Book
	var author, overview, coverPath sql.NullString
	var isHiddenInt, metaLockedInt int

	err := row.Scan(
		&b.ID, &b.SourceID, &b.Title, &author, &overview, &coverPath, &b.FolderPath,
		&isHiddenInt, &metaLockedInt, &b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	b.Author = author.String
	b.Overview = overview.String
	b.CoverPath = coverPath.String
	b.IsHidden = isHiddenInt == 1
	b.MetadataLocked = metaLockedInt == 1
	return &b, nil
}

func (r *BookRepository) Create(ctx context.Context, b *model.Book) (int64, error) {
	query := `
	INSERT INTO books (source_id, title, author, overview, cover_path, folder_path, is_hidden, metadata_locked)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	isHiddenInt := 0
	if b.IsHidden {
		isHiddenInt = 1
	}
	metaLockedInt := 0
	if b.MetadataLocked {
		metaLockedInt = 1
	}

	res, err := r.db.Writer.ExecContext(ctx, query, b.SourceID, b.Title, b.Author, b.Overview, b.CoverPath, b.FolderPath, isHiddenInt, metaLockedInt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *BookRepository) UpdateMetadata(ctx context.Context, b *model.Book) error {
	query := `
	UPDATE books 
	SET title = ?, author = ?, overview = ?, cover_path = ?, metadata_locked = 1
	WHERE book_id = ?;`

	_, err := r.db.Writer.ExecContext(ctx, query, b.Title, b.Author, b.Overview, b.CoverPath, b.ID)
	return err
}

func (r *BookRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM books WHERE book_id = ?;`
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

func (r *BookRepository) DeleteFile(ctx context.Context, fileID int64) error {
	query := `DELETE FROM book_files WHERE file_id = ?;`
	_, err := r.db.Writer.ExecContext(ctx, query, fileID)
	return err
}

func (r *BookRepository) ListFilesByBookID(ctx context.Context, bookID int64, userID int64, showHidden bool) ([]*model.BookFile, error) {
	query := `
	SELECT 
	    f.file_id, f.book_id, f.title, f.custom_title, f.relative_path, f.file_size, f.mtime, f.format,
	    f.order_index, f.is_hidden, f.is_missing,
	    COALESCE(p.position_data, 'unread') AS reading_status,
	    COALESCE(p.is_finished, 0) AS is_finished
	FROM book_files f
	LEFT JOIN progress p ON f.file_id = p.file_id AND p.media_type = 'book' AND p.user_id = ?
	WHERE f.book_id = ? AND (f.is_hidden = 0 OR ?)
	ORDER BY f.order_index ASC, f.title ASC;`

	rows, err := r.db.Reader.QueryContext(ctx, query, userID, bookID, showHidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*model.BookFile
	for rows.Next() {
		var f model.BookFile
		var customTitle sql.NullString
		var isHiddenInt, isMissingInt, isFinishedInt int

		err := rows.Scan(
			&f.ID, &f.BookID, &f.Title, &customTitle, &f.RelativePath, &f.FileSize, &f.MTime, &f.Format,
			&f.OrderIndex, &isHiddenInt, &isMissingInt, &f.ReadingStatus, &isFinishedInt,
		)
		if err != nil {
			return nil, err
		}
		f.CustomTitle = customTitle.String
		f.IsHidden = isHiddenInt == 1
		f.IsMissing = isMissingInt == 1
		f.IsFinished = isFinishedInt == 1
		files = append(files, &f)
	}
	return files, rows.Err()
}

func (r *BookRepository) GetFileByID(ctx context.Context, fileID int64) (*model.BookFile, error) {
	query := `
	SELECT file_id, book_id, title, custom_title, relative_path, file_size, mtime, format, order_index, is_hidden, is_missing
	FROM book_files
	WHERE file_id = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, fileID)
	var f model.BookFile
	var customTitle sql.NullString
	var isHiddenInt, isMissingInt int

	err := row.Scan(
		&f.ID, &f.BookID, &f.Title, &customTitle, &f.RelativePath, &f.FileSize, &f.MTime, &f.Format,
		&f.OrderIndex, &isHiddenInt, &isMissingInt,
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

func (r *BookRepository) GetFileWithContainer(ctx context.Context, fileID int64) (*model.BookFile, *model.Book, *model.StorageSource, error) {
	query := `
	SELECT 
	    f.file_id, f.book_id, f.title, f.custom_title, f.relative_path, f.file_size, f.mtime, f.format,
	    f.order_index, f.is_hidden, f.is_missing,
	    b.book_id, b.source_id, b.title, b.author, b.folder_path,
	    s.source_id, s.name, s.media_type, s.folder_path, s.is_active
	FROM book_files f
	JOIN books b ON f.book_id = b.book_id
	JOIN storage_sources s ON b.source_id = s.source_id
	WHERE f.file_id = ?;`

	row := r.db.Reader.QueryRowContext(ctx, query, fileID)
	var f model.BookFile
	var b model.Book
	var s model.StorageSource
	var customTitle, author sql.NullString
	var fHidden, fMissing, sActive int

	err := row.Scan(
		&f.ID, &f.BookID, &f.Title, &customTitle, &f.RelativePath, &f.FileSize, &f.MTime, &f.Format,
		&f.OrderIndex, &fHidden, &fMissing,
		&b.ID, &b.SourceID, &b.Title, &author, &b.FolderPath,
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
	b.Author = author.String
	s.IsActive = sActive == 1
	return &f, &b, &s, nil
}

func (r *BookRepository) CreateFile(ctx context.Context, f *model.BookFile) (int64, error) {
	query := `
	INSERT INTO book_files (book_id, title, custom_title, relative_path, file_size, mtime, format, order_index, is_hidden, is_missing)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	res, err := r.db.Writer.ExecContext(ctx, query,
		f.BookID, f.Title, f.CustomTitle, f.RelativePath, f.FileSize, f.MTime, f.Format,
		f.OrderIndex, 0, 0,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *BookRepository) UpdateFile(ctx context.Context, f *model.BookFile) error {
	query := `
	UPDATE book_files
	SET custom_title = ?, order_index = ?, is_hidden = ?
	WHERE file_id = ?;`

	isHiddenInt := 0
	if f.IsHidden {
		isHiddenInt = 1
	}
	_, err := r.db.Writer.ExecContext(ctx, query, f.CustomTitle, f.OrderIndex, isHiddenInt, f.ID)
	return err
}

func (r *BookRepository) SetFileMissing(ctx context.Context, fileID int64, missing bool) error {
	query := `UPDATE book_files SET is_missing = ? WHERE file_id = ?;`
	val := 0
	if missing {
		val = 1
	}
	_, err := r.db.Writer.ExecContext(ctx, query, val, fileID)
	return err
}

func (r *BookRepository) MarkAllFilesMissingForSource(ctx context.Context, sourceID int64, missing bool) error {
	query := `
	UPDATE book_files 
	SET is_missing = ? 
	WHERE book_id IN (SELECT book_id FROM books WHERE source_id = ?);`
	val := 0
	if missing {
		val = 1
	}
	_, err := r.db.Writer.ExecContext(ctx, query, val, sourceID)
	return err
}
