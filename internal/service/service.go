package service

import "doodocsbackendchallenge/internal/config"

//go:generate mockgen -source=service.go -destination=mocks/mockservice.go

type Service struct {
	File
}

func NewServices(cfg *config.Config) *Service {
	return &Service{
		File: newFileService(cfg),
	}
}
