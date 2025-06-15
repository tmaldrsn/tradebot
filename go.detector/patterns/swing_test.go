package patterns

import (
	"testing"

	"github.com/tmaldrsn/tradebot/go.detector/core"
)

func TestDetectSwingPoints(t *testing.T) {
	candles := []core.Candle{
		{High: 10, Low: 5},
		{High: 12, Low: 4},
		{High: 11, Low: 6},
	}
	swings := DetectSwingPoints(candles)

	if len(swings) != 2 {
		t.Errorf("Expected 2 swing points, got %d", len(swings))
	}
}
