package services

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/repo"
)

type TagService interface {
	Create(tag domain.Tag, username string) error
	FindAllPostsByCategory(category, page string) (*domain.PostList, error)
	FindAllTags() (*domain.TagList, error)
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

func (s DefaultTagService) FindAllPostsByCategory(category, page string) (*domain.PostList, error) {
	postList, err := s.repo.FindAllPostsByCategory(category, page)
	if err != nil {
		return nil, err
	}
	return postList, nil
}

func (s DefaultTagService) FindAllTags() (*domain.TagList, error) {
	tags, err := s.repo.FindAllTags()
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func NewTagService(repository repo.TagRepo) DefaultTagService {
	return DefaultTagService{repository}
}