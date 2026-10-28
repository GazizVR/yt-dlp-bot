package ytdlp

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	bin       string
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
	jsRuntime string,
	opts ...Option,
) (*Client, error) {
	client := &Client{
		bin:       bin,
		jsRuntime: jsRuntime,
	}
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func (c *Client) baseArgs() []string {
	args := []string{
		"-o", "m%(id)s.%(ext)s",
		"--remote-components", "ejs:github",
		"--quiet",
		"--no-warnings",
		"--print", "after_move:filepath",
	}
	if len(c.cookies) > 0 {
		args = append(args, "--cookies")
		args = append(args, c.cookies)
	}
	if len(c.jsRuntime) > 0 {
		args = append(args, "--js-runtime")
		args = append(args, c.jsRuntime)
	}
	if len(c.browser) > 0 {
		args = append(args, "--cookies-from-browser")
		args = append(args, c.browser)
	}
	return args
}

func (c *Client) runBin(args ...string) (string, error) {
	path, err := exec.LookPath(c.bin)
	if err != nil {
		return "", err
	}
	args = append(args, c.baseArgs()...)
	cmd := exec.Command(path, args...)
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("err: %w, stdErr: %s", err, stdErr.String())
	}
	return string(out), nil
}

func (c *Client) DownloadVideo(url, outDir string) (string, error) {
	path, err := c.runBin(
		url,
		"-P", outDir,
		"--remux-video", "mp4",
	)
	if err != nil {
		return "", err
	}
	newPath := strings.ReplaceAll(path, "\n", "")
	return newPath, nil
}

func (c *Client) DownloadAudio(url, outDir string) (string, error) {
	path, err := c.runBin(
		url,
		"-P", outDir,
		"-x",
		"-f", "ba[ext=m4a]/ba",
	)
	if err != nil {
		return "", err
	}
	newPath := strings.ReplaceAll(path, "\n", "")
	return newPath, nil
}
