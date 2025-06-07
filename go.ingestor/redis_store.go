package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func StoreCandle(rdb *redis.Client, candle Candle) error {
	key := fmt.Sprintf("candle:%s:%s:%d", candle.Ticker, candle.Interval, candle.Timestamp)

	data, err := json.Marshal(candle)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, data, 0).Err()
}
