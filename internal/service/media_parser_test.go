package service

import (
	"testing"
)

func TestParseEpisodeInfo(t *testing.T) {
	tests := []struct {
		filename    string
		wantSeason  int
		wantEpisode int
		wantHasEp   bool
	}{
		{
			filename:    "Sakura Card Captor Episodio 01 en Español latino.mp4",
			wantSeason:  1,
			wantEpisode: 1,
			wantHasEp:   true,
		},
		{
			filename:    "Sakura Card Captor Episodio 02.mp4",
			wantSeason:  1,
			wantEpisode: 2,
			wantHasEp:   true,
		},
		{
			filename:    "Capítulo 15 - La carta sellada.mkv",
			wantSeason:  1,
			wantEpisode: 15,
			wantHasEp:   true,
		},
		{
			filename:    "Neon Genesis Evangelion S01E26.mkv",
			wantSeason:  1,
			wantEpisode: 26,
			wantHasEp:   true,
		},
		{
			filename:    "Attack on Titan 2x12 [1080p].mkv",
			wantSeason:  2,
			wantEpisode: 12,
			wantHasEp:   true,
		},
		{
			filename:    "Cowboy Bebop - 05 - Ballad of Fallen Angels.mp4",
			wantSeason:  1,
			wantEpisode: 5,
			wantHasEp:   true,
		},
		{
			filename:    "The Matrix (1999).mp4",
			wantSeason:  1,
			wantEpisode: 0,
			wantHasEp:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			s, e, has := ParseEpisodeInfo(tt.filename)
			if s != tt.wantSeason || e != tt.wantEpisode || has != tt.wantHasEp {
				t.Errorf("ParseEpisodeInfo(%q) = (%d, %d, %v), want (%d, %d, %v)",
					tt.filename, s, e, has, tt.wantSeason, tt.wantEpisode, tt.wantHasEp)
			}
		})
	}
}
