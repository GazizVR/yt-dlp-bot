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
		cfg.BotToken,
		telegrambot.WithBaseURL(cfg.BotAPIBaseURL),
		telegrambot.WithFilesDir(cfg.BotFilesDir),
	)
	if err != nil {
		log.Fatal(err)
	}
	ytdlp, err := ytdlp.NewClient(
		cfg.YtDlpBin,
		cfg.FfmpegBin,
		cfg.JsRuntime,
		ytdlp.WithBrowser(cfg.Browser),
        ytdlp.WithCookies(cfg.Cookies),
	)
    if err != nil {
        log.Fatal(err)
    }
	service := service.NewService(bot, ytdlp,ytdlp)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
