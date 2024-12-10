package main

import (
	"log"

	"github.com/SubhamMurarka/Schotky/Config"
	es "github.com/SubhamMurarka/Schotky/ES"
	kafka "github.com/SubhamMurarka/Schotky/Kafka"
	models "github.com/SubhamMurarka/Schotky/Models"
)

func main() {
	// Load configurations
	kConfig := models.KafkaConfig{
		Host:  Config.Cfg.KAFKA_HOST,
		Port:  Config.Cfg.KAFKA_PORT,
		Topic: Config.Cfg.TOPIC,
	}

	// Initialize Elasticsearch client
	esClient, err := es.InitializeElasticsearch(Config.Cfg.ES_URL)
	if err != nil {
		log.Printf("Failed to initialize Elasticsearch: %v", err)
	}

	// Initialize Kafka consumer and start processing
	err = kafka.StartConsumer(kConfig, esClient, 10) // Batch size of 10
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
}
