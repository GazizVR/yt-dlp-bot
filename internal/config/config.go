package config

import "os"

type Config struct {
	TelegramToken      string
	TelegramApiBaseURL string
	YtDlpBinPath       string
	CookiesPath        string
	JsRuntime          string
	Browser            string
}

func Load() *Config {
	return &Config{
		TelegramToken:      os.Getenv("TOKEN"),
		TelegramApiBaseURL: os.Getenv("BASE_URL"),
		YtDlpBinPath:       os.Getenv("YT_DLP_PATH"),
		CookiesPath:        os.Getenv("COOKIES_PATH"),
		JsRuntime:          os.Getenv("JS_RUNTIME"),
		Browser:            os.Getenv("BROWSER"),
	}
}
