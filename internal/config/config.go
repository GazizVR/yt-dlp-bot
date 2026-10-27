package config

import "os"

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

func Load() *Config {
	ytDlpBin := os.Getenv("YT_DLP_BIN")
	if ytDlpBin == "" {
		ytDlpBin = "yt-dlp"
	}
    ffmpegBin := os.Getenv("FFMPEG_BIN")
    if ffmpegBin == "" {
        ffmpegBin = "ffmpeg"
    }
	return &Config{
		BotToken:      os.Getenv("BOT_TOKEN"),
		BotAPIBaseURL: os.Getenv("BOT_API_BASE_URL"),
		BotFilesDir:   os.Getenv("BOT_FILES_DIR"),
		YtDlpBin:      ytDlpBin,
        FfmpegBin:     ffmpegBin,
		Cookies:       os.Getenv("COOKIES"),
		JsRuntime:     os.Getenv("JS_RUNTIME"),
		Browser:       os.Getenv("BROWSER"),
	}
}
