package database

import (
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

var MongoConn *Connection

func ConnectToDB() {
	//uri := os.Getenv("DB_URI")

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

	client, err := mongo.Connect(ctx, dbOptions.ApplyURI("mongodb://admin-mongo-srv:27017/blog"))
	if err != nil {
		panic(err)
	}

	// create database
	db := client.Database("blog")

	// create collection
	postsCollection := db.Collection("posts")
	tagsCollection := db.Collection("tags")
	blackListCollection := db.Collection("blacklist")
	adminCollection := db.Collection("admin")

	MongoConn = &Connection{client, postsCollection, tagsCollection,blackListCollection, adminCollection,db}
}
