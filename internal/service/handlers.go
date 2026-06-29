package service

import "fmt"

const (
	StartText = "🔗 Отправьте ссылку на видео"
	SendText  = "⏳ Подождите, загружаем..."
	ErrorText = "❌ Ошибка установки, попробуйте снова"
    ButtonText = ""
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
	msg, err := s.Tg.SendMessage(
		chatId,
		SendText,
	)
	if err != nil {
		s.Tg.SendMessage(
			chatId,
			ErrorText,
		)
		return err
	}
	videoFile, err := s.Dlp.DownloadVideo(
		"tmp",
		url,
	)
	if err != nil {
		s.Tg.DeleteMessage(chatId, msg.Result.Id)
		s.Tg.SendMessage(
			chatId,
			ErrorText,
		)
		return err
	}
	s.Tg.DeleteMessage(chatId, msg.Result.Id)
	_, err = s.Tg.SendVideoWithButton(
        chatId,
        *videoFile,
        ButtonText,
        "Download audio",
    )
	if err != nil {
		s.Tg.DeleteMessage(chatId, msg.Result.Id)
		s.Tg.SendMessage(
			chatId,
			ErrorText,
		)
		return err
	}
	return nil
}
