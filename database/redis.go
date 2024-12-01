package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"net/url"
	"os"
)

var RedisClient *redis.Client
var Ctx = context.Background() // Change variable name to Ctx to be more descriptive

func ConnectRedis() {
	redisURL := os.Getenv("REDIS_URL")

	// Parse the Redis URL
	options, err := parseRedisURL(redisURL)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}

	// Create a new Redis client with the parsed options
	RedisClient = redis.NewClient(options)

	// Ping Redis to verify connection
	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Connected to Redis")
}

func parseRedisURL(redisURL string) (*redis.Options, error) {
	u, err := url.Parse(redisURL)
	if err != nil {
		return nil, err
	}

	options := &redis.Options{
		Addr: u.Host,
	}

	if u.User != nil {
		password, _ := u.User.Password()
		options.Password = password
	}

	return options, nil
}
