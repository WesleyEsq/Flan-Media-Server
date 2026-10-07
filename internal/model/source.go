package model

import "time"

type MediaType string

const (
	MediaTypeVideo MediaType = "video"
	MediaTypeBook  MediaType = "book"
)

type StorageSource struct {
	ID              int64     `json:"source_id"`
	Name            string    `json:"name"`
	MediaType       MediaType `json:"media_type"`
	FolderPath      string    `json:"folder_path"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	FreeSpaceBytes  uint64    `json:"free_space_bytes,omitempty"`
	TotalSpaceBytes uint64    `json:"total_space_bytes,omitempty"`
	IsOnline        bool      `json:"is_online"`
}
