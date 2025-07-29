package logging

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger() *slog.Logger {
	options := slog.HandlerOptions{Level: getLogLevel()}
	handler := slog.NewJSONHandler(os.Stdout, &options)

	return slog.New(handler)
}

func getLogLevel() slog.Level {
	level := os.Getenv("LDMA_LOG_LEVEL")

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
