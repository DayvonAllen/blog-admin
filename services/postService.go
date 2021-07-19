package services

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	FindAllPosts(page string) (*domain.PostList, error)
	FeaturedPosts() (*domain.PostList, error)
	Create(post domain.Post, username string) error
	UpdateByTitle(post domain.PostUpdateDto, username string) error
	UpdateVisibility(post domain.PostUpdateVisibilityDto, username string) error
	FindPostById(id primitive.ObjectID) (*domain.PostDto, error)
}

type DefaultPostService struct {
	repo repo.PostRepo
}

func (s DefaultPostService) FindAllPosts(page string) (*domain.PostList, error) {
	postList, err := s.repo.FindAllPosts(page)
	if err != nil {
		return nil, err
	}
	return postList, nil
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

func (s DefaultPostService) FeaturedPosts() (*domain.PostList, error) {
	postList, err := s.repo.FeaturedPosts()
	if err != nil {
		return nil, err
	}
	return postList, nil
}

func (s DefaultPostService) FindPostById(id primitive.ObjectID) (*domain.PostDto, error) {
	post, err := s.repo.FindPostById(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func NewPostService(repository repo.PostRepo) DefaultPostService {
	return DefaultPostService{repository}
}