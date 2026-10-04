package service

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/WesleyEsq/Flan-Media-Server/internal/database"
	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
)

type sfCall struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type singleFlightGroup struct {
	mu sync.Mutex
	m  map[string]*sfCall
}

func (g *singleFlightGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error, bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*sfCall)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}
	c := new(sfCall)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err, false
}

type ScannerService struct {
	videoRepo  *repository.VideoRepository
	bookRepo   *repository.BookRepository
	sourceRepo *repository.SourceRepository
	sf         singleFlightGroup

	statusMu   sync.RWMutex
	lastStatus *model.ScanStatus
}

func NewScannerService(
	videoRepo *repository.VideoRepository,
	bookRepo *repository.BookRepository,
	sourceRepo *repository.SourceRepository,
) *ScannerService {
	return &ScannerService{
		videoRepo:  videoRepo,
		bookRepo:   bookRepo,
		sourceRepo: sourceRepo,
		lastStatus: &model.ScanStatus{},
	}
}

func (s *ScannerService) GetStatus() model.ScanStatus {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return *s.lastStatus
}

// ScanSource performs an ephemeral crawl of a storage source without writing to SQLite,
// returning gathered candidate items for the interactive review wizard.
func (s *ScannerService) ScanSource(ctx context.Context, source *model.StorageSource) ([]*model.CandidateContainer, error) {
	// 1. Verify .flan-keep marker file
	if err := database.VerifyMarker(source.FolderPath); err != nil {
		return nil, fmt.Errorf("%w in %s", model.ErrMissingMarker, source.FolderPath)
	}

	// 2. Prune any previously ingested containers that no longer exist on disk
	_ = s.pruneDeletedContainersForSource(ctx, source)

	if source.MediaType == model.MediaTypeVideo {
		return s.scanVideoSource(source)
	}
	return s.scanBookSource(source)
}

func (s *ScannerService) scanVideoSource(source *model.StorageSource) ([]*model.CandidateContainer, error) {
	entries, err := os.ReadDir(source.FolderPath)
	if err != nil {
		return nil, err
	}

	var candidates []*model.CandidateContainer

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(source.FolderPath, name)

		if !entry.IsDir() {
			// Flat video file in source root
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".mp4" || ext == ".mkv" || ext == ".webm" {
				info, err := entry.Info()
				if err != nil {
					continue
				}

				title, year := CleanTitleAndYear(name)
				dur := s.probeVideoDuration(fullPath, ext)

				c := &model.CandidateContainer{
					FolderPath:  name,
					Title:       title,
					ReleaseYear: year,
					VideoType:   model.VideoTypeMovie,
					Files: []*model.DiscoveredFile{
						{
							RelativePath:    "",
							Title:           title,
							FileSize:        info.Size(),
							MTime:           info.ModTime().Unix(),
							Format:          strings.TrimPrefix(ext, "."),
							DurationSeconds: dur,
							SeasonNumber:    1,
							EpisodeNumber:   0,
							OrderIndex:      0,
						},
					},
				}
				candidates = append(candidates, c)
			}
		} else {
			// Subdirectory container (Series or Movie folder)
			c, err := s.scanVideoDirectory(source.FolderPath, name)
			if err == nil && c != nil && len(c.Files) > 0 {
				candidates = append(candidates, c)
			}
		}
	}

	// Sort candidates alphabetically by title
	sort.Slice(candidates, func(i, j int) bool {
		return NaturalLess(candidates[i].Title, candidates[j].Title)
	})

	return candidates, nil
}

func (s *ScannerService) scanVideoDirectory(sourceRoot, dirName string) (*model.CandidateContainer, error) {
	containerPath := filepath.Join(sourceRoot, dirName)
	cleanTitle, year := CleanTitleAndYear(dirName)

	var files []*model.DiscoveredFile
	var coverPath string
	hasMultipleEpisodes := false

	err := filepath.WalkDir(containerPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		name := d.Name()
		if strings.HasPrefix(name, ".") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(name))

		// Check for local cover artwork
		if coverPath == "" && (ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp") {
			base := strings.ToLower(strings.TrimSuffix(name, ext))
			if base == "poster" || base == "cover" || base == "folder" || base == strings.ToLower(dirName) {
				relCover, _ := filepath.Rel(sourceRoot, path)
				coverPath = relCover
			}
		}

		if ext == ".mp4" || ext == ".mkv" || ext == ".webm" {
			info, err := d.Info()
			if err != nil {
				return nil
			}

			relPath, err := filepath.Rel(containerPath, path)
			if err != nil {
				return nil
			}

			season, episode, hasEp := ParseEpisodeInfo(name)
			if hasEp {
				hasMultipleEpisodes = true
			}

			var orderIndex int
			if season <= 1 {
				orderIndex = episode
			} else {
				orderIndex = (season * 1000) + episode
			}
			fileTitle, _ := CleanTitleAndYear(name)
			dur := s.probeVideoDuration(path, ext)

			files = append(files, &model.DiscoveredFile{
				RelativePath:    relPath,
				Title:           fileTitle,
				FileSize:        info.Size(),
				MTime:           info.ModTime().Unix(),
				Format:          strings.TrimPrefix(ext, "."),
				DurationSeconds: dur,
				SeasonNumber:    season,
				EpisodeNumber:   episode,
				OrderIndex:      orderIndex,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	// Classify container
	videoType := model.VideoTypeMovie
	if len(files) > 1 || hasMultipleEpisodes {
		videoType = model.VideoTypeSeries
	}

	// Sort files by orderIndex or naturally by title
	sort.Slice(files, func(i, j int) bool {
		if files[i].OrderIndex != files[j].OrderIndex {
			return files[i].OrderIndex < files[j].OrderIndex
		}
		return NaturalLess(files[i].Title, files[j].Title)
	})

	return &model.CandidateContainer{
		FolderPath:  dirName,
		Title:       cleanTitle,
		ReleaseYear: year,
		VideoType:   videoType,
		CoverPath:   coverPath,
		Files:       files,
	}, nil
}

func (s *ScannerService) probeVideoDuration(filePath, ext string) int {
	f, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer f.Close()

	if ext == ".mp4" {
		if d, err := ExtractMP4Duration(f); err == nil && d > 0 {
			return d
		}
	} else if ext == ".mkv" || ext == ".webm" {
		if d, err := ExtractMKVDuration(f); err == nil && d > 0 {
			return d
		}
	}
	return 0
}

func (s *ScannerService) scanBookSource(source *model.StorageSource) ([]*model.CandidateContainer, error) {
	entries, err := os.ReadDir(source.FolderPath)
	if err != nil {
		return nil, err
	}

	var candidates []*model.CandidateContainer

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".epub" || ext == ".pdf" {
				info, err := entry.Info()
				if err != nil {
					continue
				}

				title, _ := CleanTitleAndYear(name)
				c := &model.CandidateContainer{
					FolderPath: name,
					Title:      title,
					Files: []*model.DiscoveredFile{
						{
							RelativePath: "",
							Title:        title,
							FileSize:     info.Size(),
							MTime:        info.ModTime().Unix(),
							Format:       strings.TrimPrefix(ext, "."),
							OrderIndex:   0,
						},
					},
				}
				candidates = append(candidates, c)
			}
		} else {
			// Subdirectory container
			c, err := s.scanBookDirectory(source.FolderPath, name)
			if err == nil && c != nil && len(c.Files) > 0 {
				candidates = append(candidates, c)
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return NaturalLess(candidates[i].Title, candidates[j].Title)
	})

	return candidates, nil
}

func (s *ScannerService) scanBookDirectory(sourceRoot, dirName string) (*model.CandidateContainer, error) {
	containerPath := filepath.Join(sourceRoot, dirName)
	cleanTitle, _ := CleanTitleAndYear(dirName)

	var files []*model.DiscoveredFile
	var coverPath string

	err := filepath.WalkDir(containerPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		name := d.Name()
		if strings.HasPrefix(name, ".") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(name))

		if coverPath == "" && (ext == ".jpg" || ext == ".jpeg" || ext == ".png") {
			base := strings.ToLower(strings.TrimSuffix(name, ext))
			if base == "poster" || base == "cover" || base == "folder" {
				relCover, _ := filepath.Rel(sourceRoot, path)
				coverPath = relCover
			}
		}

		if ext == ".epub" || ext == ".pdf" {
			info, err := d.Info()
			if err != nil {
				return nil
			}

			relPath, err := filepath.Rel(containerPath, path)
			if err != nil {
				return nil
			}

			fileTitle, _ := CleanTitleAndYear(name)
			files = append(files, &model.DiscoveredFile{
				RelativePath: relPath,
				Title:        fileTitle,
				FileSize:     info.Size(),
				MTime:        info.ModTime().Unix(),
				Format:       strings.TrimPrefix(ext, "."),
				OrderIndex:   len(files) + 1,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	sort.Slice(files, func(i, j int) bool {
		return NaturalLess(files[i].Title, files[j].Title)
	})

	return &model.CandidateContainer{
		FolderPath: dirName,
		Title:      cleanTitle,
		CoverPath:  coverPath,
		Files:      files,
	}, nil
}

// IngestApproved commits user-approved items into the catalog with metadata_locked = 1
func (s *ScannerService) IngestApproved(ctx context.Context, sourceID int64, items []*model.CandidateContainer) error {
	source, err := s.sourceRepo.GetByID(ctx, sourceID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if source.MediaType == model.MediaTypeVideo {
			existing, err := s.videoRepo.GetBySourceAndFolder(ctx, sourceID, item.FolderPath)
			var videoID int64
			if err != nil || existing == nil {
				// Insert new video container
				v := &model.Video{
					SourceID:       sourceID,
					Title:          item.Title,
					VideoType:      item.VideoType,
					ReleaseYear:    item.ReleaseYear,
					CoverPath:      item.CoverPath,
					FolderPath:     item.FolderPath,
					MetadataLocked: true,
				}
				newID, err := s.videoRepo.Create(ctx, v)
				if err != nil {
					continue
				}
				videoID = newID
			} else {
				videoID = existing.ID
			}

			// Ingest files
			for _, f := range item.Files {
				vf := &model.VideoFile{
					VideoID:         videoID,
					Title:           f.Title,
					RelativePath:    f.RelativePath,
					FileSize:        f.FileSize,
					MTime:           f.MTime,
					Format:          f.Format,
					DurationSeconds: f.DurationSeconds,
					SeasonNumber:    f.SeasonNumber,
					EpisodeNumber:   f.EpisodeNumber,
					OrderIndex:      f.OrderIndex,
				}
				_, _ = s.videoRepo.CreateFile(ctx, vf)
			}
		} else {
			// Book container
			existing, err := s.bookRepo.GetBySourceAndFolder(ctx, sourceID, item.FolderPath)
			var bookID int64
			if err != nil || existing == nil {
				b := &model.Book{
					SourceID:       sourceID,
					Title:          item.Title,
					Author:         item.Author,
					CoverPath:      item.CoverPath,
					FolderPath:     item.FolderPath,
					MetadataLocked: true,
				}
				newID, err := s.bookRepo.Create(ctx, b)
				if err != nil {
					continue
				}
				bookID = newID
			} else {
				bookID = existing.ID
			}

			for _, f := range item.Files {
				bf := &model.BookFile{
					BookID:       bookID,
					Title:        f.Title,
					RelativePath: f.RelativePath,
					FileSize:     f.FileSize,
					MTime:        f.MTime,
					Format:       model.BookFormat(f.Format),
					OrderIndex:   f.OrderIndex,
				}
				_, _ = s.bookRepo.CreateFile(ctx, bf)
			}
		}
	}
	return nil
}

// MaintenanceScan verifies existing indexed files on disk, updates mtime, and flags missing files
func (s *ScannerService) MaintenanceScan(ctx context.Context) (*model.ScanStatus, error) {
	val, err, _ := s.sf.Do("maintenance_scan", func() (interface{}, error) {
		s.statusMu.Lock()
		s.lastStatus = &model.ScanStatus{IsScanning: true}
		s.statusMu.Unlock()

		defer func() {
			s.statusMu.Lock()
			s.lastStatus.IsScanning = false
			s.statusMu.Unlock()
		}()

		sources, err := s.sourceRepo.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		totalScanned := 0
		errorsCount := 0

		for _, src := range sources {
			if !src.IsActive {
				continue
			}

			s.statusMu.Lock()
			s.lastStatus.CurrentPath = src.FolderPath
			s.statusMu.Unlock()

			// Check .flan-keep
			if err := database.VerifyMarker(src.FolderPath); err != nil {
				// Mark all files missing for this source
				if src.MediaType == model.MediaTypeVideo {
					_ = s.videoRepo.MarkAllFilesMissingForSource(ctx, src.ID, true)
				} else {
					_ = s.bookRepo.MarkAllFilesMissingForSource(ctx, src.ID, true)
				}
				continue
			}

			if src.MediaType == model.MediaTypeVideo {
				videos, err := s.videoRepo.List(ctx, "all", true)
				if err != nil {
					errorsCount++
					continue
				}

				for _, v := range videos {
					if v.SourceID != src.ID {
						continue
					}

					var containerDir string
					if v.FolderPath == "" {
						containerDir = src.FolderPath
					} else {
						containerDir = filepath.Join(src.FolderPath, v.FolderPath)
					}

					// If the entire container folder on disk was deleted, remove from DB
					if _, err := os.Stat(containerDir); os.IsNotExist(err) {
						_ = s.videoRepo.Delete(ctx, v.ID)
						continue
					}

					files, err := s.videoRepo.ListFilesByVideoID(ctx, v.ID, 0, true)
					if err != nil {
						continue
					}

					activeCount := 0
					for _, f := range files {
						totalScanned++
						var diskPath string
						if f.RelativePath == "" {
							diskPath = containerDir
						} else {
							diskPath = filepath.Join(containerDir, f.RelativePath)
						}

						if _, err := os.Stat(diskPath); os.IsNotExist(err) {
							_ = s.videoRepo.DeleteFile(ctx, f.ID)
						} else {
							_ = s.videoRepo.SetFileMissing(ctx, f.ID, false)
							activeCount++
						}
					}

					// If all files were removed, delete the container
					if len(files) > 0 && activeCount == 0 {
						_ = s.videoRepo.Delete(ctx, v.ID)
					}
				}
			} else {
				books, err := s.bookRepo.List(ctx, "all", true)
				if err != nil {
					errorsCount++
					continue
				}

				for _, b := range books {
					if b.SourceID != src.ID {
						continue
					}

					var containerDir string
					if b.FolderPath == "" {
						containerDir = src.FolderPath
					} else {
						containerDir = filepath.Join(src.FolderPath, b.FolderPath)
					}

					// If the container folder on disk was deleted, remove from DB
					if _, err := os.Stat(containerDir); os.IsNotExist(err) {
						_ = s.bookRepo.Delete(ctx, b.ID)
						continue
					}

					files, err := s.bookRepo.ListFilesByBookID(ctx, b.ID, 0, true)
					if err != nil {
						continue
					}

					activeCount := 0
					for _, f := range files {
						totalScanned++
						var diskPath string
						if f.RelativePath == "" {
							diskPath = containerDir
						} else {
							diskPath = filepath.Join(containerDir, f.RelativePath)
						}

						if _, err := os.Stat(diskPath); os.IsNotExist(err) {
							_ = s.bookRepo.DeleteFile(ctx, f.ID)
						} else {
							_ = s.bookRepo.SetFileMissing(ctx, f.ID, false)
							activeCount++
						}
					}

					if len(files) > 0 && activeCount == 0 {
						_ = s.bookRepo.Delete(ctx, b.ID)
					}
				}
			}
		}

		s.statusMu.Lock()
		s.lastStatus.ItemsScanned = totalScanned
		s.lastStatus.Errors = errorsCount
		statusCopy := *s.lastStatus
		s.statusMu.Unlock()

		return &statusCopy, nil
	})

	if err != nil {
		return nil, err
	}
	return val.(*model.ScanStatus), nil
}

func (s *ScannerService) pruneDeletedContainersForSource(ctx context.Context, source *model.StorageSource) error {
	if source.MediaType == model.MediaTypeVideo {
		videos, err := s.videoRepo.List(ctx, "all", true)
		if err != nil {
			return err
		}
		for _, v := range videos {
			if v.SourceID != source.ID {
				continue
			}
			containerDir := filepath.Join(source.FolderPath, v.FolderPath)
			if v.FolderPath == "" {
				containerDir = source.FolderPath
			}
			if _, err := os.Stat(containerDir); os.IsNotExist(err) {
				_ = s.videoRepo.Delete(ctx, v.ID)
				continue
			}

			files, err := s.videoRepo.ListFilesByVideoID(ctx, v.ID, 0, true)
			if err != nil {
				continue
			}
			activeCount := 0
			for _, f := range files {
				var diskPath string
				if f.RelativePath == "" {
					diskPath = containerDir
				} else {
					diskPath = filepath.Join(containerDir, f.RelativePath)
				}
				if _, err := os.Stat(diskPath); os.IsNotExist(err) {
					_ = s.videoRepo.DeleteFile(ctx, f.ID)
				} else {
					activeCount++
				}
			}
			if len(files) > 0 && activeCount == 0 {
				_ = s.videoRepo.Delete(ctx, v.ID)
			}
		}
	} else {
		books, err := s.bookRepo.List(ctx, "all", true)
		if err != nil {
			return err
		}
		for _, b := range books {
			if b.SourceID != source.ID {
				continue
			}
			containerDir := filepath.Join(source.FolderPath, b.FolderPath)
			if b.FolderPath == "" {
				containerDir = source.FolderPath
			}
			if _, err := os.Stat(containerDir); os.IsNotExist(err) {
				_ = s.bookRepo.Delete(ctx, b.ID)
				continue
			}

			files, err := s.bookRepo.ListFilesByBookID(ctx, b.ID, 0, true)
			if err != nil {
				continue
			}
			activeCount := 0
			for _, f := range files {
				var diskPath string
				if f.RelativePath == "" {
					diskPath = containerDir
				} else {
					diskPath = filepath.Join(containerDir, f.RelativePath)
				}
				if _, err := os.Stat(diskPath); os.IsNotExist(err) {
					_ = s.bookRepo.DeleteFile(ctx, f.ID)
				} else {
					activeCount++
				}
			}
			if len(files) > 0 && activeCount == 0 {
				_ = s.bookRepo.Delete(ctx, b.ID)
			}
		}
	}
	return nil
}

