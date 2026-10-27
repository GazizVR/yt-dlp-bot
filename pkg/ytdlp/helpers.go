package ytdlp

import (
	"bytes"
	"fmt"
	"os/exec"
)

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

func runBin(bin string, args ...string) (string, error) {
	path, err := exec.LookPath(bin)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(path, args...)
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("err: %w, stdErr: %s", err, stdErr.String())
	}
	return string(out), nil
}
