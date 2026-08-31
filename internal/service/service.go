package service

import (
	"strings"

	"github.com/gazizvr/yt-dlp-bot/pkg/telegram"
	"github.com/gazizvr/yt-dlp-bot/pkg/ytdlp"
)

type Service struct {
	Tg  *telegram.Client
	Dlp *ytdlp.Client
}

func NewService(
	tg *telegram.Client,
	dlp *ytdlp.Client,
) *Service {
	return &Service{
		Tg:  tg,
		Dlp: dlp,
	}
}

func (s *Service) handleUpdate(
	update telegram.Update,
) {
	if update.Message != nil {
		if update.Message.Text == "/start" {
			if err := s.handleStartCommand(update.Message.Chat.Id); err != nil {
				return
			}
		}
		if update.Message.LinkPreview != nil {
			var url string
			if len(strings.TrimSpace(update.Message.LinkPreview.URL)) > 0 {
				url = update.Message.LinkPreview.URL
			} else {
				url = update.Message.Text
			}
			if err := s.handleMsgWURL(
				update.Message.Chat.Id,
				update.Message.Id,
				url,
			); err != nil {
				return
			}
		}
	}
	if update.Callback != nil {
		if err := s.handleCallbackQuery(
			*update.Callback,
		); err != nil {
			return
		}
	}
}

func (s *Service) Run() error {
	if err := s.Tg.ListenUpdates(
		s.handleUpdate,
		[]string{"message", "callback_query"},
	); err != nil {
		return err
	}
	return nil
}
