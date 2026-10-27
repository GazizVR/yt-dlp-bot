package service

import (
	"context"
	"log/slog"
	"strings"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

type Downloader interface {
    DownloadVideo(url, outputDir string) (string,error)
}

type Converter interface {
    VideoToAudio(videoPath string) (string,error)
}

type Service struct {
	bot        *telegrambot.Bot
	downloader Downloader
    converter Converter
}

func NewService(
	bot *telegrambot.Bot,
	downloader Downloader,
    converter Converter,
) *Service {
	return &Service{
		bot:        bot,
		downloader: downloader,
        converter: converter,
	}
}

func (s *Service) handleUpdate(
	ctx context.Context,
	update telegrambot.Update,
) error {
	if update.Message != nil {
		if update.Message.Text == "/start" {
			if err := s.handleStartCommand(ctx, update.Message.Chat.ID); err != nil {
				return err
			}
		}
		if update.Message.LinkPreviewOptions != nil {
			var url string
			if len(strings.TrimSpace(update.Message.LinkPreviewOptions.URL)) > 0 {
				url = update.Message.LinkPreviewOptions.URL
			} else {
				url = update.Message.Text
			}
			if err := s.handleMsgWURL(
				ctx,
				update.Message.Chat.ID,
				update.Message.ID,
				url,
			); err != nil {
				return err
			}
		}
	}
	if update.CallbackQuery != nil {
		if err := s.handleCallbackQuery(
			ctx, *update.CallbackQuery,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Run() error {
	var lastUpdateID int64
	ctx := context.Background()
	for {
		params := telegrambot.GetUpdatesParams{
			Offset:         lastUpdateID + 1,
			Limit:          100,
			Timeout:        45,
			AllowedUpdates: []string{"message", "callback_query"},
		}
		updates, err := s.bot.GetUpdates(ctx, params)
		if err != nil {
			return err
		}
		for _, u := range updates {
			lastUpdateID = u.ID
			go func() {
				if err := s.handleUpdate(ctx, u); err != nil {
					slog.ErrorContext(ctx, err.Error())
				}
			}()
		}
	}
}
