package scheduler

import (
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

func (j *Job) MarkRun() {
	j.NextRun = time.Now().Add(j.Interval)
}
