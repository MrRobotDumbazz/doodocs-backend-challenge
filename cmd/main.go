package main

import (
	"doodocsbackendchallenge/internal/config"
	"doodocsbackendchallenge/internal/logger"

	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger()
	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")
}
