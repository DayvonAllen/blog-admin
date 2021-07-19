package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tag struct {
	Id             primitive.ObjectID `bson:"_id" json:"-"`
	Value          string `bson:"value" json:"value"`
	AssociatedPosts []primitive.ObjectID `bson:"associatedPosts"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"-"`
}

type TagDto struct {
	Id             primitive.ObjectID `bson:"_id" json:"-"`
	Value          string `bson:"value" json:"value"`
	AssociatedPosts []primitive.ObjectID `bson:"associatedPosts" json:"-"`
}
