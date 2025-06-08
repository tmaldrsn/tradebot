package core

import "time"

type Ingestor interface {
	FetchCandles(ticker, timeframe string, from, to time.Time) ([]Candle, error)
	SourceName() string
	SourceType() string
}
