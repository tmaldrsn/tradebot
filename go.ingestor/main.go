package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/tmaldrsn/tradebot/go.ingestor/config"
	"github.com/tmaldrsn/tradebot/go.ingestor/scheduler"
	polygonrest "github.com/tmaldrsn/tradebot/go.ingestor/sources/polygon/rest"
)

// main initializes the job processing system, loads configuration, sets up Redis, starts scheduled jobs, and handles graceful shutdown on interrupt signals.
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

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH env var not set")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// cfg, err := cfg.GetSource("polygon", "rest")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ingestor := polygonrest.NewIngestor()
	jobs := scheduler.BuildScheduledJobs(cfg, ingestor)

	pool := scheduler.NewWorkerPool(ctx, 100, 4, rdb)
	sched := scheduler.NewScheduler(pool, jobs)
	sched.Start()

	// Wait for interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("Shutdown signal received.")
	cancel()
	pool.Stop()
	sched.Stop()

	if err := rdb.Close(); err != nil {
		log.Printf("error closing redis client: %v", err)
	}
}
