package service

import (
	"context"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

const (
	sendAudio  = "sa"
	againVideo = "av"
	againAudio = "aa"
)

func (s *Service) sendError(
	ctx context.Context,
	chatId int64,
	msgToReply int64,
	callbackData string,
) {
	button := &telegrambot.InlineKeyboardButton{
		Text:         TryAgainText,
		CallbackData: callbackData,
	}
	params := telegrambot.SendMessageParams{
		ChatID: chatId,
		Text:   ErrInternal,
		ReplyParams: &telegrambot.ReplyParameters{
			MessageID: &msgToReply,
		},
		ReplyMarkup: telegrambot.InlineKeyboardMarkup{
			Keyboard: [][]telegrambot.InlineKeyboardButton{{*button}},
		},
	}
	s.Tg.SendMessage(ctx, params)
}

func (s *Service) editToError(
	ctx context.Context,
	chatId int64,
	messageId int64,
	callbackData string,
) {
	button := &telegrambot.InlineKeyboardButton{
		Text:         TryAgainText,
		CallbackData: callbackData,
	}
	markup := telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{{*button}},
	}
	params := telegrambot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   messageId,
		Text:        ErrInternal,
		ReplyMarkup: &markup,
	}
	s.Tg.EditMessageText(ctx, params)
}
