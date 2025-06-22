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

	// Verify the swing points are correct
	for _, swing := range swings {
		if !(swing.Type == "swing_low" && swing.Candle.Low == 4) {
			t.Error("Expected swing low with Low=4 not found")
		}
		if !(swing.Type == "swing_high" && swing.Candle.High == 12) {
			t.Error("Expected swing high with High=12 not found")
		}
	}
}
