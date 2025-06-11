package scheduler

import (
	"log"
	"time"

	"github.com/tmaldrsn/tradebot/go.ingestor/core"
)

type Job struct {
	Ticker    string
	Timeframe string
	Ingestor  core.Ingestor
	Interval  time.Duration
	NextRun   time.Time
	Quit      chan struct{}
}

func (j *Job) ShouldRun(now time.Time) bool {
	return !now.Before(j.NextRun)
}

func (j *Job) MarkRun(workerId int) {
	j.NextRun = time.Now().Add(j.Interval)
	log.Printf("[Worker %d] Successfully fetched %s candles for %s from %s\n", workerId, j.Timeframe, j.Ticker, j.Ingestor.SourceName())
}
