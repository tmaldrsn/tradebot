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
	swingHighFound := false
	swingLowFound := false
	for _, swing := range swings {
		if swing.Type == SwingLow && swing.Candle.Low == 4 {
			swingLowFound = true
		}
		if swing.Type == SwingHigh && swing.Candle.High == 12 {
			swingHighFound = true
		}
	}

	if !swingLowFound {
		t.Error("Expected swing low with Low=4 not found")
	}
	if !swingHighFound {
		t.Error("Expected swing high with High=12 not found")
	}
}

func TestDetectSwingPointsEdgeCases(t *testing.T) {
	// Test with empty slice
	t.Run("EmptyCandles", func(t *testing.T) {
		swings := DetectSwingPoints([]core.Candle{})
		if len(swings) != 0 {
			t.Errorf("Expected 0 swing points for empty slice, got %d", len(swings))
		}
	})

	// Test with insufficient candles
	t.Run("InsufficientCandles", func(t *testing.T) {
		candles := []core.Candle{
			{High: 10, Low: 5, Timestamp: 1000},
			{High: 12, Low: 4, Timestamp: 2000},
		}
		swings := DetectSwingPoints(candles)
		if len(swings) != 0 {
			t.Errorf("Expected 0 swing points for 2 candles, got %d", len(swings))
		}
	})

	// Test with no swing points
	t.Run("NoSwingPoints", func(t *testing.T) {
		candles := []core.Candle{
			{High: 10, Low: 5, Timestamp: 1000},
			{High: 11, Low: 6, Timestamp: 2000},
			{High: 12, Low: 7, Timestamp: 3000},
		}
		swings := DetectSwingPoints(candles)
		if len(swings) != 0 {
			t.Errorf("Expected 0 swing points for ascending trend, got %d", len(swings))
		}
	})
}
