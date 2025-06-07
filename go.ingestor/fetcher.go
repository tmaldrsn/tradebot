package main

import (
	"time"
)

func FetchCandles(ticker string, interval string) ([]Candle, error) {
	now := time.Now().Truncate(time.Minute)
	return []Candle{
		{
			Ticker:    ticker,
			Timestamp: now.Unix(),
			Open:      100.0,
			High:      105.0,
			Low:       98.0,
			Close:     102.5,
			Volume:    1200,
			Interval:  interval,
			Source:    "mock",
		},
	}, nil
}
