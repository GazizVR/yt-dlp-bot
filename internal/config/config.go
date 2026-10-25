package config

import "os"

type Config struct {
	BotToken      string
	BotAPIBaseURL string
	YtDlpBinPath  string
	CookiesPath   string
	JsRuntime     string
	Browser       string
}

func Load() *Config {
	return &Config{
		BotToken:      os.Getenv("BOT_TOKEN"),
		BotAPIBaseURL: os.Getenv("BOT_API_BASE_URL"),
		YtDlpBinPath:  os.Getenv("YT_DLP_PATH"),
		CookiesPath:   os.Getenv("COOKIES_PATH"),
		JsRuntime:     os.Getenv("JS_RUNTIME"),
		Browser:       os.Getenv("BROWSER"),
	}
}
