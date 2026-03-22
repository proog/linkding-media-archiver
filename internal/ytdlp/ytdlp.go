package ytdlp

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

func NewYtdlp(downloadDir string, format string) *Ytdlp {
	return &Ytdlp{DownloadDir: downloadDir, Format: format}
}

func (ytdlp *Ytdlp) DownloadMedia(url string) (*DownloadResult, error) {
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

	var jsonDump jsonDump
	err = json.Unmarshal(output, &jsonDump)

	if err != nil {
		return nil, err
	}

	result := newDownloadResult(&jsonDump)
	logger.Debug("Downloaded media", "result", result)

	if len(result.Paths) == 0 {
		return nil, fmt.Errorf("no paths in download result: %+v", result)
	}

	return &result, nil
}

func (ytdlp *Ytdlp) cmd(url string) *exec.Cmd {
	args := []string{
		"--no-simulate",
		"--restrict-filenames",
		"--dump-single-json",
	}

	// https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#format-selection
	if len(ytdlp.Format) > 0 {
		args = append(args, "--format", ytdlp.Format)
	}

	args = append(args, url)
	cmd := exec.Command("yt-dlp", args...)

	return cmd
}

func newDownloadResult(jsonDump *jsonDump) DownloadResult {
	paths := make([]string, 0, len(jsonDump.RequestedDownloads)+len(jsonDump.Entries))

	for _, download := range jsonDump.RequestedDownloads {
		paths = append(paths, download.FilePath)
	}

	for _, entry := range jsonDump.Entries {
		for _, download := range entry.RequestedDownloads {
			paths = append(paths, download.FilePath)
		}
	}

	return DownloadResult{
		Title:       jsonDump.Title,
		Description: jsonDump.Description,
		Tags:        jsonDump.Tags,
		Paths:       paths,
	}
}
