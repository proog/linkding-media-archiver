package ytdlp

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func NewYtdlp(downloadDir string, format string) *Ytdlp {
	return &Ytdlp{DownloadDir: downloadDir, Format: format}
}

func (ytdlp *Ytdlp) DownloadMedia(url string) ([]string, error) {
	logger := slog.With("url", url)

	tempdir, err := os.MkdirTemp(ytdlp.DownloadDir, "media")

	if err != nil {
		return nil, err
	}

	cmd := ytdlp.cmd(url)
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

func (ytdlp *Ytdlp) cmd(url string) *exec.Cmd {
	args := []string{
		"--no-simulate",
		"--restrict-filenames",
		"--print",
		"after_move:filepath", // https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#outtmpl-postprocess-note
	}

	// https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#format-selection
	if len(ytdlp.Format) > 0 {
		args = append(args, "--format", ytdlp.Format)
	}

	args = append(args, url)
	cmd := exec.Command("yt-dlp", args...)

	return cmd
}
