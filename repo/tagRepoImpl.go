package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type TagRepoImpl struct {
	postPreviews   []domain.PostPreviewDto
	postList   domain.PostList
	tags []domain.TagDto
	tag domain.Tag
}

func (t TagRepoImpl) Create(tag domain.Tag, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	var wg sync.WaitGroup
	wg.Add(2)
	errorMessage := ""

	go func() {
		defer wg.Done()
		adminSearch := new(domain.Admin)
		err := conn.AdminCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(adminSearch)

		if err != nil {
			errorMessage = "unauthorized"
		}
	}()

	go func() {
		defer wg.Done()
		cur, err := conn.TagCollection.Find(context.TODO(), bson.M{"value": tag.Value})

		if err != nil {
			panic(err)
		}
		if cur.Next(context.TODO()) {
			errorMessage = "tag already exists"
		}
	}()

	wg.Wait()

	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}

	posts := make([]primitive.ObjectID, 0, 0)
	tag.Id = primitive.NewObjectID()
	tag.AssociatedPosts = posts

	_, err := conn.TagCollection.InsertOne(context.TODO(), &tag)

	if err != nil {
		return fmt.Errorf("errorMessage processing data")
	}

	go func() {
		err := SendAltKafkaMessage(&tag, 201)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func (t TagRepoImpl) UpdateTag(tagValue string, postId primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.TagCollection.FindOne(context.TODO(), bson.M{"value": tagValue}).Decode(&t.tag)

	if err != nil {
		return err
	}

	cur, err := conn.TagCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": t.tag.AssociatedPosts}})

	if err != nil {
		return err
	}

	if cur.Next(context.TODO()) {
		return fmt.Errorf("already associated")
	}

	_, err = conn.TagCollection.UpdateOne(context.TODO(), bson.D{{"value", tagValue}}, bson.M{"$push": bson.M{"associatedPosts": postId}})

	if err != nil {
		return fmt.Errorf("errorMessage processing data")
	}

	messageObj := new(domain.Tag)

	err = conn.TagCollection.FindOne(context.TODO(), bson.M{"value": tagValue}).Decode(&messageObj)

	if err != nil {
		return err
	}
	go func() {
		err := SendAltKafkaMessage(messageObj, 200)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func NewTagRepoImpl() TagRepoImpl {
	var tagRepoImpl TagRepoImpl

	return tagRepoImpl
}
