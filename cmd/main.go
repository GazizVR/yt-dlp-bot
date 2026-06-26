package main

import (
	"bot/internal/config"
	"bot/internal/service"
	"bot/pkg/telegram"
	"bot/pkg/ytdlp"
	"log"
)

func main() {
	cfg := config.Load()
	telegram := telegram.NewClient(cfg.TelegramToken, cfg.TelegramApiBaseURL)
	ytdlp := ytdlp.NewClient(cfg.YtDlpBinPath)
	service := service.NewService(telegram, ytdlp)
	log.Println("Сервер запускается")
	if err := service.Run(); err != nil {
		log.Fatalln("Ошибка запуска обработчика ошибок")
	}
}
