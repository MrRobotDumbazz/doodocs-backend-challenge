package logger

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	log := &slog.Logger{}
	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
