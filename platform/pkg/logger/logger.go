package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Config interface {
	Level() string
	AsJSON() bool
}

const (
	LoggerLevelInfo    = "INFO"
	LoggerLevelError   = "ERROR"
	LoggerLevelWarning = "WARNING"
)

func Init(cfg Config) {
	slog.SetDefault(slog.New(newHandler(parseLevel(cfg.Level()), cfg.AsJSON())))
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
	case LoggerLevelInfo:
		return slog.LevelInfo
	case LoggerLevelError:
		return slog.LevelError
	case LoggerLevelWarning:
		return slog.LevelWarn
	default:
		return slog.LevelDebug
	}
}
