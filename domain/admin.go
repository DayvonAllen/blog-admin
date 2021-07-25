package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Admin struct {
	Id                          primitive.ObjectID   `bson:"_id" json:"-"`
	Username                    string               `bson:"username" json:"username"`
	Email                       string               `bson:"email" json:"-"`
	Password                    string               `bson:"password" json:"password"`
	LastLoginIp					string				 `bson:"lastLoginIp" json:"-"`
	LastLoginIps				[]string			 `bson:"lastLoginIps" json:"-"`
	CreatedAt                   time.Time            `bson:"createdAt" json:"-"`
	UpdatedAt                   time.Time            `bson:"updatedAt" json:"-"`
}
