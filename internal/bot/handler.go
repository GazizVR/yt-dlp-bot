package bot

import (
	"context"
	"strconv"
	"strings"

	telegrambot "github.com/gazizvr/go-telegram-bot"
	"github.com/gazizvr/yt-dlp-bot/internal/bot/view"
	"github.com/gazizvr/yt-dlp-bot/internal/domain"
)

type MediaService interface {
	Create(ctx context.Context, url string) (int64, error)
	Download(ctx context.Context, id int64, action domain.Action) (filePath string, isVideo bool, err error)
	Release(ctx context.Context, id int64) error
}

type Handler struct {
	bot *telegrambot.Bot
	src MediaService
}

func NewHandler(
	bot *telegrambot.Bot,
	src MediaService,
) *Handler {
	return &Handler{
		bot: bot,
		src: src,
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

func (h *Handler) OnDownload(
	ctx context.Context,
	upd telegrambot.Update,
) error {
	url := upd.Message.LinkPreviewOptions.URL
	if url == "" {
		url = upd.Message.Text
	}
	reqID, err := h.src.Create(ctx, url)
	if err != nil {
		h.bot.SendMessage(ctx, view.ErrorMessageParams(upd.Message.Chat.ID, upd.Message.ID))
		return err
	}
	strReqID := strconv.Itoa(int(reqID))
	if _, err := h.bot.SendMessage(ctx, view.ChooseDownloadMessageParams(upd.Message.Chat.ID, upd.Message.ID, strReqID, separator)); err != nil {
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
	msg, err := h.bot.EditMessageText(ctx, view.WaitMessageParams(msgID, chatID))
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(chatID, msgID))
		return err
	}
	data := upd.CallbackQuery.Data
	if data == "" {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return ErrEmptyCallbackData
	}
	rawID, rawAction, ok := strings.Cut(data, separator)
	if !ok || rawID == "" || rawAction == "" {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return ErrInvalidCallbackData
	}
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return err
	}
	action := domain.Action(rawAction)
	filepath, isVideo, err := h.src.Download(ctx, id, action)
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.ID))
		return ErrInvalidCallbackData
	}
	var inputMediaType telegrambot.MediaType
	if isVideo {
		inputMediaType = telegrambot.MediaTypeVideo
	} else {
		inputMediaType = telegrambot.MediaTypeVideo
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
	if err := h.src.Release(ctx, id); err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.ID, data))
		return err
	}
	return nil
}
