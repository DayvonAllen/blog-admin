package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
	"time"
)

type PostRepoImpl struct {
	postPreviews []domain.PostPreviewDto
	postList     domain.PostList
	postDto      domain.PostDto
	post         domain.Post
}

func (p PostRepoImpl) FindAllPosts(page string, newPosts bool) (*domain.PostList, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newPosts {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	query := bson.M{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		cur, err := conn.PostCollection.Find(context.TODO(), query, &findOptions)

		if err != nil {
			panic(err)
		}

		if err = cur.All(context.TODO(), &p.postPreviews); err != nil {
			log.Fatal(err)
		}
		return
	}()

	go func() {
		defer wg.Done()
		count, err:= conn.PostCollection.CountDocuments(context.TODO(),query)

		if err != nil {
			panic(err)
		}

		p.postList.NumberOfPosts = count

		if p.postList.NumberOfPosts < 10 {
			p.postList.NumberOfPages = 1
		} else {
			p.postList.NumberOfPages = int(count / 10) + 1
		}
	}()

	wg.Wait()

	p.postList.Posts = p.postPreviews
	p.postList.CurrentPage = 1

	return &p.postList, nil
}

func (p PostRepoImpl) FindPostById(id primitive.ObjectID) (*domain.PostDto, error) {
	conn := database.MongoConn
	defer func(conn *database.Connection, ctx context.Context) {
		err := conn.Disconnect(ctx)
		if err != nil {

		}
	}(conn, context.TODO())

	query := bson.D{{"_id", id}}

	err := conn.PostCollection.FindOne(context.TODO(), query).Decode(&p.postDto)

	if err != nil {
		return nil, err
	}

	return &p.postDto, nil
}

func (p PostRepoImpl) Create(post domain.Post, username string) error {
	conn := database.MongoConn

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

	if len(post.Content) > 160 {
		post.Preview = string([]rune(post.Content)[:161]) + "..."
	} else {
		post.Preview = post.Content
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
	conn := database.MongoConn

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
	conn := database.MongoConn

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
