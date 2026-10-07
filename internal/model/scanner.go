package model

type DiscoveredFile struct {
	RelativePath    string `json:"relative_path"`
	Title           string `json:"title"`
	FileSize        int64  `json:"file_size"`
	MTime           int64  `json:"mtime"`
	Format          string `json:"format"`
	DurationSeconds int    `json:"duration_seconds"`
	SeasonNumber    int    `json:"season_number"`
	EpisodeNumber   int    `json:"episode_number"`
	OrderIndex      int    `json:"order_index"`
}

type CandidateContainer struct {
	FolderPath  string            `json:"folder_path"`
	Title       string            `json:"title"`
	ReleaseYear int               `json:"release_year,omitempty"`
	VideoType   VideoType         `json:"video_type,omitempty"`
	Author      string            `json:"author,omitempty"`
	CoverPath   string            `json:"cover_path,omitempty"`
	Files       []*DiscoveredFile `json:"files"`
}

type IngestionPayload struct {
	Approved []*CandidateContainer `json:"approved"`
}

type ScanStatus struct {
	IsScanning   bool   `json:"is_scanning"`
	CurrentPath  string `json:"current_path,omitempty"`
	ItemsScanned int    `json:"items_scanned"`
	Errors       int    `json:"errors"`
}
