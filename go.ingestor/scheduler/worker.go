package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tmaldrsn/tradebot/go.ingestor/core"
)

type WorkerPool struct {
	JobQueue chan Job
	Quit     chan struct{}
	Redis    *redis.Client
}

// NewWorkerPool creates and starts a WorkerPool with the specified job queue size, number of worker goroutines, and a Redis client.
func NewWorkerPool(ctx context.Context, queueSize, numWorkers int, redisClient *redis.Client) *WorkerPool {
	pool := &WorkerPool{
		JobQueue: make(chan Job, queueSize),
		Quit:     make(chan struct{}),
		Redis:    redisClient,
	}

	for i := 0; i < numWorkers; i++ {
		go pool.worker(ctx, i)
	}

	log.Print("New worker pool has been created")
	return pool
}

func (p *WorkerPool) worker(ctx context.Context, id int) {
	for {
		select {
		case job := <-p.JobQueue:
			// log.Printf("[Worker %d] Fetching %s candles for %s from %s\n", id, job.Timeframe, job.Ticker, job.Ingestor.SourceName())

			// `from` and `to` cannot be real time until I update the subscription
			// we can only get end-of-day data currently
			// from := time.Now().Add(-1 * job.Interval)
			// to := time.Now()
			from, _ := time.Parse("2006-01-02", "2025-06-01")
			to, _ := time.Parse("2006-01-02", "2025-06-02")

			candles, err := job.Ingestor.FetchCandles(job.Ticker, job.Timeframe, from, to)
			if err != nil {
				log.Printf("[Worker %d] Error fetching candles: %v", id, err)
				continue
			}

			if err := core.StoreCandles(ctx, p.Redis, candles); err != nil {
				log.Printf("[Worker %d] ❌ Failed to store candles: %v", id, err)
				continue
			}

			event := core.MarketDataFetchedEvent{
				Ticker:    job.Ticker,
				Timeframe: job.Timeframe,
				Candles:   candles,
			}

			if err := core.PublishMarketData(ctx, p.Redis, event); err != nil {
				log.Printf("[Worker %d] ⚠️  Failed to publish market data: %v", id, err)
			}

			job.MarkRun(id)
		case <-p.Quit:
			log.Printf("[Worker %d] Shutting down", id)
			return
		}
	}
}

func (p *WorkerPool) Stop() {
	close(p.Quit)
}
