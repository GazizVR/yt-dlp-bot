package view

import telegrambot "github.com/gazizvr/go-telegram-bot"

const errorText string = "❌ Внутренняя ошибка, попробуйте снова"

func ErrorMessageParams(chatID, replyMsgID int64, data string) telegrambot.SendMessageParams {
	return telegrambot.SendMessageParams{
		ChatID:      chatID,
		Text:        errorText,
		ReplyMarkup: errorMarkup(data),
		ReplyParams: &telegrambot.ReplyParameters{
			MessageID: &replyMsgID,
		},
	}
}

func ErrorEditMessageParams(chatID, msgID int64, data string) telegrambot.EditMessageTextParams {
	return telegrambot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   msgID,
		Text:        errorText,
		ReplyMarkup: errorMarkup(data),
	}
}

const waitText string = "⏳ Подождите, загружается..."

func WaitMessageParams(chatID, replyMsgID int64) telegrambot.SendMessageParams {
	return telegrambot.SendMessageParams{
		ChatID: chatID,
		Text:   string(waitText),
		ReplyParams: &telegrambot.ReplyParameters{
			MessageID: &replyMsgID,
		},
	}
}

func WaitEditMessageParams(msgID, chatID int64) telegrambot.EditMessageTextParams {
	return telegrambot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   msgID,
		Text:        string(waitText),
		ReplyMarkup: EmptyMarkup,
	}
}

func StartMessageParams(chatID int64) telegrambot.SendMessageParams {
	return telegrambot.SendMessageParams{
		ChatID: chatID,
		Text:   "🔗 Отправьте ссылку на видео",
	}
}
