package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func StoreCandle(ctx context.Context, rdb *redis.Client, candle Candle) error {
	key := fmt.Sprintf("candle:%s:%s:%d", candle.Ticker, candle.Timeframe, candle.Timestamp)

	data, err := json.Marshal(candle)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, data, 0).Err()
}

func StoreCandles(ctx context.Context, rdb *redis.Client, candles []Candle) error {
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, c := range candles {
			key := fmt.Sprintf("candle:%s:%s:%d", c.Ticker, c.Timeframe, c.Timestamp)
			data, err := json.Marshal(c)
			if err != nil {
				return err
			}
			pipe.Set(ctx, key, data, 0)
		}
		return nil
	})
	return err
}
