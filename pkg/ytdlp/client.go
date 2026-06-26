package ytdlp

import (
	"log"
	"os"
	"os/exec"
)

type Client struct {
	BinPath string
}

func NewClient(binPath string) *Client {
	return &Client{
		BinPath: binPath,
	}
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
	filePath := string(out)
	return &filePath, nil
}

func (c *Client) DownloadVideo(
	outputPath string,
	videoURL string,
) (*os.File, error) {
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
	file, err := os.Open(*path)
	if err != nil {
		log.Println("Ошибка открытия файла: ", err)
		return nil, err
	}
	return file, nil
}
