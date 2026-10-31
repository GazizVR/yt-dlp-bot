package service

import (
	"context"
	"errors"

	"github.com/gazizvr/yt-dlp-bot/internal/domain"
)

type Downloader interface {
	DownloadVideo(url, outDir string) (string, error)
	DownloadAudio(url, outDir string) (string, error)
}

type RequestStore interface {
	Get(ctx context.Context, id int64) (string, error)
	Create(ctx context.Context, url string) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type MediaService struct {
	downloader Downloader
	store      RequestStore
	fileDir    string
}

func NewMediaService(downloader Downloader, store RequestStore, filesDir string) *MediaService {
	return &MediaService{
		downloader: downloader,
		store:      store,
		fileDir:    filesDir,
	}
}

func (s *MediaService) Create(ctx context.Context, url string) (int64, error) {
	id, err := s.store.Create(ctx, url)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *MediaService) Download(
	ctx context.Context,
	id int64,
	action domain.Action,
) (filePath string, isVideo bool, err error) {
	isVideo = action == domain.DownloadVideo
	mediaURL, err := s.store.Get(ctx, id)
	if err != nil {
		return "", isVideo, err
	}
	switch action {
	case domain.DownloadVideo:
		filePath, err = s.downloader.DownloadVideo(mediaURL, s.fileDir)
		if err != nil {
			return "", isVideo, err
		}
	case domain.DownloadAudio:
		filePath, err = s.downloader.DownloadAudio(mediaURL, s.fileDir)
		if err != nil {
			return "", isVideo, err
		}
	default:
		return "", isVideo, errors.New("unexpected action")
	}
	return filePath, isVideo, nil
}

func (s *MediaService) Release(ctx context.Context, id int64) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
