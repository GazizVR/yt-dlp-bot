package service

import (
	"strings"
	"time"

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
) error {
	if update.Message != nil {
		if update.Message.Text == "/start" {
			if err := s.handleStartCommand(update.Message.Chat.Id); err != nil {
				return err
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
				return err
			}
		}
	}
	if update.Callback != nil {
		if err := s.handleCallbackQuery(
			*update.Callback,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Run() error {
	var lastUpdateId int64
	for {
		response, err := s.Tg.GetUpdates(
			lastUpdateId,
			100,
			20,
			[]string{"message", "callback_query"},
		)
		if err != nil {
			time.Sleep(800 * time.Millisecond)
			continue
		}
		for _, u := range response.Result {
			lastUpdateId = u.Id + 1
			go s.handleUpdate(u)
		}
	}
}
