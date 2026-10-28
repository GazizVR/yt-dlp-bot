package bot

import (
	"context"
	"strconv"
	"strings"

	telegrambot "github.com/gazizvr/go-telegram-bot"
	"github.com/gazizvr/yt-dlp-bot/internal/bot/view"
	"github.com/gazizvr/yt-dlp-bot/internal/domain"
)

type Downloader interface {
	DownloadVideo(url, outDir string) (string, error)
	DownloadAudio(url, outDir string) (string, error)
}

type RequestStore interface {
	Get(ctx context.Context, id int64) (string, error)
	Create(ctx context.Context, mediaURL string) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	bot      *telegrambot.Bot
	src      Downloader
	filesDir string
	store    RequestStore
}

func NewHandler(
	bot *telegrambot.Bot,
	fileDir string,
	src Downloader,
	store RequestStore,
) *Handler {
	return &Handler{
		bot:      bot,
		filesDir: fileDir,
		src:      src,
		store:    store,
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

const separator string = "-"

func (h *Handler) OnURL(
	ctx context.Context,
	upd telegrambot.Update,
) error {
	url := upd.Message.LinkPreviewOptions.URL
	if url == "" {
		url = upd.Message.Text
	}
	reqID, err := h.store.Create(ctx, url)
	if err != nil {
		h.bot.SendMessage(ctx, view.ErrorMessageParams(upd.Message.Chat.ID, upd.Message.ID))
		return err
	}
	strReqID := strconv.Itoa(int(reqID))
	if _, err := h.bot.SendMessage(ctx, view.DownloadOptionsParams(upd.Message.Chat.ID, upd.Message.ID, strReqID, separator)); err != nil {
		h.bot.SendMessage(ctx, view.ErrorMessageParams(upd.Message.Chat.ID, upd.Message.ID))
		return err
	}
	return nil
}

func (h *Handler) OnCallbackQuery(
	ctx context.Context,
	upd telegrambot.Update,
) error {
	var chatID, msgID int64
	switch v := upd.CallbackQuery.Message.(type) {
	case telegrambot.InaccessibleMessage:
		chatID = v.Chat.ID
		msgID = v.ID
	case telegrambot.Message:
		chatID = v.Chat.ID
		msgID = v.ID
	default:
		return ErrMessageNotFound
	}
	msg, err := h.bot.EditMessageText(ctx, view.WaitEditMessageParams(msgID, chatID))
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(chatID, msgID))
		return err
	}
	data := upd.CallbackQuery.Data
	if data == "" {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return ErrEmptyCallbackData
	}
	separatorInd := strings.Index(data, separator)
	if separatorInd == -1 {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return ErrInvalidCallbackData
	}
	rawID := data[:separatorInd]
	rawID = strings.TrimSpace(rawID)
	if rawID == "" {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return ErrInvalidCallbackData
	}
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return err
	}
	rawAction := data[separatorInd+1:]
	rawAction = strings.TrimSpace(rawAction)
	if rawAction == "" {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return ErrInvalidCallbackData
	}
	action := domain.Action(rawAction)
	mediaPath, err := h.store.Get(ctx, id)
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.ID, data))
		return err
	}
	var filepath string
	var inputMediaType telegrambot.MediaType
	switch action {
	case domain.SendVideo:
		filepath, err = h.src.DownloadVideo(mediaPath, h.filesDir)
		if err != nil {
			h.bot.EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.ID, data))
			return err
		}
		inputMediaType = telegrambot.MediaTypeVideo
	case domain.SendAudio:
		filepath, err = h.src.DownloadAudio(mediaPath, h.filesDir)
		if err != nil {
			h.bot.EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.ID, data))
			return err
		}
		inputMediaType = telegrambot.MediaTypeAudio
	}
	params := telegrambot.EditMessageMediaParams{
		ChatID:    msg.Chat.ID,
		MessageID: msg.ID,
		Media: telegrambot.InputMedia{
			Type:  inputMediaType,
			Media: &telegrambot.LocalSource{Filepath: filepath},
		},
	}
	if _, err := h.bot.EditMessageMedia(ctx, params); err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.ID, data))
		return err
	}
	if err := h.store.Delete(ctx, id); err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.ID, data))
		return err
	}
	return nil
}
