package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"os/exec"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func StartHandler(ctx *th.Context, update telego.Update) error {
	ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
		Text:   "Отправьте ссылку, а я вам видео!",
		ChatID: update.Message.Chat.ChatID(),
	})
	return nil
}

func TextHandler(ctx *th.Context, update telego.Update) error {
	msg, _ := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
		Text:   "🔄 Идет загрузка...",
		ChatID: update.Message.Chat.ChatID(),
	})
	mediaPath := filepath.Join("tmp", fmt.Sprint("media-", msg.Chat.ChatID()))
	prgm := exec.Command(
		"./yt-dlp",
		"--force-overwrites",
		"-o", mediaPath,
		"--print", "after_move:filepath",
		"--merge-output-format", "mp4",
		update.Message.Text,
	)
	dt, err := prgm.Output()
	if err != nil {
		ctx.Bot().EditMessageText(
			ctx, &telego.EditMessageTextParams{
				Text:      "⛔️ Ошибка загрузки",
				MessageID: msg.GetMessageID(),
				ChatID:    msg.Chat.ChatID(),
			},
		)
		return err
	} else {
		path := string(dt)
		path = strings.ReplaceAll(path, "\n", "")
		path = strings.TrimSpace(path)
		video, err := os.Open(path)
		if err != nil {
			return err
		}
		ctx.Bot().DeleteMessage(
			ctx,
			&telego.DeleteMessageParams{
				MessageID: msg.MessageID,
				ChatID:    msg.Chat.ChatID(),
			},
		)
		ctx.Bot().SendVideo(
			ctx,
			&telego.SendVideoParams{
				ChatID: msg.Chat.ChatID(),
				Video: telego.InputFile{
					File: video,
				},
			},
		)
	}
	return nil
}

func main() {
	botToken := os.Getenv("BOT_TOKEN")

	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Fatalln("Ошибка создания бота: ", err)
	}
	ctx := context.Background()
	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		log.Fatalln("Ошибка создания обработчика обновлений: ", err)
	}
	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatalln("Ошибка регистрация обработчка обновлений: ", err)
	}
	defer bh.Stop()
	bh.Handle(StartHandler, th.CommandEqual("start"))
	bh.Handle(TextHandler)
	bh.Start()
}
