package core

type MarketDataFetchedEvent struct {
	Ticker    string   `json:"ticker"`
	Timeframe string   `json:"timeframe"`
	Candles   []Candle `json:"candles"`
}
