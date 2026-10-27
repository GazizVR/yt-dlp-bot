package ytdlp

import (
	"path/filepath"
	"strings"
)

type Client struct {
	bin       string
	ffmpegBin string
	jsRuntime string
	cookies   string
	browser   string
}

type Option func(*Client) error

func WithCookies(cookies string) Option {
	return func(c *Client) error {
		c.cookies = cookies
		return nil
	}
}

func WithBrowser(browser string) Option {
	return func(c *Client) error {
		c.browser = browser
		return nil
	}
}

func NewClient(
	bin string,
	ffmpegBin string,
	jsRuntime string,
	opts ...Option,
) (*Client, error) {
	client := &Client{
		bin:       bin,
		ffmpegBin: ffmpegBin,
		jsRuntime: jsRuntime,
	}
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func (c *Client) DownloadVideo(url, outputDir string) (string, error) {
	args := c.baseArgs()
	args = append(
		args,
		[]string{
			url,
			"-P", outputDir,
			"--remux-video", "mp4",
		}...,
	)
	path, err := runBin(c.bin, args...)
	if err != nil {
		return "", err
	}
	newPath := strings.ReplaceAll(path, "\n", "")
	return newPath, nil
}

func (c *Client) VideoToAudio(videoPath string) (string, error) {
	fileExt := filepath.Ext(videoPath)
	filename := strings.TrimSuffix(videoPath, fileExt)
	_, err := runBin(
		c.ffmpegBin,
		"-i", videoPath,
		"-y",
		"-vn",
		"-loglevel", "error",
		"-c:a", "libmp3lame",
		"-q:a", "2",
		filename,
	)
	if err != nil {
		return "", err
	}
	return filename, nil
}
