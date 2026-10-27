package view

import (
	telegrambot "github.com/gazizvr/go-telegram-bot"
)

func DownloadAudioMarkup(data string) *telegrambot.InlineKeyboardMarkup {
	button := telegrambot.InlineKeyboardButton{
		Text:         "​📥 Скачать аудио",
		CallbackData: data,
	}
	return &telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{{button}},
	}
}

var EmptyMarkup = &telegrambot.InlineKeyboardMarkup{
	Keyboard: [][]telegrambot.InlineKeyboardButton{},
}

func errorMarkup(data string) *telegrambot.InlineKeyboardMarkup {
	button := &telegrambot.InlineKeyboardButton{
		Text:         "🔄 Еще раз",
		CallbackData: data,
	}
	return &telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{{*button}},
	}
}
