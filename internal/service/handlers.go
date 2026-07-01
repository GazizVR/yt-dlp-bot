package service

import (
	"bot/pkg/telegram"
	"fmt"
	"strings"
)

const (
	StartText         = "🔗 Отправьте ссылку на видео"
	WaitText          = "⏳ Подождите, загружается..."
	DownloadAudioText = "​📥 Скачать аудио"
	ErrorText         = "❌ Внутренняя ошибка, попробуйте снова"
	TryAgainText      = "🔄 Еще раз"
)

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
		WaitText,
		nil,
		&messageId,
	)
	if err != nil {
		s.sendError(
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", againVideo, url),
		)
		return err
	}
	if err := s.sendMedia(
		chatId,
		msg.Result.Id,
		true,
		url,
	); err != nil {
		return err
	}
	return nil
}

func (s *Service) handleCallbackQuery(
	callback telegram.CallbackQuery,
) error {
	url := callback.Data[strings.Index(callback.Data, "-")+1:]
	action := callback.Data[:strings.Index(callback.Data, "-")]
	switch action {
	case sendAudio:
		markup := telegram.NewInlineMarkup([]telegram.InlineButton{})
		s.Tg.EditMessageReplyMarkup(
			callback.Message.Chat.Id,
			callback.Message.Id,
			*markup,
		)
		msg, _ := s.Tg.SendMessage(
			callback.Message.Chat.Id,
			WaitText,
			nil,
			&callback.Message.Id,
		)
		s.sendMedia(
			msg.Result.Chat.Id,
			msg.Result.Id,
			false,
			url,
		)
	case againVideo:
		s.Tg.EditMessageText(
			callback.Message.Chat.Id,
			callback.Message.Id,
			WaitText,
			nil,
		)
		s.sendMedia(
			callback.Message.Chat.Id,
			callback.Message.Id,
			true,
			url,
		)
	case againAudio:
		s.Tg.EditMessageText(
			callback.Message.Chat.Id,
			callback.Message.Id,
			WaitText,
			nil,
		)
		s.sendMedia(
			callback.Message.Chat.Id,
			callback.Message.Id,
			false,
			url,
		)
	}
	s.Tg.AnswerCallbackQuery(callback.Id)
	return nil
}
