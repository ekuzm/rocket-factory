package logger

import (
	"log/slog"
	"os"
	"strings"
)

func Init(level string, asJSON bool) {
	slog.SetDefault(slog.New(newHandler(parseLevel(level), asJSON)))
}

func newHandler(level slog.Level, asJSON bool) slog.Handler {
	options := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	if asJSON {
		return slog.NewJSONHandler(os.Stdout, options)
	}

	return slog.NewTextHandler(os.Stdout, options)
}

func parseLevel(level string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "INFO":
		return slog.LevelInfo
	case "ERROR":
		return slog.LevelError
	case "WARNING", "WARN":
		return slog.LevelWarn
	default:
		return slog.LevelDebug
	}
}
