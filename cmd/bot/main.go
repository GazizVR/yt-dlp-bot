package main

import (
	"context"
	"log"

	"github.com/gazizvr/yt-dlp-bot/internal/bot"
	"github.com/gazizvr/yt-dlp-bot/internal/config"
	"github.com/gazizvr/yt-dlp-bot/internal/repository"
	"github.com/gazizvr/yt-dlp-bot/pkg/ytdlp"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ytdlp, err := ytdlp.NewClient(
		cfg.YtDlpBin,
		cfg.JsRuntime,
		ytdlp.WithBrowser(cfg.Browser),
		ytdlp.WithCookies(cfg.Cookies),
	)
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	requestStore := repository.NewRequestStore(rdb)
	bot, err := bot.NewBot(cfg.BotAPIBaseURL, cfg.BotToken, cfg.BotFilesDir, ytdlp, requestStore)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err := bot.ListenUpdates(ctx); err != nil {
		log.Fatal(err)
	}
}
