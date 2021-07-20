package repo

import (
	"com.aharakitchen/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRepo interface {
	Create(post domain.Post, username string) error
	UpdateByTitle(post domain.PostUpdateDto, username string) error
	UpdateVisibility(post domain.PostUpdateVisibilityDto, username string) error
	FindAllPosts(page string, newPosts bool) (*domain.PostList, error)
	FindPostById(id primitive.ObjectID) (*domain.PostDto, error)
}
