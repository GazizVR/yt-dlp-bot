package config

import "os"

type Config struct {
	TelegramToken      string
	TelegramApiBaseURL string
}

func Load() *Config {
	return &Config{
		TelegramToken:      os.Getenv("TOKEN"),
		TelegramApiBaseURL: os.Getenv("BASE_URL"),
	}
}
