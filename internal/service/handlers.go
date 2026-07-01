package service

import (
	"bot/pkg/telegram"
	"fmt"
	"strings"
)

const (
	StartText         = "🔗 Отправьте ссылку на видео"
	SendText          = "⏳ Подождите, загружаем..."
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
		SendText,
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
	videoFile, err := s.Dlp.DownloadVideo(
		"tmp",
		url,
	)
	if err != nil {
		s.editToError(
			chatId,
			msg.Result.Id,
			fmt.Sprintf("%s-%s", againVideo, url),
		)
		return err
	}
	button := &telegram.InlineButton{
		Text: DownloadAudioText,
		Data: fmt.Sprintf("%s-%s", sendAudio, url),
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
		s.editToError(
			chatId,
			msg.Result.Id,
			fmt.Sprintf("%s-%s", againVideo, url),
		)
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
		audio, err := s.Dlp.DownloadAudio("tmp", url)
		if err != nil {
			s.sendError(
				callback.Message.Chat.Id,
				callback.Message.Id,
				fmt.Sprintf("%s-%s", againAudio, url),
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
				fmt.Sprintf("%s-%s", againAudio, url),
			)
			return err
		}
	}
	s.Tg.AnserCallbackQuery(callback.Id)
	return nil
}
