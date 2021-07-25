package database

import (
	"com.aharakitchen/app/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Connection struct {
	*mongo.Client
	PostCollection     *mongo.Collection
	TagCollection    *mongo.Collection
	BlackListCollection    *mongo.Collection
	AdminCollection *mongo.Collection
	*mongo.Database
}

func ConnectToDB() (*Connection, error) {
	p := config.Config("DB_PORT")
	n := config.Config("DB_NAME")
	h := config.Config("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//max := uint64(1000)
	//min := uint64(1)
	//idleTime := time.Second * 10
	dbOptions := options.ClientOptions{
		//MaxPoolSize: &max,
		//MinPoolSize: &min,
		//MaxConnIdleTime: &idleTime,
	}

	client, err := mongo.Connect(ctx, dbOptions.ApplyURI(n+h+p))
	if err != nil {
		return nil, err
	}

	// create database
	db := client.Database("blog-admin-service")

	// create collection
	postsCollection := db.Collection("posts")
	tagsCollection := db.Collection("tags")
	blackListCollection := db.Collection("blacklist")
	adminCollection := db.Collection("admin")

	dbConnection := &Connection{client, postsCollection, tagsCollection,blackListCollection, adminCollection,db}
	return dbConnection, nil
}
