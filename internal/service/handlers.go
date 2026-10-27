package service

import (
	"context"
	"errors"
	"strings"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

const (
	StartText         = "🔗 Отправьте ссылку на видео"
	WaitText          = "⏳ Подождите, загружается..."
	DownloadAudioText = "​📥 Скачать аудио"
)

const (
	TryAgainText  = "🔄 Еще раз"
	ErrInternal   = "❌ Внутренняя ошибка, попробуйте снова"
	ErrUrlTooLong = "✂️ Ссылка слишком длинная!"
)

func (s *Service) handleStartCommand(
	ctx context.Context,
	chatID int64,
) error {
	params := telegrambot.SendMessageParams{
		ChatID: chatID,
		Text:   StartText,
	}
	if _, err := s.bot.SendMessage(ctx, params); err != nil {
		return err
	}
	return nil
}

type callbackAction string

const (
	againVideo callbackAction = "av"
	sendAudio  callbackAction = "sa"
	againAudio callbackAction = "aa"
)

func (s *Service) handleMsgWURL(
	ctx context.Context,
	chatID int64,
	messageID int64,
	url string,
) error {
	if len([]byte(url)) > 61 {
		params := telegrambot.SendMessageParams{
			ChatID: chatID,
			Text:   ErrUrlTooLong,
			ReplyParams: &telegrambot.ReplyParameters{
				MessageID: &messageID,
			},
		}
		s.bot.SendMessage(ctx, params)
		return errors.New(ErrUrlTooLong)
	}
	params := telegrambot.SendMessageParams{
		ChatID: chatID,
		Text:   WaitText,
		ReplyParams: &telegrambot.ReplyParameters{
			MessageID: &messageID,
		},
	}
	msg, err := s.bot.SendMessage(ctx, params)
	if err != nil {
		s.sendError(ctx, chatID, messageID, againVideo, url)
		return err
	}
	return s.editToVideo(ctx, chatID, msg.ID, url)
}

func (s *Service) handleCallbackQuery(
	ctx context.Context,
	callback telegrambot.CallbackQuery,
) error {
	data := callback.Data[strings.Index(callback.Data, "-")+1:]
	action := callbackAction(callback.Data[:strings.Index(callback.Data, "-")])
	var message telegrambot.InaccessibleMessage
	switch v := callback.Message.(type) {
	case telegrambot.InaccessibleMessage:
		message = v
	case telegrambot.Message:
		message.ID = v.ID
		message.Chat = v.Chat
	case nil:
		return errors.New("message is nil")
	default:
		return errors.New("unsupported message type")
	}
	switch action {
	case againVideo:
		msg, err := s.bot.EditMessageText(
			ctx,
			telegrambot.EditMessageTextParams{
				ChatID:    message.Chat.ID,
				MessageID: message.ID,
				Text:      WaitText,
			},
		)
		if err != nil {
			s.editToError(ctx, message.Chat.ID, message.ID, againVideo, data)
			return err
		}
		if err := s.editToVideo(ctx, msg.Chat.ID, msg.ID, data); err != nil {
			return err
		}
	case sendAudio:
		markup := telegrambot.InlineKeyboardMarkup{
			Keyboard: [][]telegrambot.InlineKeyboardButton{},
		}
		msg, err := s.bot.EditMessageReplyMarkup(
			ctx,
			telegrambot.EditMessageReplyMarkupParams{
				ChatID:      message.Chat.ID,
				MessageID:   message.ID,
				ReplyMarkup: &markup,
			},
		)
		if err != nil {
			s.editToError(ctx, message.Chat.ID, message.ID, sendAudio, data)
			return err
		}
		msg, err = s.bot.SendMessage(
			ctx,
			telegrambot.SendMessageParams{
				ChatID: msg.Chat.ID,
				Text:   WaitText,
				ReplyParams: &telegrambot.ReplyParameters{
					MessageID: &msg.ID,
				},
			},
		)
		if err != nil {
			s.editToError(ctx, msg.Chat.ID, msg.ID, againAudio, data)
			return err
		}
		if err := s.editToAudio(ctx, msg.Chat.ID, msg.ID, data); err != nil {
			return err
		}
	case againAudio:
		msg, err := s.bot.EditMessageText(
			ctx,
			telegrambot.EditMessageTextParams{
				ChatID:    message.Chat.ID,
				MessageID: message.ID,
				Text:      WaitText,
			},
		)
		if err != nil {
			s.editToError(ctx, message.Chat.ID, message.ID, againAudio, data)
			return err
		}
		if err := s.editToAudio(ctx, msg.Chat.ID, msg.ID, data); err != nil {
			return err
		}
	}
	s.bot.AnswerCallbackQuery(ctx, telegrambot.AnswerCallbackQueryParams{CallbackQueryID: callback.ID})
	return nil
}
