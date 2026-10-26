package ytdlp

import (
	"bytes"
	"log"
	"os/exec"
)

func (c *Client) baseArgs() []string {
	args := []string{
		"-o", "%(id)s-audio.%(ext)s",
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

func (c *Client) runBin(args ...string) (*string, error) {
	path, err := exec.LookPath(c.bin)
	if err != nil {
		log.Println("Бинарник yt-dlp не найден: ", err)
		return nil, err
	}
	args = append(c.baseArgs(), args...)
	cmd := exec.Command(path, args...)
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	out, err := cmd.Output()
	if err != nil {
		log.Printf(
			"Ошибка при работе программы: %s, %s",
			err.Error(),
			stdErr.String(),
		)
		return nil, err
	}
	filePath := string(out)
	return &filePath, nil
}
