package service

import (
	"blogs/internal/model/entity"
	"blogs/internal/repository"
)

type aboutService struct {
	aboutRepo repository.AboutRepository
}

func NewAboutService(aboutRepo repository.AboutRepository) AboutService {
	return &aboutService{aboutRepo: aboutRepo}
}

func (s *aboutService) Get() (*entity.About, error) {
	return s.aboutRepo.Get()
}

func (s *aboutService) Update(about *entity.About) error {
	return s.aboutRepo.Upsert(about)
}
