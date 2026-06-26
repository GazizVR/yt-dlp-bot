package config

import "os"

type Config struct {
	TelegramToken      string
	TelegramApiBaseURL string
	YtDlpBinPath       string
}

func Load() *Config {
	return &Config{
		TelegramToken:      os.Getenv("TOKEN"),
		TelegramApiBaseURL: os.Getenv("BASE_URL"),
		YtDlpBinPath:       os.Getenv("YT_DLP_PATH"),
	}
}
