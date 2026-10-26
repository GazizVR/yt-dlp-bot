package ytdlp

import (
	"strings"
)

type Client struct {
	bin       string
	cookies   string
	jsRuntime string
	browser   string
}

func NewClient(
	binPath string,
	cookiesPath string,
	jsRuntime string,
	browser string,
) *Client {
	return &Client{
		bin:       binPath,
		cookies:   cookiesPath,
		jsRuntime: jsRuntime,
		browser:   browser,
	}
}

func (c *Client) DownloadVideo(
	outputPath string,
	mediaURL string,
) (*string, error) {
	path, err := c.runBin(
		"-P", outputPath,
		"--remux-video", "mp4",
		mediaURL,
	)
	if err != nil {
		return nil, err
	}
	newPath := strings.ReplaceAll(*path, "\n", "")
	return &newPath, nil
}

func (c *Client) DownloadAudio(
	outputPath string,
	mediaURL string,
) (*string, error) {
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
	return &newPath, nil
}
