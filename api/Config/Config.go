package Config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ChildLocks        string
	ParentPath        string
	GlobalCounterPath string
	ParentLockPath    string
	ServerPath        string
	ServerPort        string
	ServerRangePath   string
	Servers           []string
	SessionTimeout    time.Duration
	CounterRange      string
	AwsRegion         string
	TableName         string
	DaxEndpoint       string
	UseDax            string
	ExpiryTime        time.Duration
	APP_PORT          string
}

var Cfg AppConfig

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	ServerEnv := os.Getenv("SERVER")
	ServerList := strings.Split(ServerEnv, ",")

	// Parse SESSION_TIMEOUT environment variable into a time.Duration
	sessionTimeout, err := time.ParseDuration(os.Getenv("SESSION_TIMEOUT"))
	if err != nil {
		sessionTimeout = 30 * time.Minute // Default to 30 minutes if parsing fails
	}

	expiryTime, err := time.ParseDuration(os.Getenv("EXPIRY_TIME"))
	if err != nil {
		expiryTime = 24 * time.Hour * 366 * 10 // Default to 10 years if parsing fails
	}

	Cfg = AppConfig{
		ChildLocks:        os.Getenv("CHILD_LOCKS"),
		ParentPath:        os.Getenv("PARENT_PATH"),
		GlobalCounterPath: os.Getenv("GLOBAL_COUNTER_PATH"),
		ParentLockPath:    os.Getenv("PARENT_LOCK_PATH"),
		ServerPath:        os.Getenv("SERVER_PATH"),
		ServerPort:        os.Getenv("SERVER_ID"),
		SessionTimeout:    sessionTimeout,
		CounterRange:      os.Getenv("COUNTER_RANGE"),
		Servers:           ServerList,
		AwsRegion:         os.Getenv("AWS_REGION"),
		TableName:         os.Getenv("TABLE_NAME"),
		DaxEndpoint:       os.Getenv("DAX_ENDPOINT"),
		UseDax:            os.Getenv("USE_DAX"),
		ExpiryTime:        expiryTime,
		APP_PORT:          os.Getenv("APP_PORT"),
	}
}
