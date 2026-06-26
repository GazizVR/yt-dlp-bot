package service

import (
	"bot/pkg/telegram"
	"time"
)

type Service struct {
	Cl *telegram.Client
}

func NewService(client *telegram.Client) *Service {
	return &Service{
		Cl: client,
	}
}

const StartText = "Hello, world!"

func (s *Service) handleUpdate(
	update telegram.Update,
) error {
	if update.Message.Text == "/start" {
		if _, err := s.Cl.SendMessage(
			update.Message.Chat.Id,
			StartText,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) handleUpdates(
	offset int64,
) (*int64, error) {
	lastUpdateId := offset
	response, err := s.Cl.GetUpdates(
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
		if err := s.handleUpdate(u); err != nil {
			return nil, err
		}
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
