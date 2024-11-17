package service

import (
	"log/slog"

	"doodocsbackendchallenge/internal/config"
)

type Service struct {
	File
}

func NewServices(log *slog.Logger, cfg config.Config) *Service {
	return &Service{
		File: newFileService(log, cfg),
	}
}
