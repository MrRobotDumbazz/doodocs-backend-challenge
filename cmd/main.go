package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"doodocsbackendchallenge/internal/config"
	"doodocsbackendchallenge/internal/delivery"
	"doodocsbackendchallenge/internal/logger"
	"doodocsbackendchallenge/internal/server"
	"doodocsbackendchallenge/internal/service"
)

func main() {
	const op = "main:"
	log := logger.SetupLogger()
	log = log.With(
		slog.String("op", op),
	)
	cfg, err := config.MustLoad()
	if err != nil {
		log.Error(err.Error())
		return
	}
	services := service.NewServices(log, cfg)
	handlers := delivery.NewHandler(services, log)
	log.Debug("logger debug mode enabled")
	server := new(server.Server)
	log.Info("Starting server", slog.String("address", ":8080"))
	go func() {
		if err := server.Start(":8080", handlers.Handlers()); err != nil {
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	log.Info("server started", op)
	<-quit
	log.Info("stopping server", op)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", err.Error())
		return
	}
	log.Info("server stopped", op)
}
