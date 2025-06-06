package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis:6379"
	}

	opts := &redis.Options{
		Addr: redisURL,
		DB:   0,
	}

	client := redis.NewClient(opts)

	// ping test
	err := client.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("✅ Connected to Redis")

	// Example usage
	err = client.Set(ctx, "foo", "bar", time.Minute).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value for 'foo': %s\r", val)
}
