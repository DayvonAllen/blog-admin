package repo

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/kafkaProducer"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack/v5"
)

func ProcessMessage(message domain.Message) error {

	//if message.ResourceType == "posts" {
	//	// 201 is the created messageType
	//	if message.MessageType == 201 {
	//		user := message.User
	//		err := UserRepoImpl{}.Create(&user)
	//
	//		if err != nil {
	//			return err
	//		}
	//		return nil
	//	}
	//
	//	// 200 is the updated messageType
	//	if message.MessageType == 200 {
	//		user := message.User
	//
	//		err := UserRepoImpl{}.UpdateByID(&user)
	//		if err != nil {
	//			return err
	//		}
	//		return nil
	//	}
	//
	//	// 204 is the deleted messageType
	//	if message.MessageType == 204 {
	//		user := message.User
	//
	//		err := UserRepoImpl{}.DeleteByID(&user)
	//
	//		if err != nil {
	//			return err
	//		}
	//		return nil
	//	}
	//}

	return fmt.Errorf("cannot process this message")
}

func PushUserToQueue(message []byte, topic string) error {
	producer := kafkaProducer.GetInstance()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		err = producer.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("Failed to send message to the queue")
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "post", partition, offset)
	return nil
}

func SendKafkaMessage(post *domain.Post, eventType int) error {
	um := new(domain.Message)
	um.Post = *post

	// user created/updated event
	um.MessageType = eventType
	um.ResourceType = "post"

	//turn user struct into a byte array
	b, err := msgpack.Marshal(um)

	if err != nil {
		return err
	}

	err = PushUserToQueue(b, "post")

	if err != nil {
		return err
	}

	return nil
}

func SendAltKafkaMessage(tag *domain.Tag, eventType int) error {
	um := new(domain.Message)
	um.Tag = *tag

	// user created/updated event
	um.MessageType = eventType
	um.ResourceType = "tag"

	//turn user struct into a byte array
	b, err := msgpack.Marshal(um)

	if err != nil {
		return err
	}

	err = PushUserToQueue(b, "tag")

	if err != nil {
		return err
	}

	return nil
}