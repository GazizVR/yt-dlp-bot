package service

import (
	"context"
	"errors"
	"fmt"
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
	if _, err := s.Tg.SendMessage(ctx, params); err != nil {
		return err
	}
	return nil
}

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
		s.Tg.SendMessage(ctx, params)
		return errors.New(ErrUrlTooLong)
	}
	params := telegrambot.SendMessageParams{
		ChatID: chatID,
		Text:   WaitText,
		ReplyParams: &telegrambot.ReplyParameters{
			MessageID: &messageID,
		},
	}
	msg, err := s.Tg.SendMessage(ctx, params)
	if err != nil {
		s.sendError(
			ctx,
			chatID,
			messageID,
			fmt.Sprintf("%s-%s", againVideo, url),
		)
		return err
	}
	return s.sendVideo(ctx, chatID, msg.ID, url)
}

func (s *Service) handleCallbackQuery(
	ctx context.Context,
	callback telegrambot.CallbackQuery,
) error {
	url := callback.Data[strings.Index(callback.Data, "-")+1:]
	action := callback.Data[:strings.Index(callback.Data, "-")]
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
	case sendAudio:
		markup := telegrambot.InlineKeyboardMarkup{
			Keyboard: [][]telegrambot.InlineKeyboardButton{},
		}
		s.Tg.EditMessageReplyMarkup(
			ctx,
			telegrambot.EditMessageReplyMarkupParams{
				ChatID:      message.Chat.ID,
				MessageID:   message.ID,
				ReplyMarkup: &markup,
			},
		)
		msg, _ := s.Tg.SendMessage(
			ctx,
			telegrambot.SendMessageParams{
				ChatID: message.Chat.ID,
				Text:   WaitText,
				ReplyParams: &telegrambot.ReplyParameters{
					MessageID: &message.ID,
				},
			},
		)
		s.sendAudio(
			ctx,
			msg.Chat.ID,
			msg.ID,
			url,
		)
	case againVideo:
		s.Tg.EditMessageText(
			ctx,
			telegrambot.EditMessageTextParams{
				ChatID:    message.Chat.ID,
				MessageID: message.ID,
				Text:      WaitText,
			},
		)
		s.sendVideo(
			ctx,
			message.Chat.ID,
			message.ID,
			url,
		)
	case againAudio:
		s.Tg.EditMessageText(
			ctx,
			telegrambot.EditMessageTextParams{
				ChatID:    message.Chat.ID,
				MessageID: message.ID,
				Text:      WaitText,
			},
		)
		s.sendAudio(
			ctx,
			message.Chat.ID,
			message.ID,
			url,
		)
	}
	s.Tg.AnswerCallbackQuery(ctx, telegrambot.AnswerCallbackQueryParams{CallbackQueryID: callback.ID})
	return nil
}
