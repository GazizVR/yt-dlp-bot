package ytdlp

import (
	"log"
	"os/exec"
)

type Client struct {
	BinPath string
}

func (c *Client) runBin(args ...string) (*string, error) {
	path, err := exec.LookPath(c.BinPath)
	if err != nil {
		log.Println("Бинарник yt-dlp не найден: ", err)
		return nil, err
	}
	cmd := exec.Command(path, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Printf(
			"Ошибка при работе программы: %s\nВывод: %s\n",
			err.Error(),
			string(out),
		)
		return nil, err
	}
	return nil, nil
}

func (c *Client) DownloadVideo(
	outputPath string,
	videoURL string,
) (*string, error) {
	var videoPath string
	path, err := c.runBin(
		"-o", outputPath,
		"--recode-video", "mp4",
		"--quiet",
		"--no-warnings",
		"--print", "after_move:filepath",
		videoURL,
	)
	if err != nil {
		return nil, err
	}
	videoPath = *path
	return &videoPath, nil
}
