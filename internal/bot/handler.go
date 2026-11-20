package bot

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/gazizvr/yt-dlp-bot/internal/bot/view"
	"github.com/gazizvr/yt-dlp-bot/internal/domain"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type MediaService interface {
	Create(ctx context.Context, url string) (int64, error)
	Download(ctx context.Context, id int64, action domain.Action) (filePath string, isVideo bool, err error)
	Release(ctx context.Context, id int64) error
}

type Handler struct {
	src      MediaService
	filesDir string
}

func NewHandler(
	bot *telego.Bot,
	src MediaService,
	filesDir string,
) *Handler {
	return &Handler{
		src:      src,
		filesDir: filesDir,
	}
}

func (h *Handler) OnStart(
	ctx *th.Context,
	upd telego.Update,
) error {
	params := view.StartMessageParams(upd.Message.Chat.ID)
	if _, err := ctx.Bot().SendMessage(ctx, params); err != nil {
		return err
	}
	return nil
}

const separator string = "-"

func (h *Handler) OnDownload(
	ctx *th.Context,
	msg telego.Message,
) error {
	var url string
	if msg.LinkPreviewOptions != nil {
		url = msg.LinkPreviewOptions.URL
	}
	if url == "" {
		url = msg.Text
	}
	reqID, err := h.src.Create(ctx, url)
	if err != nil {
		ctx.Bot().SendMessage(ctx, view.ErrorMessageParams(msg.Chat.ID, msg.MessageID))
		return err
	}
	strReqID := strconv.Itoa(int(reqID))
	if _, err := ctx.Bot().SendMessage(ctx, view.ChooseDownloadMessageParams(msg.Chat.ID, msg.MessageID, strReqID, separator)); err != nil {
		ctx.Bot().SendMessage(ctx, view.ErrorMessageParams(msg.Chat.ID, msg.MessageID))
		return err
	}
	return nil
}

func (h *Handler) OnCallbackQuery(
	ctx *th.Context,
	query telego.CallbackQuery,
) error {
	if query.Message == nil {
		return ErrMessageNotFound
	}
	chatID := query.Message.GetChat().ID
	msgID := query.Message.GetMessageID()

	msg, err := ctx.Bot().EditMessageText(ctx, view.WaitMessageParams(msgID, chatID))
	if err != nil {
		ctx.Bot().EditMessageText(ctx, view.ErrorEditMessageParams(chatID, msgID))
		return err
	}
	data := query.Data
	if data == "" {
		ctx.Bot().EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
		return ErrEmptyCallbackData
	}
	rawID, rawAction, ok := strings.Cut(data, separator)
	if !ok || rawID == "" || rawAction == "" {
		ctx.Bot().EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
		return ErrInvalidCallbackData
	}
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		ctx.Bot().EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
		return err
	}
	action := domain.Action(rawAction)
	filepath, isVideo, err := h.src.Download(ctx, id, action)
	if err != nil {
		ctx.Bot().EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
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
			ctx.Bot().EditMessageText(ctx, view.ErrorEditMessageParams(msg.Chat.ID, msg.MessageID))
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
	if _, err := ctx.Bot().EditMessageMedia(ctx, params); err != nil {
		ctx.Bot().EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.MessageID, data))
		return err
	}
	if err := h.src.Release(ctx, id); err != nil {
		ctx.Bot().EditMessageText(ctx, view.ErrorEditMessageMarkupParams(msg.Chat.ID, msg.MessageID, data))
		return err
	}
	return nil
}
