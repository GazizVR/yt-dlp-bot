package view

import (
	telegrambot "github.com/gazizvr/go-telegram-bot"
	"github.com/gazizvr/yt-dlp-bot/internal/domain"
)

var EmptyMarkup = &telegrambot.InlineKeyboardMarkup{
	Keyboard: [][]telegrambot.InlineKeyboardButton{},
}

func errorMarkup(callbackData string) *telegrambot.InlineKeyboardMarkup {
	button := &telegrambot.InlineKeyboardButton{
		Text:         tryAgainBtnText,
		CallbackData: callbackData,
	}
	return &telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{{*button}},
	}
}

func downloadOptionsMarkup(reqID, separator string) *telegrambot.InlineKeyboardMarkup {
	videoButton := telegrambot.InlineKeyboardButton{
		Text:         videoBtnText,
		CallbackData: reqID + separator + string(domain.SendVideo),
	}
	audioButton := telegrambot.InlineKeyboardButton{
		Text:         audioBtnText,
		CallbackData: reqID + separator + string(domain.SendAudio),
	}
	return &telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{
			{videoButton, audioButton},
		},
	}
}
