package view

import (
	"github.com/gazizvr/yt-dlp-bot/internal/domain"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

var EmptyMarkup = &telego.InlineKeyboardMarkup{
	InlineKeyboard: [][]telego.InlineKeyboardButton{},
}

func errorMarkup(callbackData string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(tryAgainBtnText).WithCallbackData(callbackData),
		),
	)
}

func downloadOptionsMarkup(reqID, separator string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(videoBtnText).WithCallbackData(reqID+separator+string(domain.DownloadVideo)),
			tu.InlineKeyboardButton(audioBtnText).WithCallbackData(reqID+separator+string(domain.DownloadAudio)),
		),
	)
}
