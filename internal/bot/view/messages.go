package view

import telegrambot "github.com/gazizvr/go-telegram-bot"

func ErrorMessageParams(chatID, replyMsgID int64) telegrambot.SendMessageParams {
	return telegrambot.SendMessageParams{
		ChatID: chatID,
		Text:   errorText,
		ReplyParams: &telegrambot.ReplyParameters{
			MessageID: &replyMsgID,
		},
	}
}

func ErrorEditMessageParams(chatID, msgID int64) telegrambot.EditMessageTextParams {
	return telegrambot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   msgID,
		Text:        errorText,
		ReplyMarkup: EmptyMarkup,
	}
}

func ErrorEditMessageMarkupParams(chatID, msgID int64, callbackData string) telegrambot.EditMessageTextParams {
	return telegrambot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   msgID,
		Text:        errorText,
		ReplyMarkup: errorMarkup(callbackData),
	}
}

func WaitMessageParams(msgID, chatID int64) telegrambot.EditMessageTextParams {
	return telegrambot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   msgID,
		Text:        waitText,
		ReplyMarkup: EmptyMarkup,
	}
}

func StartMessageParams(chatID int64) telegrambot.SendMessageParams {
	return telegrambot.SendMessageParams{
		ChatID: chatID,
		Text:   startText,
	}
}

func ChooseDownloadMessageParams(chatID, replyMsgID int64, reqID, separator string) telegrambot.SendMessageParams {
	return telegrambot.SendMessageParams{
		ChatID:      chatID,
		Text:        chooseDownloadText,
		ReplyMarkup: downloadOptionsMarkup(reqID, separator),
		ReplyParams: &telegrambot.ReplyParameters{
			MessageID: &replyMsgID,
		},
	}
}
