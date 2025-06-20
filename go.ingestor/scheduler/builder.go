package scheduler

import (
	"log"
	"time"

	"github.com/tmaldrsn/tradebot/go.ingestor/config"
	"github.com/tmaldrsn/tradebot/go.ingestor/core"
)

// BuildScheduledJobs converts loaded config + source-specific ingestor into scheduled jobs
func BuildScheduledJobs(cfg *config.Config, ingestor core.Ingestor) []ScheduledJob {
	var jobs []ScheduledJob

	for _, source := range cfg.Sources {
		if source.Name != ingestor.SourceName() {
			continue // skip if not matching the ingestor we’re wiring
		}
		for _, ticker := range source.Tickers {
			// interval, err := time.ParseDuration(ticker.PollingInterval)
			interval, err := time.ParseDuration(ticker.Timeframe)
			if err != nil {
				panic(err)
			}

			job := Job{
				Ticker:    ticker.Ticker,
				Timeframe: ticker.Timeframe,
				Ingestor:  ingestor,
			}
			jobs = append(jobs, ScheduledJob{
				Job:      job,
				Interval: interval,
			})
			log.Printf("Built a job: %s, %s, %s", ticker.Ticker, ticker.Timeframe, interval)
		}
	}
	return jobs
}
