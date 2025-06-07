package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Job struct {
	Ticker    string
	Timeframe string
}

type TickerConfig struct {
	Ticker       string
	Timeframe    string
	PollInterval time.Duration // how often to poll
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup, rdb *redis.Client) {
	defer wg.Done()

	for job := range jobs {
		log.Printf("[Worker %d] Processing %s", id, job.Ticker)

		candles, err := FetchCandles(job.Ticker, job.Timeframe)
		if err != nil {
			log.Printf("[Worker %d] Error fetching for %s: %v", id, job.Ticker, err)
			continue
		}

		StoreCandles(rdb, candles)
	}
}

func startScheduler(configs []TickerConfig, jobQueue chan<- Job) {
	for _, cfg := range configs {
		go func(c TickerConfig) {
			ticker := time.NewTicker(c.PollInterval)
			defer ticker.Stop()

			for {
				jobQueue <- Job{Ticker: c.Ticker, Timeframe: c.Timeframe}
				<-ticker.C
			}
		}(cfg)
	}
}

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})

	configs := []TickerConfig{
		{"X:BTCUSD", "1h", 1 * time.Minute},
		{"X:ETHUSD", "30m", 2 * time.Minute},
		{"AAPL", "10m", 5 * time.Minute},
	}

	jobQueue := make(chan Job, 100)
	startScheduler(configs, jobQueue)

	var wg sync.WaitGroup
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobQueue, &wg, rdb)
	}

	select {} // Block forever or implement graceful shutdown
}
