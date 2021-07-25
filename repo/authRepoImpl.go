package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
}

func (a AuthRepoImpl) Login(username string, password string, ip string, ips []string) (*domain.Admin, string, error) {
	var login domain.Authentication
	var admin domain.Admin

	conn, err := database.ConnectToDB()
	defer func(conn *database.Connection, ctx context.Context) {
		err := conn.Disconnect(ctx)
		if err != nil {

		}
	}(conn, context.TODO())

	if err != nil {
		return nil, "", err
	}

	opts := options.FindOne()
	err = conn.AdminCollection.FindOne(context.TODO(), bson.D{{"username",
		username}}, opts).Decode(&admin)

	if err != nil {
		return nil, "", fmt.Errorf("error finding by username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))

	if err != nil {
		return nil, "", fmt.Errorf("error comparing password")
	}

	token, err := login.GenerateJWT(admin)

	if err != nil {
		return nil, "", fmt.Errorf("error generating token")
	}

	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"lastLoginIp", ip}, {"lastLoginIps", ips}}}}

	_, err = conn.AdminCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return nil, "", err
	}

	return &admin, token, nil
}

func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}
