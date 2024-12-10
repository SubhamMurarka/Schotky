package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	KAFKA_HOST string
	KAFKA_PORT string
	TOPIC      string
	ES_URL     string
}

var Cfg AppConfig

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	Cfg = AppConfig{
		KAFKA_HOST: os.Getenv("KAFKA_HOST"),
		KAFKA_PORT: os.Getenv("KAFKA_PORT"),
		TOPIC:      os.Getenv("KAFKA_TOPIC"),
		ES_URL:     os.Getenv("ES_URL"),
	}
}
