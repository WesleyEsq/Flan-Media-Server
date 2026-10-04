package model

import "time"

type VideoType string

const (
	VideoTypeMovie  VideoType = "movie"
	VideoTypeSeries VideoType = "series"
)

type Video struct {
	ID             int64        `json:"video_id"`
	SourceID       int64        `json:"source_id"`
	Title          string       `json:"title"`
	VideoType      VideoType    `json:"video_type"`
	ReleaseYear    int          `json:"release_year,omitempty"`
	Overview       string       `json:"overview,omitempty"`
	CoverPath      string       `json:"cover_path,omitempty"`
	FolderPath     string       `json:"folder_path"`
	IsHidden       bool         `json:"is_hidden"`
	MetadataLocked bool         `json:"metadata_locked"`
	CreatedAt      time.Time    `json:"created_at"`
	Files          []*VideoFile `json:"files,omitempty"`

	// View/Display helper fields
	SourceName  string `json:"source_name,omitempty"`
	ActiveFiles int    `json:"active_files,omitempty"`
}

type VideoFile struct {
	ID              int64  `json:"file_id"`
	VideoID         int64  `json:"video_id"`
	Title           string `json:"title"`
	CustomTitle     string `json:"custom_title,omitempty"`
	RelativePath    string `json:"relative_path"`
	FileSize        int64  `json:"file_size"`
	MTime           int64  `json:"mtime"`
	Format          string `json:"format"`
	DurationSeconds int    `json:"duration_seconds"`
	SeasonNumber    int    `json:"season_number"`
	EpisodeNumber   int    `json:"episode_number"`
	OrderIndex      int    `json:"order_index"`
	IsHidden        bool   `json:"is_hidden"`
	IsMissing       bool   `json:"is_missing"`

	// Progress state for current user
	PositionSeconds float64 `json:"position_seconds,omitempty"`
	Percentage      float64 `json:"percentage,omitempty"`
	IsFinished      bool    `json:"is_finished,omitempty"`
}

func (f *VideoFile) DisplayTitle() string {
	if f.CustomTitle != "" {
		return f.CustomTitle
	}
	return f.Title
}
