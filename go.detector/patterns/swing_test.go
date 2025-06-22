package patterns

import (
	"testing"

	"github.com/tmaldrsn/tradebot/go.detector/core"
)

func TestDetectSwingPoints(t *testing.T) {
	candles := []core.Candle{
		{High: 10, Low: 5, Timestamp: 1000},
		{High: 12, Low: 4, Timestamp: 2000},
		{High: 11, Low: 6, Timestamp: 3000},
	}
	swings := DetectSwingPoints(candles)

	if len(swings) != 2 {
		t.Errorf("Expected 2 swing points, got %d", len(swings))
	}
}
