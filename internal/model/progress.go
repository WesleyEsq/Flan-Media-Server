package model

import "time"

type Progress struct {
	UserID       int64     `json:"user_id"`
	MediaType    MediaType `json:"media_type"`
	FileID       int64     `json:"file_id"`
	PositionData string    `json:"position_data"`
	Percentage   float64   `json:"percentage"`
	IsFinished   bool      `json:"is_finished"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ContinueWatchingItem struct {
	MediaID         int64     `json:"media_id"`
	MediaType       MediaType `json:"media_type"`
	FileID          int64     `json:"file_id"`
	Title           string    `json:"title"`
	Subtitle        string    `json:"subtitle"`
	CoverPath       string    `json:"cover_path,omitempty"`
	PositionSeconds float64   `json:"position_seconds"`
	Percentage      float64   `json:"percentage"`
	FormattedPos    string    `json:"formatted_pos,omitempty"`
	StatusTag       string    `json:"status_tag"` // "In Progress" or "Reading"
	UpdatedAt       time.Time `json:"updated_at"`
}
