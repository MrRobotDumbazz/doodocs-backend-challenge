package service

import "doodocsbackendchallenge/internal/repository"

type Service struct {
	Upload
}

func NewServices(repositories *repository.Repository) *Service {
	return &Service{
		Upload: newUploadService(),
	}
}
