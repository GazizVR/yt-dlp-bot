package service

import (
	"bot/pkg/telegram"
)

const (
	StartText  = "🔗 Отправьте ссылку на видео"
	SendText   = "⏳ Подождите, загружаем..."
	ErrorText  = "❌ Внутренняя ошибка, попробуйте снова"
	ButtonText = "​📥 Скачать аудио"
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
		s.Tg.EditMessageText(
			chatId,
			msg.Result.Id,
			ErrorText,
		)
		return err
	}
	_, err = s.Tg.EditMessageToVideo(
		chatId,
		msg.Result.Id,
		*videoFile,
		ButtonText,
		url,
	)
	if err != nil {
		s.Tg.EditMessageText(
			chatId,
			msg.Result.Id,
			ErrorText,
		)
		return err
	}
	return nil
}

func (s *Service) handleCallbackQuery(
	callback telegram.CallbackQuery,
) error {
	url := callback.Data
	s.Tg.AnserCallbackQuery(callback.Id)
	s.Tg.DeleteVideoKeyboard(
		callback.Message.Chat.Id,
		callback.Message.Id,
	)
	audio, err := s.Dlp.DownloadAudio("tmp", url)
	if err != nil {
		s.Tg.SendMessage(
			callback.Message.Chat.Id,
			ErrorText,
		)
		return err
	}
	_, err = s.Tg.SendAudio(callback.Message.Chat.Id, *audio)
	if err != nil {
		s.Tg.SendMessage(
			callback.Message.Chat.Id,
			ErrorText,
		)
		return err
	}
	return nil
}
