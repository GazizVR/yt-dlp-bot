package service

import (
	"context"
	"fmt"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

func (s *Service) sendError(
	ctx context.Context,
	chatId int64,
	msgToReply int64,
    action string,
	data string,
) {
	button := &telegrambot.InlineKeyboardButton{
		Text:         TryAgainText,
		CallbackData: fmt.Sprintf("%s-%s",action,data),
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
	s.bot.SendMessage(ctx, params)
}

func (s *Service) editToError(
	ctx context.Context,
	chatId int64,
	messageId int64,
    action string,
	data string,
) {
	button := &telegrambot.InlineKeyboardButton{
		Text:         TryAgainText,
		CallbackData: fmt.Sprintf("%s-%s",action,data),
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
	s.bot.EditMessageText(ctx, params)
}
