package service

import "doodocsbackendchallenge/internal/config"

type Service struct {
	File
}

func NewServices(cfg *config.Config) *Service {
	return &Service{
		File: newFileService(cfg),
	}
}
