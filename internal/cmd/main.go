package main

import (
	"context"
	"doodocsbackendchallenge/internal/config"
	"doodocsbackendchallenge/internal/delivery"
	"doodocsbackendchallenge/internal/logger"
	"doodocsbackendchallenge/internal/server"
	"doodocsbackendchallenge/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger()
	services := service.NewServices(cfg)
	handlers := delivery.NewHandler(services, log)
	log.Debug("logger debug mode enabled")
	log.Info("Starting server", slog.String("address", cfg.Address))
	server := new(server.Server)
	go func() {
		if err := server.Start(cfg, handlers.Handlers()); err != nil {
			log.Error(err.Error())
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	log.Info("server started")
	<-quit
	log.Info("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", logger.Err(err))
		return
	}
	log.Info("server stopped")
}
