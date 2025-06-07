package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

func FetchCandles(ticker string, interval string) ([]Candle, error) {
	api_key := os.Getenv("POLYGON_API_KEY")
	c := polygon.New(api_key)

	from, err := time.Parse("2006-01-02", "2025-02-09")
	if err != nil {
		log.Fatalf("Error parsing 'from' date: %v", err)
	}
	to, err := time.Parse("2006-01-02", "2025-02-10")
	if err != nil {
		log.Fatalf("Error parsing 'to' date: %v", err)
	}

	params := models.ListAggsParams{
		Ticker:     "X:BTCUSD",
		Multiplier: 5,
		Timespan:   "minute",
		From:       models.Millis(from),
		To:         models.Millis(to),
	}.
		WithAdjusted(true).
		WithOrder(models.Order("asc")).
		WithLimit(1000)

	iter := c.ListAggs(context.Background(), params)

	candles := []Candle{}

	for iter.Next() {
		result := iter.Item()
		candles = append(candles, Candle{
			Ticker:    ticker,
			Timestamp: time.Time(result.Timestamp).UnixMilli() / 1000, // ms → s
			Open:      result.Open,
			High:      result.High,
			Low:       result.Low,
			Close:     result.Close,
			Volume:    result.Volume,
			Interval:  "5m",
			Source:    "polygon",
		})
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over candles: %w", err)
	}

	return candles, nil
}
