package polygonrest

import (
	"context"
	"os"
	"time"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"

	"github.com/tmaldrsn/tradebot/go.ingestor/core"
	"github.com/tmaldrsn/tradebot/go.ingestor/utils"
)

type PolygonRestIngestor struct {
	client *polygon.Client
}

func NewIngestor() *PolygonRestIngestor {
	apiKey := os.Getenv("POLYGON_API_KEY")
	if apiKey == "" {
		panic("No Polygon API key found!!")
	}

	client := polygon.New(apiKey)
	return &PolygonRestIngestor{
		client: client,
	}
}

func (p *PolygonRestIngestor) SourceName() string {
	return "polygon"
}

func (p *PolygonRestIngestor) SourceType() string {
	return "rest"
}

func (p *PolygonRestIngestor) FetchCandles(ticker, timeframe string, from, to time.Time) ([]core.Candle, error) {
	pTimeframe, err := utils.ParseTimeframe(timeframe)
	if err != nil {
		panic(err)
	}

	params := models.ListAggsParams{
		Ticker:     ticker,
		Multiplier: pTimeframe.Multiplier,
		Timespan:   pTimeframe.Timespan,
		From:       models.Millis(from),
		To:         models.Millis(to),
	}.
		WithAdjusted(true).
		WithOrder(models.Order("asc")).
		WithLimit(5000)

	iter := p.client.ListAggs(context.Background(), params)

	var candles []core.Candle

	for iter.Next() {
		result := iter.Item()
		candle := core.Candle{
			Ticker:    ticker,
			Timestamp: time.Time(result.Timestamp).UnixMilli() / 1000,
			Open:      result.Open,
			High:      result.High,
			Low:       result.Low,
			Close:     result.Close,
			Volume:    result.Volume,
			Timeframe: timeframe,
			Source:    p.SourceName(),
		}
		candles = append(candles, candle)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return candles, nil
}
