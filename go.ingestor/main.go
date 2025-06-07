package main

import (
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
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

	tickerSymbol := "AAPL"
	interval := "1m"

	// Loop to fetch and store candles every minute
	for {
		candles, err := FetchCandles(tickerSymbol, interval)
		if err != nil {
			fmt.Println("❌ Fetch error:", err)
			continue
		}

		for _, c := range candles {
			if err := StoreCandle(rdb, c); err != nil {
				fmt.Println("❌ Store error:", err)
			} else {
				fmt.Printf("✅ Stored candle for %s at %d\n", c.Ticker, c.Timestamp)
			}
		}

		time.Sleep(60 * time.Second)
	}
}
