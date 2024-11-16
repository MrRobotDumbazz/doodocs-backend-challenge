package service

import "log/slog"

type Service struct {
	File
}

func NewServices(log *slog.Logger) *Service {
	return &Service{
		File: newFileService(log),
	}
}
