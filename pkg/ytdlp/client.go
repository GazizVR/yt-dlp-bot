package ytdlp

import (
	"log"
	"os"
	"strings"
)

type Client struct {
	BinPath     string
	CookiesPath string
	JsRuntime   string
	Browser     string
}

func NewClient(
	binPath string,
	cookiesPath string,
	jsRuntime string,
	browser string,
) *Client {
	return &Client{
		BinPath:     binPath,
		CookiesPath: cookiesPath,
		JsRuntime:   jsRuntime,
		Browser:     browser,
	}
}

func (c *Client) DownloadVideo(
	outputPath string,
	mediaURL string,
) (*os.File, error) {
	path, err := c.runBin(
		"-P", outputPath,
		"--remux-video", "mp4",
		mediaURL,
	)
	if err != nil {
		return nil, err
	}
	newPath := strings.ReplaceAll(*path, "\n", "")
	file, err := os.Open(newPath)
	if err != nil {
		log.Printf("Ошибка открытия %q: %v", newPath, err)
		return nil, err
	}
	return file, nil
}

func (c *Client) DownloadAudio(
	outputPath string,
	mediaURL string,
) (*os.File, error) {
	path, err := c.runBin(
		"-P", outputPath,
		"-x",
		"--audio-format", "mp3",
		mediaURL,
	)
	if err != nil {
		return nil, err
	}
	newPath := strings.ReplaceAll(*path, "\n", "")
	file, err := os.Open(newPath)
	if err != nil {
		log.Printf("Ошибка открытия %q: %v", newPath, err)
		return nil, err
	}
	return file, nil
}
