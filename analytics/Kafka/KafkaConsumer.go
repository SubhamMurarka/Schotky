package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	es "github.com/SubhamMurarka/Schotky/ES"
	models "github.com/SubhamMurarka/Schotky/Models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ip2location/ip2location-go"
	"github.com/mssola/user_agent"
)

var db4 *ip2location.DB
var db6 *ip2location.DB

// StartConsumer initializes the Kafka consumer and processes messages in batches.
func StartConsumer(kConfig models.KafkaConfig, esClient *elasticsearch.Client, batchSize int) error {
	// Initialize Kafka consumer
	consumer, err := initializeKafkaConsumer(kConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize Kafka consumer: %w", err)
	}
	defer consumer.Close()

	// db4, err = ip2location.OpenDB("analytics/Kafka/IP2LOCATION-LITE-DB3.BIN")
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer db4.Close()

	// db6, err = ip2location.OpenDB("analytics/Kafka/IP2LOCATION-LITE-DB3.IPV6.BIN")
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer db6.Close()

	return processMessages(consumer, kConfig.Topic, esClient, batchSize)
}

func initializeKafkaConsumer(kConfig models.KafkaConfig) (sarama.Consumer, error) {
	url := fmt.Sprintf("%s:%s", kConfig.Host, kConfig.Port)
	brokerUrls := []string{url}
	fmt.Println(url)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	var consumer sarama.Consumer
	var err error
	for retries := 10; retries > 0; retries-- {
		consumer, err = sarama.NewConsumer(brokerUrls, config)
		if err == nil {
			log.Println("kafka consumer connected successfully")
			return consumer, nil
		}
		log.Printf("Error creating Kafka consumer, retrying: %v", err)
		time.Sleep(10 * time.Second)
	}
	return nil, fmt.Errorf("failed to create Kafka consumer after retries: %w", err)
}

func processMessages(consumer sarama.Consumer, topic string, esClient *elasticsearch.Client, batchSize int) error {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("error starting partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	// Signal handling for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	batch := make([]models.ProcessedData, 0, batchSize)

	for {
		select {
		case <-stopChan:
			log.Println("Received shutdown signal, stopping consumer.")
			return nil
		case message := <-partitionConsumer.Messages():
			var data models.AnalyticsData
			if err := json.Unmarshal(message.Value, &data); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			fmt.Println("data : ", data)
			processed := processAnalyticsData(data)
			fmt.Println("processed : ", processed)
			batch = append(batch, processed)

			if len(batch) >= batchSize {
				if err := es.WriteBatchToElasticsearch(batch, esClient); err != nil {
					log.Printf("Error writing batch to Elasticsearch: %v", err)
				}
				batch = batch[:0]
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming messages: %v", err)
		}
	}
}

func processAnalyticsData(data models.AnalyticsData) models.ProcessedData {
	var dataObj models.ProcessedData

	// IpDetail := GetIpDetails(data.IPAddress)

	// dataObj.Country = IpDetail.Country
	// dataObj.City = IpDetail.City

	if data.Referrer != "Not Available" {
		dataObj.Referrer = data.Referrer
	}

	if data.ShortURL != "Not Available" {
		dataObj.ShortURL = data.ShortURL
	}

	dataObj.Timestamp = data.Timestamp

	ua := user_agent.New(data.UserAgent)
	browser, _ := ua.Browser()
	dataObj.Browser = browser
	dataObj.OS = ua.OS()
	dataObj.Device = ua.Platform()

	return dataObj
}

func GetIpDetails(ip string) models.IP {
	ipV := checkIPAddress(ip)

	if ipV == "IPv4" {
		record, err := db4.Get_all(ip)
		if err != nil {
			log.Println(err)
		}
		return models.IP{
			Country: record.Country_short,
			City:    record.City,
			Region:  record.Region,
		}
	} else if ipV == "IPv6" {
		record, err := db6.Get_all(ip)
		if err != nil {
			log.Println(err)
		}
		return models.IP{
			Country: record.Country_short,
			City:    record.City,
			Region:  record.Region,
		}
	} else {
		return models.IP{}
	}
}

func checkIPAddress(ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "Invalid"
	}
	if parsedIP.To4() != nil {
		return "IPv4"
	}
	return "IPv6"
}
