package model

import "time"

type BookFormat string

const (
	BookFormatEPUB BookFormat = "epub"
	BookFormatPDF  BookFormat = "pdf"
)

type Book struct {
	ID             int64       `json:"book_id"`
	SourceID       int64       `json:"source_id"`
	Title          string      `json:"title"`
	Author         string      `json:"author,omitempty"`
	Overview       string      `json:"overview,omitempty"`
	CoverPath      string      `json:"cover_path,omitempty"`
	FolderPath     string      `json:"folder_path"`
	IsHidden       bool        `json:"is_hidden"`
	MetadataLocked bool        `json:"metadata_locked"`
	CreatedAt      time.Time   `json:"created_at"`
	Files          []*BookFile `json:"files,omitempty"`

	// View/Display helper fields
	SourceName  string `json:"source_name,omitempty"`
	ActiveFiles int    `json:"active_files,omitempty"`
}

type BookFile struct {
	ID           int64      `json:"file_id"`
	BookID       int64      `json:"book_id"`
	Title        string     `json:"title"`
	CustomTitle  string     `json:"custom_title,omitempty"`
	RelativePath string     `json:"relative_path"`
	FileSize     int64      `json:"file_size"`
	MTime        int64      `json:"mtime"`
	Format       BookFormat `json:"format"`
	OrderIndex   int        `json:"order_index"`
	IsHidden     bool       `json:"is_hidden"`
	IsMissing    bool       `json:"is_missing"`

	// Progress state for current user: "unread", "reading", "finished"
	ReadingStatus string `json:"reading_status,omitempty"`
	IsFinished    bool   `json:"is_finished,omitempty"`
}

func (f *BookFile) DisplayTitle() string {
	if f.CustomTitle != "" {
		return f.CustomTitle
	}
	return f.Title
}
