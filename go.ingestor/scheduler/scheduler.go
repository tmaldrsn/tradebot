package scheduler

import (
	"log"
	"time"
)

type Scheduler struct {
	pool *WorkerPool
	jobs []ScheduledJob
	quit chan struct{}
}

type ScheduledJob struct {
	Job      Job
	Interval time.Duration
}

func NewScheduler(pool *WorkerPool, jobs []ScheduledJob) *Scheduler {
	return &Scheduler{
		pool: pool,
		jobs: jobs,
		quit: make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Second) // check every second
	defer ticker.Stop()

	for {
		select {
		case <-s.quit:
			log.Println("[Scheduler] Shutting down")
			return
		case now := <-ticker.C:
			for i := range s.jobs {
				job := &s.jobs[i]
				if now == job.Job.NextRun || now.After(job.Job.NextRun) {
					s.pool.JobQueue <- job.Job
					job.Job.NextRun = now.Add(job.Interval)
				}
			}
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.quit)
}
