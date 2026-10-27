package config

import (
	"errors"
	"os"
)

type Config struct {
	BotToken      string
	BotAPIBaseURL string
	BotFilesDir   string
	YtDlpBin      string
	Cookies       string
	JsRuntime     string
	Browser       string
	FfmpegBin     string
}

func Load() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, errors.New("bot token must not be empty")
	}
	jsRuntime := os.Getenv("JS_RUNTIME")
	if jsRuntime == "" {
		return nil, errors.New("js runtime must not be empty")
	}
	ytDlpBin := os.Getenv("YT_DLP_BIN")
	if ytDlpBin == "" {
		ytDlpBin = "yt-dlp"
	}
	ffmpegBin := os.Getenv("FFMPEG_BIN")
	if ffmpegBin == "" {
		ffmpegBin = "ffmpeg"
	}
	return &Config{
		BotToken:      botToken,
		BotAPIBaseURL: os.Getenv("BOT_API_BASE_URL"),
		BotFilesDir:   os.Getenv("BOT_FILES_DIR"),
		YtDlpBin:      ytDlpBin,
		FfmpegBin:     ffmpegBin,
		Cookies:       os.Getenv("COOKIES"),
		JsRuntime:     jsRuntime,
		Browser:       os.Getenv("BROWSER"),
	}, nil
}
