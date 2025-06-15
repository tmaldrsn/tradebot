package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/tmaldrsn/tradebot/go.detector/core"
)

type MarketDataFetchedEvent struct {
	Ticker    string `json:"ticker"`
	Timeframe string `json:"timeframe"`
}

func HandleMarketDataMessage(payload string) {
	var event MarketDataFetchedEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Printf("❌ Failed to parse candle message: %v", err)
		return
	}

	log.Printf("🕯️ Detected event: %+v", event)

	// TODO: Run pattern detection here (swing point, etc.)
}

func GetCandlesByTickerAndTimeframe(rdb *redis.Client, ticker string, timeframe string) ([]core.Candle, error) {
	var candles []core.Candle

	ctx := context.Background()
	pattern := fmt.Sprintf("candle:%s:%s:*", ticker, timeframe)

	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			log.Printf("❌ Failed to get key %s: %v", key, err)
			continue
		}

		var candle core.Candle
		if err := json.Unmarshal([]byte(val), &candle); err != nil {
			log.Printf("❌ Failed to unmarshal candle: %v", err)
			continue
		}

		candles = append(candles, candle)
	}

	if err := iter.Err(); err != nil {
		log.Fatalf("iteration error: %w", err)
	}

	return candles, nil
}
