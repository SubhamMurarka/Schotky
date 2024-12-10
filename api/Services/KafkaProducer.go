package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/SubhamMurarka/Schotky/Config"
	models "github.com/SubhamMurarka/Schotky/Models"
)

var producer sarama.SyncProducer

func InitKafkaProducer() {
	Host := Config.Cfg.KAFKA_HOST
	Port := Config.Cfg.KAFKA_PORT

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 2
	config.Producer.Flush.Messages = 1000
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.Return.Successes = true
	fmt.Println(Host, Port)
	var err error
	for i := 0; i < 5; i++ { // Try to connect up to 5 times
		producer, err = sarama.NewSyncProducer([]string{Host + ":" + Port}, config)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to Kafka, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to Kafka after retries:", err)
	}
}

func PublishMessage(message *models.AnalyticsData) {
	Topic := Config.Cfg.TOPIC
	log.Printf("Kafka topic: %s", Topic)
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	fmt.Println("mbytes: ", messageBytes)
	fmt.Println("messages : ", message)

	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.ByteEncoder(messageBytes),
	})

	if err != nil {
		log.Println("Error publishing message to Kafka:", err)
		return
	}

	fmt.Printf("data inserted in partition %d and with offset value %d", partition, offset)
}

func CloseKafka() {
	var err error
	if producer != nil {
		err = producer.Close()
	}
	if err != nil {
		fmt.Println("Kafka not closed", err)
	}
}
