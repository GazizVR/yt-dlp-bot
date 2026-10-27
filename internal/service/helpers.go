package service

import (
	"context"
	"fmt"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

func (s *Service) sendVideo(
	ctx context.Context,
	chatID int64,
	messageID int64,
	url string,
) error {
	mediaPath, err := s.downloader.DownloadVideo(url,"tmp")
	if err != nil {
		s.editToError(ctx,chatID,messageID,againVideo, url)
		return err
	}
	inputMedia := telegrambot.InputMedia{
		Media: &telegrambot.LocalSource{
			Filepath: mediaPath,
		},
		Type: telegrambot.MediaTypeVideo,
	}
	button := &telegrambot.InlineKeyboardButton{
		Text:         DownloadAudioText,
		CallbackData: fmt.Sprintf("%s-%s", sendAudio, mediaPath),
	}
	markup := &telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{{*button}},
	}
	if _, err := s.bot.EditMessageMedia(
		ctx,
		telegrambot.EditMessageMediaParams{
			ChatID:      chatID,
			MessageID:   messageID,
			Media:       inputMedia,
			ReplyMarkup: markup,
		},
	); err != nil {
		s.editToError(ctx,chatID,messageID,againVideo, url)
		return err
	}
	return nil
}

func (s *Service) sendAudio(
	ctx context.Context,
	chatID int64,
	messageID int64,
	mediaPath string,
) error {
    markup := telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{},
	}
	msg, err := s.bot.EditMessageReplyMarkup(
		ctx,
		telegrambot.EditMessageReplyMarkupParams{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: &markup,
		},
	)
    if err != nil {
        s.editToError(ctx,chatID,messageID,sendAudio,mediaPath)
        return err
    }
	msg, err = s.bot.SendMessage(
		ctx,
		telegrambot.SendMessageParams{
			ChatID: msg.Chat.ID,
			Text:   WaitText,
			ReplyParams: &telegrambot.ReplyParameters{
				MessageID: &msg.ID,
			},
		},
	)
    if err != nil {
        s.editToError(ctx,chatID,messageID,sendAudio, mediaPath)
        return err
    }
	audioPath, err := s.converter.VideoToAudio(mediaPath)
	if err != nil {
		s.editToError(ctx,chatID,messageID,sendAudio, mediaPath)
		return err
	}
	inputMedia := telegrambot.InputMedia{
		Media: &telegrambot.LocalSource{
			Filepath: audioPath,
		},
		Type: telegrambot.MediaTypeAudio,
	}
	if _, err := s.bot.EditMessageMedia(
		ctx,
		telegrambot.EditMessageMediaParams{
			ChatID:    chatID,
			MessageID: messageID,
			Media:     inputMedia,
		},
	); err != nil {
		s.editToError(ctx,chatID,messageID,sendAudio, mediaPath)
		return err
	}
	return nil
}
