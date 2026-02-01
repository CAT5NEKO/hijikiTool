package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/application/ports"
	"github.com/CAT5NEKO/hijikiTool/internal/application/usecases"
	"github.com/CAT5NEKO/hijikiTool/internal/domain"
)

type Job struct {
	ID       string
	Schedule domain.Schedule
	Content  string
}

type Scheduler struct {
	clock      ports.Clock
	repository ports.PostRecordRepository
	poster     ports.Poster
	jobs       []Job
	useCase    *usecases.SchedulePostUseCase
	tolerance  time.Duration
}

func New(
	clock ports.Clock,
	repository ports.PostRecordRepository,
	poster ports.Poster,
	jobs []Job,
) *Scheduler {
	guard := domain.NewPostGuard()
	useCase := usecases.NewSchedulePostUseCase(clock, repository, poster, guard)

	return &Scheduler{
		clock:      clock,
		repository: repository,
		poster:     poster,
		jobs:       jobs,
		useCase:    useCase,
		tolerance:  time.Minute,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	for {
		s.RunOnce()

		sleepDuration := s.NextWakeUpDuration()
		if sleepDuration < s.tolerance {
			sleepDuration = s.tolerance
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleepDuration):
		}
	}
}

func (s *Scheduler) RunOnce() {
	for _, job := range s.jobs {
		if s.useCase.ShouldExecuteNow(job.ID, job.Schedule, s.tolerance) {
			if err := s.useCase.Execute(job.ID, job.Schedule, job.Content); err != nil {
				log.Printf("Failed to execute job %s: %v", job.ID, err)
			}
		}
	}
}

func (s *Scheduler) NextWakeUpDuration() time.Duration {
	now := s.clock.Now()
	var minDuration time.Duration

	for i, job := range s.jobs {
		duration := job.Schedule.DurationUntil(now)
		if i == 0 || duration < minDuration {
			minDuration = duration
		}
	}

	if minDuration <= 0 {
		return s.tolerance
	}

	return minDuration
}
