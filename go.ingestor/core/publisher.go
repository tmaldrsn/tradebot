package core

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

const MarketDataFetchedTopic = "marketdata:fetched"

func PublishMarketData(ctx context.Context, rdb *redis.Client, event MarketDataFetchedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = rdb.Publish(ctx, MarketDataFetchedTopic, data).Err()
	if err != nil {
		return err
	}

	return nil
}
