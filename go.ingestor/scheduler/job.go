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
	return j.NextRun.Sub(now) >= j.Interval
}

func (j *Job) MarkRun(worker_id int) {
	j.NextRun = time.Now().Add(j.Interval)
	log.Printf("[Worker %d] Successfully fetched %s candles for %s from %s\n", worker_id, j.Timeframe, j.Ticker, j.Ingestor.SourceName())
}
