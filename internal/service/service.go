package service

import "doodocsbackendchallenge/internal/repository"

type Service struct {
	UploadFile
}

func NewServices(repositories *repository.Repository) *Service {
	return &Service{
		UploadFile: newUploadService(),
	}
}
