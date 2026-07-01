package service

import (
	"bot/pkg/telegram"
	"fmt"
	"os"
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

func (s *Service) sendMedia(
	chatId int64,
	messageId int64,
	isVideo bool,
	url string,
) error {
	var (
		mediaFile     *os.File
		err           error
		markup        *telegram.InlineMarkup
		formFieldName string
		errText       string
	)
	if isVideo {
		errText = againVideo
		formFieldName = "video"
		mediaFile, err = s.Dlp.DownloadVideo(
			"tmp",
			url,
		)
		button := &telegram.InlineButton{
			Text: DownloadAudioText,
			Data: fmt.Sprintf("%s-%s", sendAudio, url),
		}
		markup = telegram.NewInlineMarkup(
			[]telegram.InlineButton{*button},
		)
	} else {
		errText = againAudio
		formFieldName = "audio"
		mediaFile, err = s.Dlp.DownloadAudio(
			"tmp",
			url,
		)
		markup = nil
	}

	if err != nil {
		s.editToError(
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", errText, url),
		)
		return err
	}
	_, err = s.Tg.EditMessageMedia(
		chatId,
		messageId,
		formFieldName,
		*mediaFile,
		markup,
	)
	if err != nil {
		s.editToError(
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", errText, url),
		)
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
	s.Tg.AnserCallbackQuery(callback.Id)
	return nil
}
