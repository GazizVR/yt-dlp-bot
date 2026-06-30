package ytdlp

import (
	"log"
	"os"
	"os/exec"
	"strings"
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
	mediaURL string,
) (*os.File, error) {
	path, err := c.runBin(
		"-P", outputPath,
		"--js-runtime", "node",
		"--remote-components", "ejs:github",
		"--quiet",
		"--no-warnings",
		"--print", "after_move:filepath",
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
		"--js-runtime", "node",
		"--remote-components", "ejs:github",
		"--quiet",
		"--no-warnings",
		"--print", "after_move:filepath",
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
