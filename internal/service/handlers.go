package service

const StartText = "Hello, world!"

func (s *Service) handleStartCommand(
	chatId int64,
) error {
	if _, err := s.Cl.SendMessage(
		chatId,
		StartText,
	); err != nil {
		return err
	}
	return nil
}
