package bot

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/gazizvr/yt-dlp-bot/internal/bot/view"
	"github.com/gazizvr/yt-dlp-bot/internal/domain"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type MediaService interface {
	Create(ctx context.Context, url string) (int64, error)
	Download(ctx context.Context, id int64, action domain.Action) (filePath string, isVideo bool, err error)
	Release(ctx context.Context, id int64) error
}

type Handler struct {
	bot      *telego.Bot
	src      MediaService
	filesDir string
}

func NewHandler(
	bot *telego.Bot,
	src MediaService,
	filesDir string,
) *Handler {
	return &Handler{
		bot:      bot,
		src:      src,
		filesDir: filesDir,
	}
}

func (h *Handler) OnStart(
	ctx context.Context,
	upd telego.Update,
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
	upd telego.Update,
) error {
	var url string
	if upd.Message.LinkPreviewOptions != nil {
		url = upd.Message.LinkPreviewOptions.URL
	}
	if url == "" {
		url = upd.Message.Text
	}
	reqID, err := h.src.Create(ctx, url)
	if err != nil {
		h.bot.SendMessage(ctx, view.ErrorMessageParams(upd.Message.Chat.ID, upd.Message.MessageID))
		return err
	}
	strReqID := strconv.Itoa(int(reqID))
	if _, err := h.bot.SendMessage(ctx, view.ChooseDownloadMessageParams(upd.Message.Chat.ID, upd.Message.MessageID, strReqID, separator)); err != nil {
		h.bot.SendMessage(ctx, view.ErrorMessageParams(upd.Message.Chat.ID, upd.Message.MessageID))
		return err
	}
	return nil
}

func (h *Handler) OnCallbackQuery(
	ctx context.Context,
	upd telego.Update,
) error {
	if upd.CallbackQuery == nil || upd.CallbackQuery.Message == nil {
		return ErrMessageNotFound
	}
	chatID := upd.CallbackQuery.Message.GetChat().ID
	msgID := upd.CallbackQuery.Message.GetMessageID()

	msg, err := h.bot.EditMessageText(ctx, view.WaitMessageParams(msgID, chatID))
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(chatID, msgID))
		return err
	}
	data := upd.CallbackQuery.Data
	if data == "" {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
		return ErrEmptyCallbackData
	}
	rawID, rawAction, ok := strings.Cut(data, separator)
	if !ok || rawID == "" || rawAction == "" {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
		return ErrInvalidCallbackData
	}
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
		return err
	}
	action := domain.Action(rawAction)
	filepath, isVideo, err := h.src.Download(ctx, id, action)
	if err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
		return err
	}

	var inputFile telego.InputFile
	if h.filesDir != "" {
		fileURL := filepath
		if !strings.HasPrefix(fileURL, "file://") {
			fileURL = "file://" + fileURL
		}
		inputFile = telego.InputFile{URL: fileURL}
	} else {
		file, err := os.Open(filepath)
		if err != nil {
			h.bot.EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
			return err
		}
		defer file.Close()
		inputFile = tu.File(file)
	}

	var media telego.InputMedia
	if isVideo {
		media = &telego.InputMediaVideo{
			Type:              telego.MediaTypeVideo,
			Media:             inputFile,
			SupportsStreaming: true,
		}
	} else {
		media = &telego.InputMediaAudio{
			Type:  telego.MediaTypeAudio,
			Media: inputFile,
		}
	}
	params := &telego.EditMessageMediaParams{
		ChatID:    tu.ID(msg.Chat.ID),
		MessageID: msg.MessageID,
		Media:     media,
	}
	if _, err := h.bot.EditMessageMedia(ctx, params); err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.MessageID, data))
		return err
	}
	if err := h.src.Release(ctx, id); err != nil {
		h.bot.EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.MessageID, data))
		return err
	}
	return nil
}
