package view

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func ErrorMessageParams(chatID int64, replyMsgID int) *telego.SendMessageParams {
	return &telego.SendMessageParams{
		ChatID: tu.ID(chatID),
		Text:   errorText,
		ReplyParameters: &telego.ReplyParameters{
			MessageID: replyMsgID,
		},
	}
}

func ErrorEditMessageParams(chatID int64, msgID int) *telego.EditMessageTextParams {
	return &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   msgID,
		Text:        errorText,
		ReplyMarkup: EmptyMarkup,
	}
}

func ErrorEditMessageMarkupParams(chatID int64, msgID int, callbackData string) *telego.EditMessageTextParams {
	return &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   msgID,
		Text:        errorText,
		ReplyMarkup: errorMarkup(callbackData),
	}
}

func WaitMessageParams(msgID int, chatID int64) *telego.EditMessageTextParams {
	return &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   msgID,
		Text:        waitText,
		ReplyMarkup: EmptyMarkup,
	}
}

func StartMessageParams(chatID int64) *telego.SendMessageParams {
	return &telego.SendMessageParams{
		ChatID: tu.ID(chatID),
		Text:   startText,
	}
}

func ChooseDownloadMessageParams(chatID int64, replyMsgID int, reqID, separator string) *telego.SendMessageParams {
	return &telego.SendMessageParams{
		ChatID:      tu.ID(chatID),
		Text:        chooseDownloadText,
		ReplyMarkup: downloadOptionsMarkup(reqID, separator),
		ReplyParameters: &telego.ReplyParameters{
			MessageID: replyMsgID,
		},
	}
}
