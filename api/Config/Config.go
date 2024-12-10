package Config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type RedisClientConfig struct {
	Addr     string
	Password string
	DB       int
}

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
	KAFKA_HOST        string
	KAFKA_PORT        string
	TOPIC             string
	REDISCONN         []RedisClientConfig
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
	RedisServers := os.Getenv("REDIS_SERVER")
	RedisServerList := strings.Split(RedisServers, ",")

	var RedisClients []RedisClientConfig

	for _, val := range RedisServerList {
		RedisStruct := RedisClientConfig{
			Addr:     val,
			Password: "",
			DB:       0,
		}
		RedisClients = append(RedisClients, RedisStruct)
	}

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
		KAFKA_HOST:        os.Getenv("KAFKA_HOST"),
		KAFKA_PORT:        os.Getenv("KAFKA_PORT"),
		REDISCONN:         RedisClients,
		TOPIC:             os.Getenv("KAFKA_TOPIC"),
	}
}
