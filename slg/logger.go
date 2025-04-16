package slg

import (
	"log/slog"
	"os"
	"strings"
)

var Logger = initLogger()

func initLogger() *slog.Logger {
	level := slog.LevelInfo

	if lvlStr, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch strings.ToLower(lvlStr) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	logger := slog.New(handler)
	logger.Info("Logger initialized", "level", level.String())
	return logger
}
