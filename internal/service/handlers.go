package service

import "fmt"

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
) error {
	if _, err := s.Tg.SendMessage(
		chatId,
		SendText,
	); err != nil {
		return err
	}
	videoFile, err := s.Dlp.DownloadVideo(
		fmt.Sprint("tmp/fine-", chatId),
		url,
	)
	if err != nil {
		return err
	}
	_, err = s.Tg.SendVideo(chatId, *videoFile)
	if err != nil {
		return err
	}
	return nil
}
