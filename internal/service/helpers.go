package service

import (
	"bot/pkg/telegram"
	"fmt"
	"os"
)

func (s *Service) sendMedia(
	chatId int64,
	messageId int64,
	isVideo bool,
	url string,
) error {
	var (
		mediaFile     *os.File
		err           error
		markup        *telegram.InlineMarkup
		formFieldName string
		errText       string
	)
	if isVideo {
		errText = againVideo
		formFieldName = "video"
		mediaFile, err = s.Dlp.DownloadVideo(
			"tmp",
			url,
		)
		button := &telegram.InlineButton{
			Text: DownloadAudioText,
			Data: fmt.Sprintf("%s-%s", sendAudio, url),
		}
		markup = telegram.NewInlineMarkup(
			[]telegram.InlineButton{*button},
		)
	} else {
		errText = againAudio
		formFieldName = "audio"
		mediaFile, err = s.Dlp.DownloadAudio(
			"tmp",
			url,
		)
		markup = nil
	}

	if err != nil {
		s.editToError(
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", errText, url),
		)
		return err
	}
	_, err = s.Tg.EditMessageMedia(
		chatId,
		messageId,
		formFieldName,
		*mediaFile,
		markup,
	)
	if err != nil {
		s.editToError(
			chatId,
			messageId,
			fmt.Sprintf("%s-%s", errText, url),
		)
		return err
	}
	return nil
}
