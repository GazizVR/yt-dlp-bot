package service

import "bot/pkg/telegram"

type Service struct {
	Client *telegram.Client
}

func NewService(client *telegram.Client) *Service {
	return &Service{
		Client: client,
	}
}

func (s *Service) StartUpdateHandle() error {
	return nil
}
