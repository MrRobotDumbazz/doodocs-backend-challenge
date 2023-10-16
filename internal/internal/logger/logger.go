package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

func SetupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
