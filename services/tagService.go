package services

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/repo"
)

type TagService interface {
	Create(tag domain.Tag, username string) error
}

type DefaultTagService struct {
	repo repo.TagRepo
}

func (s DefaultTagService) Create(tag domain.Tag, username string) error {
	err := s.repo.Create(tag, username)
	if err != nil {
		return err
	}
	return nil
}

func NewTagService(repository repo.TagRepo) DefaultTagService {
	return DefaultTagService{repository}
}