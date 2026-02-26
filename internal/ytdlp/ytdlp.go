package ytdlp

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func NewYtdlp(downloadDir string, format string) *Ytdlp {
	return &Ytdlp{DownloadDir: downloadDir, Format: format}
}

func (ytdlp *Ytdlp) DownloadMedia(url string) ([]DownloadResult, error) {
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

	results := make([]DownloadResult, 0, 10)
	for line := range strings.Lines(string(output)) {
		mediaPath := strings.TrimSpace(line)
		result, err := createDownloadResult(mediaPath)

		if err != nil {
			return nil, err
		}

		results = append(results, *result)
	}

	logger.Debug("Downloaded media", "results", results)

	return results, nil
}

func (ytdlp *Ytdlp) cmd(url string) *exec.Cmd {
	args := []string{
		"--no-simulate",
		"--restrict-filenames",
		"--write-info-json",
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

func createDownloadResult(mediaPath string) (*DownloadResult, error) {
	logger := slog.With("mediaPath", mediaPath)
	logger.Debug("Creating download result")

	infoJsonPath := strings.TrimSuffix(mediaPath, filepath.Ext(mediaPath)) + ".info.json"
	infoJson, err := os.ReadFile(infoJsonPath)

	if err != nil {
		return nil, err
	}

	var metadata Metadata
	if err := json.Unmarshal(infoJson, &metadata); err != nil {
		return nil, err
	}

	result := DownloadResult{Path: mediaPath, Metadata: metadata}
	logger.Debug("Created download result")

	return &result, nil
}
