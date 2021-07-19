package repo

import (
	"com.aharakitchen/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagRepo interface {
	FindAllPostsByCategory(category, page string) (*domain.PostList, error)
	FindAllTags() (*[]domain.TagDto, error)
	Create(tag domain.Tag, username string) error
	UpdateTag(tagValue string, postId primitive.ObjectID) error
}
