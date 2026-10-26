package service

import (
	"context"
	"fmt"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

func (s *Service) sendVideo(
	ctx context.Context,
	chatId int64,
	messageId int64,
	url string,
) error {
	errText := againVideo
	mediaFilePath, err := s.downloader.DownloadVideo("tmp", url)
	if err != nil {
		s.editToError(
			ctx,
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", errText, url),
		)
		return err
	}
	inputMedia := telegrambot.InputMedia{
		Media: &telegrambot.LocalSource{
			Filepath: *mediaFilePath,
		},
		Type: telegrambot.MediaTypeVideo,
	}
	button := &telegrambot.InlineKeyboardButton{
		Text:         DownloadAudioText,
		CallbackData: fmt.Sprintf("%s-%s", sendAudio, url),
	}
	markup := &telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{{*button}},
	}
	if _, err := s.bot.EditMessageMedia(
		ctx,
		telegrambot.EditMessageMediaParams{
			ChatID:      chatId,
			MessageID:   messageId,
			Media:       inputMedia,
			ReplyMarkup: markup,
		},
	); err != nil {
		s.editToError(
			ctx,
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", errText, url),
		)
		return err
	}
	return nil
}

func (s *Service) sendAudio(
	ctx context.Context,
	chatId int64,
	messageId int64,
	url string,
) error {
	errText := againAudio
	mediaFilePath, err := s.downloader.DownloadAudio("tmp", url)
	if err != nil {
		s.editToError(
			ctx,
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", errText, url),
		)
		return err
	}
	inputMedia := telegrambot.InputMedia{
		Media: &telegrambot.LocalSource{
			Filepath: *mediaFilePath,
		},
		Type: telegrambot.MediaTypeAudio,
	}
	if _, err := s.bot.EditMessageMedia(
		ctx,
		telegrambot.EditMessageMediaParams{
			ChatID:    chatId,
			MessageID: messageId,
			Media:     inputMedia,
		},
	); err != nil {
		s.editToError(
			ctx,
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", errText, url),
		)
		return err
	}
	return nil
}
