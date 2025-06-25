package patterns

import (
	"sort"

	"github.com/tmaldrsn/tradebot/go.detector/core"
)

const (
	SwingHigh = "swing_high"
	SwingLow  = "swing_low"
)

type SwingPoint struct {
	Candle core.Candle
	Type   string // "SwingHigh" or "SwingLow"
}

func DetectSwingPoints(candles []core.Candle) []SwingPoint {
	if len(candles) < 3 {
		return []SwingPoint{}
	}

	// Create a copy to avoid modifying the input slice
	sortedCandles := make([]core.Candle, len(candles))
	copy(sortedCandles, candles)

	sort.Slice(sortedCandles, func(i, j int) bool {
		return candles[i].Timestamp < candles[j].Timestamp
	})

	var swings []SwingPoint
	for i := 1; i < len(sortedCandles)-1; i++ {
		prev := sortedCandles[i-1]
		curr := sortedCandles[i]
		next := sortedCandles[i+1]

		if curr.Low < prev.Low && curr.Low < next.Low {
			swings = append(swings, SwingPoint{Candle: curr, Type: SwingLow})
		}
		if curr.High > prev.High && curr.High > next.High {
			swings = append(swings, SwingPoint{Candle: curr, Type: SwingHigh})
		}
	}
	return swings
}
