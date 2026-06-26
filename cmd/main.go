package main

import (
	"bot/internal/config"
	"bot/internal/service"
	"bot/pkg/telegram"
	"log"
)

func main() {
	cfg := config.Load()
	client := telegram.NewClient(cfg.TelegramToken, cfg.TelegramApiBaseURL)
	service := service.NewService(client)
	log.Println("Сервер запускается")
	if err := service.StartUpdateHandle(); err != nil {
		log.Fatalln("Ошибка запуска обработчика ошибок")
	}
}
