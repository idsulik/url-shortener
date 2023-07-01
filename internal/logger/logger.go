package logger

import (
	"golang.org/x/exp/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New(env string) *Logger {
	switch env {
	case "dev":
		return &Logger{
			Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		}
	case "prod":
		return &Logger{
			Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})),
		}
	}

	return &Logger{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}
}
