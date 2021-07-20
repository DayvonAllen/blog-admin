package repo

import (
	"com.aharakitchen/app/domain"
)

type PostRepo interface {
	Create(post domain.Post, username string) error
	UpdateByTitle(post domain.PostUpdateDto, username string) error
	UpdateVisibility(post domain.PostUpdateVisibilityDto, username string) error
}
