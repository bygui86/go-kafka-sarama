package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"samples/go-kafka/producer/user/domain"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gocql/gocql"
)

const (
	kafkaConn = "dockerhost:9092"
	topic     = "gomessages"
)

func main() {
	// create producer
	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	// read command line input
	var counter int = 1
	for {
		// fmt.Print("Enter msg: ")
		// msg, _ := reader.ReadString('\n')

		// publish without goroutine
		user := domain.User{
			ID:        gocql.TimeUUID(),
			FirstName: "john",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Age:       counter,
			City:      "NY",
		}
		publish(user, producer)

		// publish with goroutine
		// go publish(msg, producer)
		counter++

		time.Sleep(15 * time.Second)
	}
}

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// async producer
	//prd, err := sarama.NewAsyncProducer([]string{kafkaConn}, config)

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

func publish(user domain.User, producer sarama.SyncProducer) {

	message, jsonErr := json.Marshal(user)
	if jsonErr != nil {
		fmt.Println("Error json marshal:", jsonErr.Error())
		panic(jsonErr)
	}

	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, publishErr := producer.SendMessage(msg)
	if publishErr != nil {
		fmt.Println("Error publish:", publishErr.Error())
		panic(publishErr)
	}

	// publish async
	//producer.Input() <- &sarama.ProducerMessage{

	fmt.Println("Partition: ", partition)
	fmt.Println("Offset: ", offset)
	fmt.Println("User firstname:", user.FirstName)
	fmt.Println("User lastname:", user.LastName)
	fmt.Println("User age:", user.Age)
}
