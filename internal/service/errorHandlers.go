package service

import "bot/pkg/telegram"

const (
	sendAudio  = "sendaudio"
	againVideo = "againVideo"
	againAudio = "againAudio"
)

func (s *Service) sendError(
	chatId int64,
	msgToReply int64,
	callbackData string,
) {
	button := &telegram.InlineButton{
		Text: TryAgainText,
		Data: callbackData,
	}
	markup := telegram.NewInlineMarkup(
		[]telegram.InlineButton{*button},
	)
	s.Tg.SendMessage(
		chatId,
		ErrorText,
		markup,
		&msgToReply,
	)
}

func (s *Service) editToError(
	chatId int64,
	messageId int64,
	callbackData string,
) {
	button := &telegram.InlineButton{
		Text: TryAgainText,
		Data: callbackData,
	}
	markup := telegram.NewInlineMarkup(
		[]telegram.InlineButton{*button},
	)
	s.Tg.EditMessageText(
		chatId,
		messageId,
		ErrorText,
		markup,
	)
}
