package ytdlp

import (
	"log/slog"
	"os/exec"
	"strings"
)

func NewYtdlp(downloadDir string) *Ytdlp {
	return &Ytdlp{DownloadDir: downloadDir}
}

func (ytdlp *Ytdlp) DownloadVideo(url string) (string, error) {
	logger := slog.With("url", url)

	args := []string{
		"--no-simulate",
		"--restrict-filenames",
		"--print",
		"after_move:filepath", // https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#outtmpl-postprocess-note
		url,
	}

	cmd := exec.Command("yt-dlp", args...)
	cmd.Dir = ytdlp.DownloadDir

	logger.Info("Downloading video", "command", cmd)
	output, err := cmd.Output()

	if err != nil {
		exitErr, ok := err.(*exec.ExitError)

		var stderr string
		if ok {
			stderr = string(exitErr.Stderr)
		}

		logger.Error("yt-dlp error", "stderr", stderr)
		return "", err
	}

	path := strings.TrimSpace(string(output))
	logger.Info("Downloaded video", "path", path)

	return path, nil
}
