package service

import (
	"bot/pkg/telegram"
	"bot/pkg/ytdlp"
	"time"
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
	if update.Message.Text == "/start" {
		if err := s.handleStartCommand(update.Message.Chat.Id); err != nil {
			return err
		}
	}
	if update.Message.LinkPreview != nil {
		if msgId, err := s.handleMsgWURL(
			update.Message.Chat.Id,
			update.Message.LinkPreview.URL,
		); err != nil {
			if msgId != nil {
				s.Tg.EditMessageText(
					update.Message.Chat.Id,
					*msgId,
					ErrorText,
				)
			} else {
				s.Tg.SendMessage(
					update.Message.Chat.Id,
					ErrorText,
				)
			}
			return err
		}
	}
	return nil
}

func (s *Service) handleUpdates(
	offset int64,
) (*int64, error) {
	lastUpdateId := offset
	response, err := s.Tg.GetUpdates(
		offset,
		100,
		60,
		[]string{"message"},
	)
	if err != nil {
		return nil, err
	}
	for _, u := range response.Result {
		lastUpdateId = u.Id
		go s.handleUpdate(u)
	}
	lastUpdateId++
	return &lastUpdateId, nil
}

func (s *Service) Run() error {
	var lastUpdateId int64
	for {
		updateId, err := s.handleUpdates(lastUpdateId)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		lastUpdateId = *updateId
		time.Sleep(1 * time.Second)
	}
}
