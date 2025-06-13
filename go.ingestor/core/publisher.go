package core

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

var pubsubCtx = context.Background()

const MarketDataFetchedTopic = "marketdata:fetched"

func PublishMarketData(rdb *redis.Client, event MarketDataFetchedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = rdb.Publish(pubsubCtx, MarketDataFetchedTopic, data).Err()
	if err != nil {
		log.Printf("❌ Failed to publish market data event: %v", err)
		return err
	}

	return nil
}
