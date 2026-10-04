package tests

import (
	"sort"
	"strings"
	"testing"

	"github.com/WesleyEsq/Flan-Media-Server/internal/service"
)

func TestCleanTitleAndYear(t *testing.T) {
	tests := []struct {
		input        string
		expectedName string
		expectedYear int
	}{
		{
			input:        "Inception.2010.1080p.BluRay.x264.mp4",
			expectedName: "Inception",
			expectedYear: 2010,
		},
		{
			input:        "Blade.Runner.1982.Final.Cut.2160p.mkv",
			expectedName: "Blade Runner Final Cut",
			expectedYear: 1982,
		},
		{
			input:        "The.Matrix.1999.Remux.1080p.mkv",
			expectedName: "The Matrix",
			expectedYear: 1999,
		},
		{
			input:        "Dune.Part.Two.2024.UHD.mkv",
			expectedName: "Dune Part Two",
			expectedYear: 2024,
		},
		{
			input:        "Clean Book Title.pdf",
			expectedName: "Clean Book Title",
			expectedYear: 0,
		},
	}

	for _, tc := range tests {
		title, year := service.CleanTitleAndYear(tc.input)
		if title != tc.expectedName {
			t.Errorf("CleanTitleAndYear(%q) title = %q; want %q", tc.input, title, tc.expectedName)
		}
		if year != tc.expectedYear {
			t.Errorf("CleanTitleAndYear(%q) year = %d; want %d", tc.input, year, tc.expectedYear)
		}
	}
}

func TestParseEpisodeInfo(t *testing.T) {
	tests := []struct {
		filename    string
		wantSeason  int
		wantEpisode int
		wantHas     bool
	}{
		{
			filename:    "Breaking.Bad.S01E03.720p.mkv",
			wantSeason:  1,
			wantEpisode: 3,
			wantHas:     true,
		},
		{
			filename:    "The.Wire.S04E12.mkv",
			wantSeason:  4,
			wantEpisode: 12,
			wantHas:     true,
		},
		{
			filename:    "Show.Name.2x08.mp4",
			wantSeason:  2,
			wantEpisode: 8,
			wantHas:     true,
		},
		{
			filename:    "Anime.Series.ep05.mkv",
			wantSeason:  1,
			wantEpisode: 5,
			wantHas:     true,
		},
		{
			filename:    "Standalone.Movie.2020.mkv",
			wantSeason:  1,
			wantEpisode: 0,
			wantHas:     false,
		},
	}

	for _, tc := range tests {
		s, e, has := service.ParseEpisodeInfo(tc.filename)
		if s != tc.wantSeason || e != tc.wantEpisode || has != tc.wantHas {
			t.Errorf("ParseEpisodeInfo(%q) = (%d, %d, %v); want (%d, %d, %v)",
				tc.filename, s, e, has, tc.wantSeason, tc.wantEpisode, tc.wantHas)
		}
	}
}

func TestNaturalSort(t *testing.T) {
	items := []string{
		"Episode 10",
		"Episode 1",
		"Episode 20",
		"Episode 2",
		"Episode 100",
	}

	expected := []string{
		"Episode 1",
		"Episode 2",
		"Episode 10",
		"Episode 20",
		"Episode 100",
	}

	sort.Slice(items, func(i, j int) bool {
		return service.NaturalLess(items[i], items[j])
	})

	for i := range items {
		if items[i] != expected[i] {
			t.Fatalf("NaturalSort result at [%d] = %q; want %q", i, items[i], expected[i])
		}
	}
}

func TestConvertSRTToWebVTT(t *testing.T) {
	srtInput := `1
00:01:20,000 --> 00:01:23,500
Hello world!

2
00:01:24,100 --> 00:01:26,800
This is a test subtitle.
`

	vttOutput := string(service.ConvertSRTToWebVTT([]byte(srtInput)))

	if !strings.HasPrefix(vttOutput, "WEBVTT") {
		t.Errorf("Expected WEBVTT header, got: %s", vttOutput)
	}

	if !strings.Contains(vttOutput, "00:01:20.000 --> 00:01:23.500") {
		t.Errorf("Expected dot decimal timestamp, got: %s", vttOutput)
	}

	if !strings.Contains(vttOutput, "Hello world!") {
		t.Errorf("Expected subtitle text preserved, got: %s", vttOutput)
	}
}
