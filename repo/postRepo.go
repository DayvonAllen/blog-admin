package repo

import (
	"com.aharakitchen/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRepo interface {
	FindAllPosts(page string) (*domain.PostList, error)
	Create(post domain.Post, username string) error
	UpdateByTitle(post domain.PostUpdateDto, username string) error
	FeaturedPosts() (*domain.PostList, error)
	FindPostById(id primitive.ObjectID) (*domain.PostDto, error)
}
