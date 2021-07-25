package kafkaProducer

import (
	"github.com/Shopify/sarama"
	"strings"
	"sync"
)

type Connection struct {
	sarama.SyncProducer
}

var kafkaConnection *Connection
var once sync.Once

// GetInstance creates one instance and always returns that one instance
func GetInstance() *Connection {
	// only executes this once
	once.Do(func() {
		_, err := connectProducer()
		if err != nil {
			panic(err)
		}
	})
	return kafkaConnection
}

func connectProducer() (sarama.SyncProducer,error) {
	//names := os.Getenv("KAFKA_SERVICE_NAMES")

	namesArr := make([]string, 0, 10)

	for _, v := range strings.Split("kafka:9092", ",") {
		namesArr = append(namesArr, v)
	}

	brokersUrl := namesArr

	newConfig := sarama.NewConfig()
	newConfig.Producer.Return.Successes = true
	newConfig.Producer.RequiredAcks = sarama.WaitForAll
	newConfig.Producer.Retry.Max = 7
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	conn, err := sarama.NewSyncProducer(brokersUrl, newConfig)
	if err != nil {
		panic(err)
	}

	kafkaConnection = &Connection{conn}

	return kafkaConnection, nil
}