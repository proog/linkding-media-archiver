package logging

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger(logLevel string) *slog.Logger {
	options := slog.HandlerOptions{Level: getLogLevel(logLevel)}
	handler := slog.NewJSONHandler(os.Stdout, &options)

	return slog.New(handler)
}

func getLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
