package service

import (
	"context"
	"fmt"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

func (s *Service) editToVideo(
	ctx context.Context,
	chatID int64,
	messageID int64,
	url string,
) error {
	mediaPath, err := s.downloader.DownloadVideo(url, s.filesDir)
	if err != nil {
		s.editToError(ctx, chatID, messageID, againVideo, url)
		return err
	}
	inputMedia := telegrambot.InputMedia{
		Media: &telegrambot.LocalSource{
			Filepath: mediaPath,
		},
		Type: telegrambot.MediaTypeVideo,
	}
	button := telegrambot.InlineKeyboardButton{
		Text:         DownloadAudioText,
		CallbackData: fmt.Sprintf("%s-%s", sendAudio, mediaPath),
	}
	markup := &telegrambot.InlineKeyboardMarkup{
		Keyboard: [][]telegrambot.InlineKeyboardButton{{button}},
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
		s.editToError(ctx, chatID, messageID, againVideo, url)
		return err
	}
	return nil
}

func (s *Service) editToAudio(
	ctx context.Context,
	chatID int64,
	messageID int64,
	videoPath string,
	replyMsgID int64,
) error {
	audioPath, err := s.converter.VideoToAudio(videoPath, s.filesDir)
	if err != nil {
		s.editToError(ctx, chatID, messageID, againAudio, videoPath)
		return err
	}
	if err := s.bot.DeleteMessage(
		ctx,
		telegrambot.DeleteMessageParams{
			ChatID:    chatID,
			MessageID: messageID,
		},
	); err != nil {
		s.editToError(ctx, chatID, messageID, againAudio, videoPath)
		return err
	}
	if _, err := s.bot.SendVoice(
		ctx,
		telegrambot.SendVoiceParams{
			ChatID:      chatID,
			Voice:       &telegrambot.LocalSource{Filepath: audioPath},
			ReplyParams: &telegrambot.ReplyParameters{MessageID: &replyMsgID},
		},
	); err != nil {
		s.editToError(ctx, chatID, messageID, againAudio, videoPath)
		return err
	}
	return nil
}
