package services

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/repo"
)

type PostService interface {
	Create(post domain.Post, username string) error
	UpdateByTitle(post domain.PostUpdateDto, username string) error
	UpdateVisibility(post domain.PostUpdateVisibilityDto, username string) error
}

type DefaultPostService struct {
	repo repo.PostRepo
}

func (s DefaultPostService) Create(post domain.Post, username string) error {
	err := s.repo.Create(post, username)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultPostService) UpdateByTitle(post domain.PostUpdateDto, username string) error {
	err := s.repo.UpdateByTitle(post, username)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultPostService) UpdateVisibility(post domain.PostUpdateVisibilityDto, username string) error {
	err := s.repo.UpdateVisibility(post, username)
	if err != nil {
		return err
	}
	return nil
}

func NewPostService(repository repo.PostRepo) DefaultPostService {
	return DefaultPostService{repository}
}