package main

import (
	"log"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/yt-dlp-bot/internal/config"
	"github.com/gazizvr/yt-dlp-bot/internal/service"
	"github.com/gazizvr/yt-dlp-bot/pkg/ytdlp"
)

func main() {
	cfg := config.Load()
	telegram := telegram.NewClient(cfg.TelegramToken, cfg.TelegramApiBaseURL)
	ytdlp := ytdlp.NewClient(
		cfg.YtDlpBinPath,
		cfg.CookiesPath,
		cfg.JsRuntime,
		cfg.Browser,
	)
	service := service.NewService(telegram, ytdlp)
	log.Println("Сервер запускается")
	if err := service.Run(); err != nil {
		log.Fatalln("Ошибка запуска обработчика ошибок")
	}
}
