package main

type Candle struct {
	Ticker    string  `json:"ticker"`
	Timestamp int64   `json:"timestamp"` // UNIX epoch (start of candle)
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	Interval  string  `json:"interval"` // e.g., "1m", "5m"
	Source    string  `json:"source"`
}
