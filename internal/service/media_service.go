package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	"github.com/WesleyEsq/Flan-Media-Server/internal/repository"
)

type MediaService struct {
	videoRepo    *repository.VideoRepository
	bookRepo     *repository.BookRepository
	progressRepo *repository.ProgressRepository
	sourceRepo   *repository.SourceRepository
	dataDir      string
}

func NewMediaService(
	videoRepo *repository.VideoRepository,
	bookRepo *repository.BookRepository,
	progressRepo *repository.ProgressRepository,
	sourceRepo *repository.SourceRepository,
	dataDir string,
) *MediaService {
	return &MediaService{
		videoRepo:    videoRepo,
		bookRepo:     bookRepo,
		progressRepo: progressRepo,
		sourceRepo:   sourceRepo,
		dataDir:      dataDir,
	}
}

func (s *MediaService) ListVideos(ctx context.Context, filterType string, showHidden bool) ([]*model.Video, error) {
	return s.videoRepo.List(ctx, filterType, showHidden)
}

func (s *MediaService) GetVideo(ctx context.Context, videoID int64, userID int64, showHidden bool) (*model.Video, []*model.VideoFile, error) {
	v, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil {
		return nil, nil, err
	}
	files, err := s.videoRepo.ListFilesByVideoID(ctx, videoID, userID, showHidden)
	if err != nil {
		return nil, nil, err
	}
	v.Files = files
	return v, files, nil
}

func (s *MediaService) ListBooks(ctx context.Context, formatFilter string, showHidden bool) ([]*model.Book, error) {
	return s.bookRepo.List(ctx, formatFilter, showHidden)
}

func (s *MediaService) GetBook(ctx context.Context, bookID int64, userID int64, showHidden bool) (*model.Book, []*model.BookFile, error) {
	b, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, nil, err
	}
	files, err := s.bookRepo.ListFilesByBookID(ctx, bookID, userID, showHidden)
	if err != nil {
		return nil, nil, err
	}
	b.Files = files
	return b, files, nil
}

func (s *MediaService) Search(ctx context.Context, query string, showHidden bool) ([]*model.Video, []*model.Book, error) {
	videos, err := s.videoRepo.Search(ctx, query, showHidden)
	if err != nil {
		return nil, nil, err
	}
	books, err := s.bookRepo.Search(ctx, query, showHidden)
	if err != nil {
		return nil, nil, err
	}
	return videos, books, nil
}

func (s *MediaService) GetContinueWatching(ctx context.Context, userID int64) ([]*model.ContinueWatchingItem, error) {
	return s.progressRepo.ListContinueWatching(ctx, userID, 10)
}

func (s *MediaService) GetJumpBackInBooks(ctx context.Context, userID int64) ([]*model.ContinueWatchingItem, error) {
	return s.progressRepo.ListJumpBackInBooks(ctx, userID, 10)
}

func (s *MediaService) SaveProgress(ctx context.Context, p *model.Progress) error {
	return s.progressRepo.Upsert(ctx, p)
}

func (s *MediaService) GetProgress(ctx context.Context, userID int64, mediaType model.MediaType, fileID int64) (*model.Progress, error) {
	return s.progressRepo.Get(ctx, userID, mediaType, fileID)
}

func (s *MediaService) UpdateVideoMetadata(ctx context.Context, v *model.Video, files []*model.VideoFile) error {
	if err := s.videoRepo.UpdateMetadata(ctx, v); err != nil {
		return err
	}
	if len(files) > 0 {
		return s.videoRepo.BatchUpdateFiles(ctx, files)
	}
	return nil
}

func (s *MediaService) UpdateBookMetadata(ctx context.Context, b *model.Book) error {
	return s.bookRepo.UpdateMetadata(ctx, b)
}

func (s *MediaService) SaveCoverArt(ctx context.Context, mediaType model.MediaType, id int64, coverBytes []byte, ext string) (string, error) {
	coversDir := filepath.Join(s.dataDir, "covers", string(mediaType))
	if err := os.MkdirAll(coversDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d%s", id, ext)
	targetPath := filepath.Join(coversDir, filename)
	if err := os.WriteFile(targetPath, coverBytes, 0644); err != nil {
		return "", err
	}

	relPath := fmt.Sprintf("data/covers/%s/%s", mediaType, filename)
	if mediaType == model.MediaTypeVideo {
		v, err := s.videoRepo.GetByID(ctx, id)
		if err == nil {
			v.CoverPath = relPath
			_ = s.videoRepo.UpdateMetadata(ctx, v)
		}
	} else {
		b, err := s.bookRepo.GetByID(ctx, id)
		if err == nil {
			b.CoverPath = relPath
			_ = s.bookRepo.UpdateMetadata(ctx, b)
		}
	}

	return relPath, nil
}

func (s *MediaService) RemoveCoverArt(ctx context.Context, mediaType model.MediaType, id int64) error {
	var oldCover string
	if mediaType == model.MediaTypeVideo {
		v, err := s.videoRepo.GetByID(ctx, id)
		if err != nil || v == nil {
			return fmt.Errorf("video not found: %w", err)
		}
		oldCover = v.CoverPath
		v.CoverPath = ""
		if err := s.videoRepo.UpdateMetadata(ctx, v); err != nil {
			return err
		}
	} else {
		b, err := s.bookRepo.GetByID(ctx, id)
		if err != nil || b == nil {
			return fmt.Errorf("book not found: %w", err)
		}
		oldCover = b.CoverPath
		b.CoverPath = ""
		if err := s.bookRepo.UpdateMetadata(ctx, b); err != nil {
			return err
		}
	}

	if oldCover != "" && strings.HasPrefix(oldCover, "data/covers/") {
		baseName := filepath.Base(oldCover)
		_ = os.Remove(filepath.Join(s.dataDir, "covers", string(mediaType), baseName))
	}

	return nil
}
