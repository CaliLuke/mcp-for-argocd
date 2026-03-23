package logging

import (
	"log/slog"
	"os"
	"strings"
)

var Logger = newLogger()

func newLogger() *slog.Logger {
	level := new(slog.LevelVar)
	level.Set(parseLevel(os.Getenv("LOG_LEVEL")))

	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
