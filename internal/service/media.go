package service

import (
	"context"

	"github.com/gazizvr/yt-dlp-bot/internal/domain"
)

type Downloader interface {
	DownloadVideo(url, outDir string) (string, error)
}

type Converter interface {
	VideoToAudio(videoPath, outDir string) (string, error)
}

type UserStore interface {
	Get(ctx context.Context, id int64) (*domain.User, error)
	Set(ctx context.Context, id int64, user domain.User) error
	Delete(ctx context.Context, id int64) error
}

type MediaService struct {
	downloader Downloader
	converter  Converter
	store      UserStore
	filesDir   string
}

func NewMediaService(
	downloader Downloader,
	converter Converter,
	filesDir string,
) *MediaService {
	return &MediaService{
		downloader: downloader,
		converter:  converter,
		filesDir:   filesDir,
	}
}
