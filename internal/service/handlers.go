package service

const (
	StartText = "🔗 Отправьте ссылку на видео"
	SendText  = "⏳ Подождите, загружаем..."
	ErrorText = "❌ Ошибка установки, попробуйте снова"
)

func (s *Service) handleStartCommand(
	chatId int64,
) error {
	if _, err := s.Tg.SendMessage(
		chatId,
		StartText,
	); err != nil {
		return err
	}
	return nil
}

func (s *Service) handleMsgWURL(
	chatId int64,
	url string,
) (*int64, error) {
	msg, err := s.Tg.SendMessage(
		chatId,
		SendText,
	)
	if err != nil {
		return &msg.Result.Id, err
	}
	videoFile, err := s.Dlp.DownloadVideo(
		"tmp",
		url,
	)
	if err != nil {
		return &msg.Result.Id, err
	}
	s.Tg.DeleteMessage(chatId, msg.Result.Id)
	_, err = s.Tg.SendVideo(chatId, *videoFile)
	if err != nil {
		return &msg.Result.Id, err
	}
	return nil, nil
}
