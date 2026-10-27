package bot

import (
	"context"

	telegrambot "github.com/gazizvr/go-telegram-bot"
	"github.com/gazizvr/yt-dlp-bot/internal/bot/view"
)

type MediaService interface {
	DownloadVideo(url, outDir string) (string, error)
	VideoToAudio(videoPath, outDir string) (string, error)
}

type Handler struct {
	bot      *telegrambot.Bot
	src      MediaService
	filesDir string
}

func NewHandler(bot *telegrambot.Bot, fileDir string) *Handler {
	return &Handler{
		bot:      bot,
		filesDir: fileDir,
	}
}

func (h *Handler) OnStart(
	ctx context.Context,
	upd telegrambot.Update,
) error {
	params := view.StartMessageParams(upd.Message.Chat.ID)
	if _, err := h.bot.SendMessage(ctx, params); err != nil {
		return err
	}
	return nil
}

func (h *Handler) OnURL(
	ctx context.Context,
	upd telegrambot.Update,
) error {
	url := upd.Message.LinkPreviewOptions.URL
	if url == "" {
		url = upd.Message.Text
	}
	params := view.WaitMessageParams(upd.Message.Chat.ID, upd.Message.ID)
	msg, err := h.bot.SendMessage(ctx, params)
	if err != nil {
		return err
	}
	videoPath, err := h.src.DownloadVideo(url, h.filesDir)
	if err != nil {
		return err
	}
	inputMedia := telegrambot.InputMedia{
		Media: &telegrambot.LocalSource{
			Filepath: videoPath,
		},
		Type: telegrambot.MediaTypeVideo,
	}
	if _, err := h.bot.EditMessageMedia(
		ctx,
		telegrambot.EditMessageMediaParams{
			ChatID:      msg.Chat.ID,
			MessageID:   msg.ID,
			Media:       inputMedia,
			ReplyMarkup: view.DownloadAudioMarkup(),
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) OnCallbackQuery(
	ctx context.Context,
	upd telegrambot.Update,
) error {
	var chatID int64 = 0
	var messageID int64 = 0
	msg, err := h.bot.EditMessageReplyMarkup(
		ctx,
		telegrambot.EditMessageReplyMarkupParams{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: view.EmptyMarkup,
		},
	)
	if err != nil {
		return err
	}
	params := view.WaitMessageParams(msg.Chat.ID, msg.ID)
	msg, err = h.bot.SendMessage(ctx, params)
	if err != nil {
		return err
	}
	audioPath, err := h.src.VideoToAudio(videoPath, h.filesDir)
	if err != nil {
		return err
	}
	if err := h.bot.DeleteMessage(
		ctx,
		telegrambot.DeleteMessageParams{
			ChatID:    msg.Chat.ID,
			MessageID: msg.ID,
		},
	); err != nil {
		return err
	}
	if _, err := h.bot.SendVoice(
		ctx,
		telegrambot.SendVoiceParams{
			ChatID:      msg.Chat.ID,
			Voice:       &telegrambot.LocalSource{Filepath: audioPath},
			ReplyParams: &telegrambot.ReplyParameters{MessageID: &messageID},
		},
	); err != nil {
		return err
	}
	return nil
}
