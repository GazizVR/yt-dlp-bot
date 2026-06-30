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
    if update.Message != nil {
        if update.Message.Text == "/start" {
		    if err := s.handleStartCommand(update.Message.Chat.Id); err != nil {
			    return err
		    }
	    }
	    if update.Message.LinkPreview != nil {
		    if err := s.handleMsgWURL(
			    update.Message.Chat.Id,
			    update.Message.LinkPreview.URL,
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

func (s *Service) handleUpdates(
	lastUpdateId *int64,
) error {
	response, err := s.Tg.GetUpdates(
		*lastUpdateId,
		100,
		60,
		[]string{"message","callback_query"},
	)
	if err != nil {
		return err
	}
	for _, u := range response.Result {
		*lastUpdateId = u.Id
		go s.handleUpdate(u)
	}
	*lastUpdateId += 1
	return nil
}

func (s *Service) Run() error {
	var lastUpdateId int64
	for {
		if err := s.handleUpdates(&lastUpdateId); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(1 * time.Second)
	}
}
