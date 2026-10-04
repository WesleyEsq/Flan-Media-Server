package service

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	seasonEpRegex1 = regexp.MustCompile(`(?i)[sS](\d+)[eE](\d+)`)
	seasonEpRegex2 = regexp.MustCompile(`(?i)(\d+)x(\d+)`)
	epOnlyRegex    = regexp.MustCompile(`(?i)(?:ep|episode|e)\s*(\d+)`)
	yearRegex      = regexp.MustCompile(`\b(19\d\d|20\d\d)\b`)

	// Noise tokens to strip from scene/torrent names
	noiseRegex = regexp.MustCompile(`(?i)\b(1080p|720p|480p|2160p|4k|uhd|bluray|bdrip|brrip|web-?dl|webrip|hdtv|x264|x265|hevc|h\.?264|h\.?265|aac|dts|ac3|e-?ac3|ddp?5\.1|truehd|remux|repack|proper|unrated|extended|directors\.cut|multi|sub|dub|dual|yify|eztv|rarbg|psa|galaxytv|vxt)\b`)
	bracketRegex = regexp.MustCompile(`\[.*?\]|\(.*?\)|-.*$`)
)

// CleanTitleAndYear extracts human readable title and release year from directory or file name
func CleanTitleAndYear(raw string) (string, int) {
	// Strip file extension if any
	name := raw
	if idx := strings.LastIndex(name, "."); idx != -1 && len(name)-idx <= 5 {
		ext := strings.ToLower(name[idx:])
		if ext == ".mp4" || ext == ".mkv" || ext == ".webm" || ext == ".epub" || ext == ".pdf" {
			name = name[:idx]
		}
	}

	// Extract year before removing brackets if year is in parentheses: e.g. "Blade Runner (1982)"
	year := 0
	if matches := yearRegex.FindAllString(name, -1); len(matches) > 0 {
		if y, err := strconv.Atoi(matches[len(matches)-1]); err == nil && y >= 1900 && y <= 2100 {
			year = y
		}
	}

	// Clean dots and underscores
	cleaned := strings.ReplaceAll(name, ".", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")

	// Strip release noise
	cleaned = noiseRegex.ReplaceAllString(cleaned, " ")

	// Strip year from title if present
	if year > 0 {
		cleaned = strings.ReplaceAll(cleaned, strconv.Itoa(year), " ")
	}

	// Strip trailing brackets/dashes
	cleaned = bracketRegex.ReplaceAllString(cleaned, " ")

	// Normalize spaces
	fields := strings.Fields(cleaned)
	cleanedTitle := strings.Join(fields, " ")
	cleanedTitle = strings.Trim(cleanedTitle, " -–_")

	if cleanedTitle == "" {
		cleanedTitle = raw
	}

	return cleanedTitle, year
}

// ParseEpisodeInfo detects SxxExx or episode numbers from filenames
func ParseEpisodeInfo(filename string) (season int, episode int, hasEpisode bool) {
	if m := seasonEpRegex1.FindStringSubmatch(filename); len(m) == 3 {
		s, _ := strconv.Atoi(m[1])
		e, _ := strconv.Atoi(m[2])
		return s, e, true
	}
	if m := seasonEpRegex2.FindStringSubmatch(filename); len(m) == 3 {
		s, _ := strconv.Atoi(m[1])
		e, _ := strconv.Atoi(m[2])
		return s, e, true
	}
	if m := epOnlyRegex.FindStringSubmatch(filename); len(m) == 2 {
		e, _ := strconv.Atoi(m[1])
		return 1, e, true
	}
	return 1, 0, false
}

// NaturalLess compares two strings using natural alphanumeric ordering
func NaturalLess(a, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)
	i, j := 0, 0
	for i < len(aRunes) && j < len(bRunes) {
		rA := aRunes[i]
		rB := bRunes[j]

		if unicode.IsDigit(rA) && unicode.IsDigit(rB) {
			// Extract full numbers
			startA := i
			for i < len(aRunes) && unicode.IsDigit(aRunes[i]) {
				i++
			}
			numA, _ := strconv.ParseUint(string(aRunes[startA:i]), 10, 64)

			startB := j
			for j < len(bRunes) && unicode.IsDigit(bRunes[j]) {
				j++
			}
			numB, _ := strconv.ParseUint(string(bRunes[startB:j]), 10, 64)

			if numA != numB {
				return numA < numB
			}
			continue
		}

		if unicode.ToLower(rA) != unicode.ToLower(rB) {
			return unicode.ToLower(rA) < unicode.ToLower(rB)
		}
		i++
		j++
	}
	return len(aRunes) < len(bRunes)
}

// ExtractMP4Duration parses the mvhd atom from MP4/MOV container in pure Go
func ExtractMP4Duration(r io.ReadSeeker) (int, error) {
	var buf [8]byte
	for {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return 0, err
		}
		atomSize := binary.BigEndian.Uint32(buf[0:4])
		atomType := string(buf[4:8])

		if atomSize < 8 {
			break
		}

		if atomType == "moov" {
			// Recurse into moov atom
			return parseMoovForDuration(io.LimitReader(r, int64(atomSize-8)))
		}

		// Skip atom
		if _, err := r.Seek(int64(atomSize-8), io.SeekCurrent); err != nil {
			return 0, err
		}
	}
	return 0, io.EOF
}

func parseMoovForDuration(r io.Reader) (int, error) {
	var buf [8]byte
	for {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return 0, err
		}
		atomSize := binary.BigEndian.Uint32(buf[0:4])
		atomType := string(buf[4:8])

		if atomSize < 8 {
			break
		}

		payloadSize := int64(atomSize - 8)
		if atomType == "mvhd" {
			mvhdData := make([]byte, payloadSize)
			if _, err := io.ReadFull(r, mvhdData); err != nil {
				return 0, err
			}
			version := mvhdData[0]
			var timescale uint32
			var duration uint64
			if version == 0 && len(mvhdData) >= 20 {
				timescale = binary.BigEndian.Uint32(mvhdData[12:16])
				duration = uint64(binary.BigEndian.Uint32(mvhdData[16:20]))
			} else if version == 1 && len(mvhdData) >= 28 {
				timescale = binary.BigEndian.Uint32(mvhdData[20:24])
				duration = binary.BigEndian.Uint64(mvhdData[24:32])
			}
			if timescale > 0 {
				return int(duration / uint64(timescale)), nil
			}
			return 0, nil
		}

		// Discard other moov children
		if _, err := io.CopyN(io.Discard, r, payloadSize); err != nil {
			return 0, err
		}
	}
	return 0, io.EOF
}

// ExtractMKVDuration parses Matroska / WebM EBML Segment Info Duration in pure Go
func ExtractMKVDuration(r io.ReadSeeker) (int, error) {
	// Matroska Segment Info contains Info ID 0x1549A966, TimecodeScale 0x2AD7B1, Duration 0x4489
	// Read first 64KB where Segment Info typically lives
	buf := make([]byte, 65536)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	}
	data := buf[:n]

	// Look for Segment Info element ID: 0x15 0x49 0xA9 0x66
	infoMarker := []byte{0x15, 0x49, 0xA9, 0x66}
	infoIdx := bytes.Index(data, infoMarker)
	if infoIdx == -1 {
		return 0, io.EOF
	}

	searchScope := data[infoIdx:]
	if len(searchScope) > 4096 {
		searchScope = searchScope[:4096]
	}

	// Default TimecodeScale is 1,000,000 nanoseconds = 1ms
	timecodeScale := float64(1000000)
	tcMarker := []byte{0x2A, 0xD7, 0xB1}
	if tcIdx := bytes.Index(searchScope, tcMarker); tcIdx != -1 {
		idx := tcIdx + 3
		if idx < len(searchScope) {
			sizeLen, tcSize := parseEBMLVarInt(searchScope[idx:])
			idx += sizeLen
			if idx+tcSize <= len(searchScope) {
				tcVal := uint64(0)
				for k := 0; k < tcSize; k++ {
					tcVal = (tcVal << 8) | uint64(searchScope[idx+k])
				}
				if tcVal > 0 {
					timecodeScale = float64(tcVal)
				}
			}
		}
	}

	// Look for Duration ID: 0x44 0x89
	durMarker := []byte{0x44, 0x89}
	durIdx := bytes.Index(searchScope, durMarker)
	if durIdx == -1 {
		return 0, io.EOF
	}

	idx := durIdx + 2
	if idx >= len(searchScope) {
		return 0, io.EOF
	}
	sizeLen, durSize := parseEBMLVarInt(searchScope[idx:])
	idx += sizeLen
	if idx+durSize > len(searchScope) {
		return 0, io.EOF
	}

	rawDurBytes := searchScope[idx : idx+durSize]
	var durationInTimecodes float64
	if durSize == 4 {
		bits := binary.BigEndian.Uint32(rawDurBytes)
		durationInTimecodes = float64(math.Float32frombits(bits))
	} else if durSize == 8 {
		bits := binary.BigEndian.Uint64(rawDurBytes)
		durationInTimecodes = math.Float64frombits(bits)
	}

	durationSec := (durationInTimecodes * timecodeScale) / 1e9
	return int(durationSec), nil
}

func parseEBMLVarInt(b []byte) (length int, val int) {
	if len(b) == 0 {
		return 0, 0
	}
	first := b[0]
	var mask byte = 0x80
	for length = 1; length <= 8; length++ {
		if (first & mask) != 0 {
			break
		}
		mask >>= 1
	}
	if length > len(b) {
		return 0, 0
	}
	res := int(first &^ mask)
	for i := 1; i < length; i++ {
		res = (res << 8) | int(b[i])
	}
	return length, res
}

// ConvertSRTToWebVTT translates SRT subtitle bytes into WebVTT in pure Go
func ConvertSRTToWebVTT(srtData []byte) []byte {
	var out bytes.Buffer
	out.WriteString("WEBVTT\n\n")

	content := string(srtData)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// Replace comma with dot in timestamps (e.g., 00:01:23,456 --> 00:01:23.456)
		if strings.Contains(line, "-->") {
			line = strings.ReplaceAll(line, ",", ".")
		}
		out.WriteString(line)
		out.WriteString("\n")
	}

	return out.Bytes()
}
