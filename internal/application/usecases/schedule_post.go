package usecases

import (
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/application/ports"
	"github.com/CAT5NEKO/hijikiTool/internal/domain"
)

type SchedulePostUseCase struct {
	clock      ports.Clock
	repository ports.PostRecordRepository
	poster     ports.Poster
	guard      *domain.PostGuard
}

func NewSchedulePostUseCase(
	clock ports.Clock,
	repository ports.PostRecordRepository,
	poster ports.Poster,
	guard *domain.PostGuard,
) *SchedulePostUseCase {
	return &SchedulePostUseCase{
		clock:      clock,
		repository: repository,
		poster:     poster,
		guard:      guard,
	}
}

func (u *SchedulePostUseCase) Execute(scheduleID string, schedule domain.Schedule, content string) error {
	now := u.clock.Now()
	record, err := u.repository.Find(scheduleID)
	if err != nil {
		return err
	}

	if !u.guard.CanPost(schedule, record, now) {
		return nil
	}

	if err := u.poster.Post(content); err != nil {
		return err
	}

	newRecord := domain.NewPostRecord(scheduleID, now)
	return u.repository.Save(newRecord)
}

func (u *SchedulePostUseCase) ShouldExecuteNow(scheduleID string, schedule domain.Schedule, tolerance time.Duration) bool {
	now := u.clock.Now()
	record, err := u.repository.Find(scheduleID)
	if err != nil {
		return false
	}

	if !u.guard.CanPost(schedule, record, now) {
		return false
	}

	nextTime := schedule.NextTime(now.Add(-tolerance))
	return u.isWithinTolerance(now, nextTime, tolerance)
}

func (u *SchedulePostUseCase) isWithinTolerance(now, target time.Time, tolerance time.Duration) bool {
	diff := now.Sub(target)
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
