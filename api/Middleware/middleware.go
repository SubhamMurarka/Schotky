package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SubhamMurarka/Schotky/Config"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/spaolacci/murmur3"
)

var Clients []*redis.Client

// InitializeRedisClients initializes and returns three Redis clients
func InitializeRedisClients() {
	Clients = make([]*redis.Client, 0)

	for _, config := range Config.Cfg.REDISCONN {
		client := redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password, // no password set if empty
			DB:       config.DB,       // use default DB
		})

		// Check the connection
		ctx := context.Background()
		_, err := client.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis at %s: %v", config.Addr, err)
		}

		Clients = append(Clients, client)
	}
}

// computeBucket determines the Redis bucket for a given IP
func computeBucket(ip string) uint32 {
	hash := murmur3.New32()
	hash.Write([]byte(ip))
	res := hash.Sum32()
	print("hashed value : ", res)
	return res % 3
}

// UpdateRateLimit checks and updates the rate limit for an IP
func UpdateRateLimit(ip string) error {
	bucket := computeBucket(ip)
	redisClient := Clients[bucket]
	println("bucket : ", bucket)

	println("ip : ", ip)
	println("redis is : ", redisClient)

	// Increment and set expiration if not already set
	incrementResult, err := redisClient.Incr(context.Background(), ip).Result()
	if err != nil {
		return fmt.Errorf("failed to increment Redis key: %w", err)
	}

	fmt.Println("incremented Result ", incrementResult)

	// Set expiration to 10 minutes if this is the first increment
	if incrementResult == 1 {
		_, err = redisClient.Expire(context.Background(), ip, 10*time.Second).Result()
		if err != nil {
			return fmt.Errorf("failed to set expiration: %w", err)
		}
	}

	// Check the limit (assuming 100 requests per 10 minutes)
	if incrementResult > 300 {
		return fiber.ErrTooManyRequests
	}

	return nil
}

// CloseRedisClients closes all Redis client connections
func CloseRedisClients() {
	for _, client := range Clients {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v", err)
		}
	}
}

// RateLimitMiddleware is a Fiber middleware for rate-limiting
func RateLimitMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Fetch the IP address from the request context
		ip := ctx.Get("X-Forwarded-For")
		if ip == " " {
			log.Print("ip missing")
		}

		// Call UpdateRateLimit with the IP address
		err := UpdateRateLimit(ip)
		if err != nil {
			if err == fiber.ErrTooManyRequests {
				// If rate limit is exceeded, return a 429 response
				return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too many requests, please try again later.",
				})
			}
			// Handle other errors
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error.",
			})
		}
		// Proceed to the next handler
		return ctx.Next()
	}
}
