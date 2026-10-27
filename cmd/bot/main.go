package main

import (
	"log"

	telegrambot "github.com/gazizvr/go-telegram-bot"
	"github.com/gazizvr/yt-dlp-bot/internal/config"
	"github.com/gazizvr/yt-dlp-bot/internal/service"
	"github.com/gazizvr/yt-dlp-bot/pkg/ytdlp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	bot, err := telegrambot.NewBot(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.BotAPIBaseURL != "" {
		if err := telegrambot.WithBaseURL(cfg.BotAPIBaseURL)(bot); err != nil {
			log.Fatal(err)
		}
	}
	if cfg.BotFilesDir != "" {
		if err := telegrambot.WithFilesDir(cfg.BotFilesDir)(bot); err != nil {
			log.Fatal(err)
		}
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
	service := service.NewService(bot, ytdlp, ytdlp, cfg.BotFilesDir)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
