package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"
)

type PostRepoImpl struct {
	postPreviews []domain.PostPreviewDto
	postList     domain.PostList
	postDto      domain.PostDto
	post         domain.Post
}

func (p PostRepoImpl) Create(post domain.Post, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	var wg sync.WaitGroup
	wg.Add(3)
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
		cur, err := conn.PostCollection.Find(context.TODO(), bson.M{"title": post.Title})

		if err != nil {
			panic(err)
		}
		if cur.Next(context.TODO()) {
			errorMessage = "title must be unique"
		}
	}()

	go func() {
		defer wg.Done()

		cur, err := conn.TagCollection.Find(context.TODO(), bson.M{"value": post.Tag})

		if err != nil {
			panic(err)
		}

		if !cur.Next(context.TODO()) {
			errorMessage = "you must include a valid tag"
		}
	}()

	wg.Wait()

	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}

	_, err := conn.PostCollection.InsertOne(context.TODO(), &post)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	go func() {
		err := SendKafkaMessage(&post, 201)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	go func() {
		err := TagRepoImpl{}.UpdateTag(post.Tag, post.Id)

		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (p PostRepoImpl) UpdateByTitle(post domain.PostUpdateDto, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	var wg sync.WaitGroup
	wg.Add(3)
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
		cur, err := conn.PostCollection.Find(context.TODO(), bson.M{"title": post.NewTitle})

		if err != nil {
			panic(err)
		}
		if cur.Next(context.TODO()) {
			errorMessage = "title must be unique"
		}
	}()

	go func() {
		defer wg.Done()

		cur, err := conn.TagCollection.Find(context.TODO(), bson.M{"value": post.Tag})

		if err != nil {
			panic(err)
		}

		if !cur.Next(context.TODO()) {
			errorMessage = "you must include a valid tag"
		}
	}()

	wg.Wait()

	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}

	err := conn.PostCollection.FindOneAndUpdate(context.TODO(), bson.D{{"title", post.Title}},
		bson.M{"$set": bson.M{
			"title":       post.NewTitle,
			"content":     post.Content,
			"mainImage":   post.MainImage,
			"storyImages": post.StoryImages,
			"tag":         post.Tag,
			"updated":     true,
			"updatedAt":   time.Now(),
		},
		}).Decode(&p.post)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	updatedPost := new(domain.Post)

	err = conn.PostCollection.FindOne(context.TODO(), bson.D{{"_id", p.post.Id}}).Decode(updatedPost)

	if err != nil {
		return err
	}

	go func() {
		err := SendKafkaMessage(updatedPost, 200)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func (p PostRepoImpl) UpdateVisibility(post domain.PostUpdateVisibilityDto, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	adminSearch := new(domain.Admin)
	err := conn.AdminCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(adminSearch)

	if err != nil {
		return err
	}

	err = conn.PostCollection.FindOneAndUpdate(context.TODO(), bson.D{{"title", post.Title}},
		bson.M{"$set": bson.M{
			"visible": post.Visible,
			"updated":     true,
			"updatedAt":   time.Now(),
		},
		}).Decode(&p.post)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	updatedPost := new(domain.Post)

	err = conn.PostCollection.FindOne(context.TODO(), bson.D{{"_id", p.post.Id}}).Decode(updatedPost)

	if err != nil {
		return err
	}

	go func() {
		err := SendKafkaMessage(updatedPost, 200)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func NewPostRepoImpl() PostRepoImpl {
	var postRepoImpl PostRepoImpl

	return postRepoImpl
}
