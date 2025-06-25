package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/tmaldrsn/tradebot/go.detector/core"
	"github.com/tmaldrsn/tradebot/go.detector/patterns"
)

type MarketDataFetchedEvent struct {
	Ticker    string `json:"ticker"`
	Timeframe string `json:"timeframe"`
}

func HandleMarketDataMessage(payload string, rdb *redis.Client) {
	var event MarketDataFetchedEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Printf("❌ Failed to parse candle message: %v", err)
		return
	}

	log.Printf("🕯️ Detected event: %+v", event)

	// TODO: Run pattern detection here (swing point, etc.)
	candles, err := GetCandlesByTickerAndTimeframe(rdb, event.Ticker, event.Timeframe)
	if err != nil {
		log.Printf("❌ Failed to read candles from redis: %v", err)
	}

	swingPoints := patterns.DetectSwingPoints(candles)
	log.Print(swingPoints)
}

func GetCandlesByTickerAndTimeframe(rdb *redis.Client, ticker string, timeframe string) ([]core.Candle, error) {
	ctx := context.Background()
	pattern := fmt.Sprintf("candle:%s:%s:*", ticker, timeframe)

	// Step 1: Collect keys matching pattern
	var keys []string
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	// Step 2: Use MGET to fetch all values in one round trip
	values, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("mget error: %w", err)
	}

	// Step 3: Deserialize JSON values into Candle structs
	candles := make([]core.Candle, 0, len(values))
	for i, v := range values {
		if v == nil {
			log.Printf("⚠️ Missing value for key: %s", keys[i])
			continue
		}

		strVal, ok := v.(string)
		if !ok {
			log.Printf("❌ Unexpected value type for key %s", keys[i])
			continue
		}

		var candle core.Candle
		if err := json.Unmarshal([]byte(strVal), &candle); err != nil {
			log.Printf("❌ Failed to unmarshal candle for key %s: %v", keys[i], err)
			continue
		}
		candles = append(candles, candle)
	}

	return candles, nil
}
