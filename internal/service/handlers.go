package service

import (
	"bot/pkg/telegram"
)

const (
	StartText         = "🔗 Отправьте ссылку на видео"
	SendText          = "⏳ Подождите, загружаем..."
	DownloadAudioText = "​📥 Скачать аудио"
	ErrorText         = "❌ Внутренняя ошибка, попробуйте снова"
	TryAgainText      = "🔄 Еще раз"
)

func (s *Service) sendError(
	chatId int64,
	msgToReply int64,
) {
	button := &telegram.InlineButton{
		Text: TryAgainText,
		Data: "Try again",
	}
	markup := telegram.NewInlineMarkup(
		[]telegram.InlineButton{*button},
	)
	s.Tg.SendMessage(
		chatId,
		ErrorText,
		markup,
		&msgToReply,
	)
}

func (s *Service) editToError(
	chatId int64,
	messageId int64,
) {
	button := &telegram.InlineButton{
		Text: TryAgainText,
		Data: "Try again",
	}
	markup := telegram.NewInlineMarkup(
		[]telegram.InlineButton{*button},
	)
	s.Tg.EditMessageText(
		chatId,
		messageId,
		ErrorText,
		markup,
	)
}

func (s *Service) handleStartCommand(
	chatId int64,
) error {
	if _, err := s.Tg.SendMessage(
		chatId,
		StartText,
		nil,
		nil,
	); err != nil {
		return err
	}
	return nil
}

func (s *Service) handleMsgWURL(
	chatId int64,
	messageId int64,
	url string,
) error {
	msg, err := s.Tg.SendMessage(
		chatId,
		SendText,
		nil,
		&messageId,
	)
	if err != nil {
		s.sendError(chatId, messageId)
		return err
	}
	videoFile, err := s.Dlp.DownloadVideo(
		"tmp",
		url,
	)
	if err != nil {
		s.editToError(chatId, msg.Result.Id)
		return err
	}
	button := &telegram.InlineButton{
		Text: DownloadAudioText,
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
		s.editToError(chatId, msg.Result.Id)
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
		s.sendError(
			callback.Message.Chat.Id,
			callback.Message.Id,
		)
		return err
	}
	_, err = s.Tg.SendAudio(
		callback.Message.Chat.Id,
		*audio,
		&callback.Message.Id,
	)
	if err != nil {
		s.sendError(
			callback.Message.Chat.Id,
			callback.Message.Id,
		)
		return err
	}
	s.Tg.AnserCallbackQuery(callback.Id)
	return nil
}
