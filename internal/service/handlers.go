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
		nil,
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
		nil,
	)
	if err != nil {
		s.Tg.SendMessage(
			chatId,
			ErrorText,
			nil,
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
	button := &telegram.InlineButton{
		Text: ButtonText,
		Data: url,
	}
	markup := telegram.NewInlineMarkup(
		[]telegram.InlineButton{*button},
	)
	_, err = s.Tg.EditMessageMedia(
		chatId,
		msg.Result.Id,
		*videoFile,
		markup,
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
	markup := telegram.NewInlineMarkup([]telegram.InlineButton{})
	s.Tg.EditMessageReplyMarkup(
		callback.Message.Chat.Id,
		callback.Message.Id,
		*markup,
	)
	audio, err := s.Dlp.DownloadAudio("tmp", url)
	if err != nil {
		s.Tg.SendMessage(
			callback.Message.Chat.Id,
			ErrorText,
			nil,
		)
		return err
	}
	_, err = s.Tg.SendAudio(
		callback.Message.Chat.Id,
		*audio,
		&callback.Message.Id,
	)
	if err != nil {
		s.Tg.SendMessage(
			callback.Message.Chat.Id,
			ErrorText,
			nil,
		)
		return err
	}
	s.Tg.AnserCallbackQuery(callback.Id)
	return nil
}
