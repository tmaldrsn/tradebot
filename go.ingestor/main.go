package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/tmaldrsn/tradebot/go.ingestor/config"
	"github.com/tmaldrsn/tradebot/go.ingestor/scheduler"
	polygonrest "github.com/tmaldrsn/tradebot/go.ingestor/sources/polygon/rest"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})

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

	pool := scheduler.NewWorkerPool(100, 4, rdb)
	sched := scheduler.NewScheduler(pool, jobs)
	sched.Start()

	// Wait for interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("Shutdown signal received.")
	pool.Stop()
	sched.Stop()

	if err := rdb.Close(); err != nil {
		log.Printf("error closing redis client: %v", err)
	}
}
