package main

import (
	"log"

	telegrambot "github.com/gazizvr/go-telegram-bot"
	"github.com/gazizvr/yt-dlp-bot/internal/config"
	"github.com/gazizvr/yt-dlp-bot/internal/service"
	"github.com/gazizvr/yt-dlp-bot/pkg/ytdlp"
)

func main() {
	cfg := config.Load()
	bot, err := telegrambot.NewBot(
		cfg.TelegramToken,
		telegrambot.WithBaseURL(cfg.TelegramApiBaseURL),
	)
	if err != nil {
		log.Fatal(err)
	}
	ytdlp := ytdlp.NewClient(
		cfg.YtDlpBinPath,
		cfg.CookiesPath,
		cfg.JsRuntime,
		cfg.Browser,
	)
	service := service.NewService(bot, ytdlp)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
