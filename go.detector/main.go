package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/tmaldrsn/tradebot/go.detector/handlers"
	"github.com/tmaldrsn/tradebot/go.detector/pubsub"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("cannot connect to redis at %s: %v", redisURL, err)
	}

	pubsub.Subscribe(ctx, rdb, "marketdata:fetched", handlers.HandleMarketDataMessage)

	// Wait for interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("Shutdown signal received.")
	cancel()

	if err := rdb.Close(); err != nil {
		log.Printf("error closing redis client: %v", err)
	}

}
