package handlers

import (
	"encoding/json"
	"log"
)

type MarketDataFetchedEvent struct {
	Ticker    string `json:"ticker"`
	Timeframe string `json:"timeframe"`
}

func HandleMarketDataMessage(payload string) {
	var event MarketDataFetchedEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Printf("❌ Failed to parse candle message: %v", err)
		return
	}

	log.Printf("🕯️ Detected event: %+v", event)

	// TODO: Run pattern detection here (swing point, etc.)
}
