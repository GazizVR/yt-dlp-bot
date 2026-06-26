package ytdlp

import (
	"log"
	"os/exec"
)

type Client struct {
	BinPath string
}

func (c *Client) runBin(args ...string) error {
	path, err := exec.LookPath(c.BinPath)
	if err != nil {
		log.Println("Бинарник yt-dlp не найден: ", err)
		return err
	}
	cmd := exec.Command(path, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Printf(
			"Ошибка при работе программы: %s\nВывод: %s\n",
			err.Error(),
			string(out),
		)
		return err
	}
	return nil
}

func (c *Client) DownloadVideo() (*string, error) {
	var videoPath string
	return &videoPath, nil
}
