package internals

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis() (*redis.Client, error) {
	//redis initialization
	//get redis url
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, errors.New("REDIS url is missing in env file")

	}

	//get redis username
	redisUsername := os.Getenv("REDIS_USERNAME")
	if redisUsername == "" {
		return nil, errors.New("redisUsername is missing in env file")

	}

	//get redis passowrd
	redisUserPassword := os.Getenv("REDIS_PASSWORD")
	if redisUserPassword == "" {
		return nil, errors.New("redis userPassword is missing in env file")
	}

	//create redis client with all credientials
	redisInstance := *redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisUserPassword,
		Username: redisUsername,
	})

	// Ping the Redis server to verify the connection
	_, err := redisInstance.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("REDIS Connected successfully")
	///----------------------redis initialization end

	return &redisInstance, nil
}
