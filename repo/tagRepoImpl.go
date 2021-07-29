package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
)

type TagRepoImpl struct {
	postPreviews   []domain.PostPreviewDto
	postList   domain.PostList
	tags []domain.TagDto
	tag domain.Tag
	tagList domain.TagList
}

func (t TagRepoImpl) FindAllTags() (*domain.TagList, error) {
	conn := database.MongoConn

	cur, err := conn.TagCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &t.tags); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	t.tagList.Tags = t.tags
	t.tagList.NumberOfCategories = len(t.tags)

	return &t.tagList, nil
}

func (t TagRepoImpl) FindAllPostsByCategory(category, page string) (*domain.PostList, error) {
	conn := database.MongoConn

	err := conn.TagCollection.FindOne(context.TODO(), bson.D{{"value", category}}).Decode(&t.tag)

	if err != nil {
		return nil, err
	}

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	query := bson.M{"_id": bson.M{"$in": t.tag.AssociatedPosts}}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		cur, err := conn.PostCollection.Find(context.TODO(), query, &findOptions)

		if err != nil {
			panic(err)
		}

		if err = cur.All(context.TODO(), &t.postPreviews); err != nil {
			log.Fatal(err)
		}

	}()

	go func() {
		defer wg.Done()
		count, err:= conn.PostCollection.CountDocuments(context.TODO(),query)

		if err != nil {
			panic(err)
		}

		t.postList.NumberOfPosts = count

		if t.postList.NumberOfPosts < 10 {
			t.postList.NumberOfPages = 1
		} else {
			t.postList.NumberOfPages = int(count / 10) + 1
		}
	}()

	wg.Wait()

	t.postList.Posts = t.postPreviews
	t.postList.CurrentPage = 1

	return &t.postList, nil
}

func (t TagRepoImpl) Create(tag domain.Tag, username string) error {
	conn := database.MongoConn

	err := conn.TagCollection.FindOne(context.TODO(), bson.M{"value": tag.Value}).Decode(&t.tag)

	if err != nil {
		if err == mongo.ErrNoDocuments {
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

			_, err = conn.TagCollection.InsertOne(context.TODO(), &tag)

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
		}
		return err
	}

	return fmt.Errorf("tag already exists")
}

func (t TagRepoImpl) UpdateTag(tagValue string, postId primitive.ObjectID) error {
	conn:= database.MongoConn

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
