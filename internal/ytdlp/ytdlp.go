package ytdlp

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func NewYtdlp(downloadDir string) *Ytdlp {
	return &Ytdlp{DownloadDir: downloadDir}
}

func (ytdlp *Ytdlp) DownloadMedia(url string) ([]string, error) {
	logger := slog.With("url", url)

	tempdir, err := os.MkdirTemp(ytdlp.DownloadDir, "media")

	if err != nil {
		return nil, err
	}

	args := []string{
		"--no-simulate",
		"--restrict-filenames",
		"--print",
		"after_move:filepath", // https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#outtmpl-postprocess-note
		url,
	}

	if ytdlp.MaxHeight > 0 {
		// Prefer single-file best up to MaxHeight to keep behavior similar to default.
		// Fallback selects bestvideo+bestaudio muxed up to MaxHeight when needed, and final fallback keeps the cap.
		format := fmt.Sprintf("b[height<=%d]/bv*[height<=%d]+ba/b[height<=%d]", ytdlp.MaxHeight, ytdlp.MaxHeight, ytdlp.MaxHeight)
		args = append(args, "--format", format)
	}

	args = append(args, url)
	cmd := exec.Command("yt-dlp", args...)
	cmd.Dir = tempdir

	logger.Debug("Downloading media", "command", cmd.String())
	output, err := cmd.Output()

	if err != nil {
		exitErr, ok := err.(*exec.ExitError)

		var stderr string
		if ok {
			stderr = string(exitErr.Stderr)
		}

		logger.Error("yt-dlp error", "stderr", stderr)
		return nil, err
	}

	paths := make([]string, 0, 10)
	for line := range strings.Lines(string(output)) {
		paths = append(paths, strings.TrimSpace(line))
	}

	logger.Debug("Downloaded media", "paths", paths)

	return paths, nil
}
