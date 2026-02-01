package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/application/usecases"
	"github.com/CAT5NEKO/hijikiTool/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeClock struct {
	fixedTime time.Time
}

func (c *FakeClock) Now() time.Time {
	return c.fixedTime
}

type FakePostRecordRepository struct {
	records    map[string]domain.PostRecord
	saveError  error
	saveCalled bool
	savedRecord domain.PostRecord
}

func NewFakePostRecordRepository() *FakePostRecordRepository {
	return &FakePostRecordRepository{records: make(map[string]domain.PostRecord)}
}

func (r *FakePostRecordRepository) Find(scheduleID string) (domain.PostRecord, error) {
	record, exists := r.records[scheduleID]
	if !exists {
		return domain.PostRecord{}, nil
	}
	return record, nil
}

func (r *FakePostRecordRepository) Save(record domain.PostRecord) error {
	r.saveCalled = true
	r.savedRecord = record
	if r.saveError != nil {
		return r.saveError
	}
	r.records[record.ScheduleID] = record
	return nil
}

type FakePoster struct {
	postCalled  bool
	postedContent string
	postError   error
}

func (p *FakePoster) Post(content string) error {
	p.postCalled = true
	p.postedContent = content
	return p.postError
}

func TestSchedulePostUseCase_Execute_WhenCanPost_PostsAndSavesRecord(t *testing.T) {
	clock := &FakeClock{fixedTime: time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{}
	schedule := domain.NewDailySchedule(12, 0)
	useCase := usecases.NewSchedulePostUseCase(clock, repo, poster, domain.NewPostGuard())

	err := useCase.Execute("test-schedule", schedule, "Hello World")

	require.NoError(t, err)
	assert.True(t, poster.postCalled)
	assert.Equal(t, "Hello World", poster.postedContent)
	assert.True(t, repo.saveCalled)
	assert.Equal(t, "test-schedule", repo.savedRecord.ScheduleID)
}

func TestSchedulePostUseCase_Execute_WhenAlreadyPostedToday_SkipsPost(t *testing.T) {
	clock := &FakeClock{fixedTime: time.Date(2026, 2, 1, 13, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	repo.records["test-schedule"] = domain.NewPostRecord("test-schedule", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC))
	poster := &FakePoster{}
	schedule := domain.NewDailySchedule(12, 0)
	useCase := usecases.NewSchedulePostUseCase(clock, repo, poster, domain.NewPostGuard())

	err := useCase.Execute("test-schedule", schedule, "Hello World")

	require.NoError(t, err)
	assert.False(t, poster.postCalled)
	assert.False(t, repo.saveCalled)
}

func TestSchedulePostUseCase_Execute_WhenPosterFails_ReturnsError(t *testing.T) {
	clock := &FakeClock{fixedTime: time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{postError: errors.New("network error")}
	schedule := domain.NewDailySchedule(12, 0)
	useCase := usecases.NewSchedulePostUseCase(clock, repo, poster, domain.NewPostGuard())

	err := useCase.Execute("test-schedule", schedule, "Hello World")

	require.Error(t, err)
	assert.False(t, repo.saveCalled)
}

func TestSchedulePostUseCase_Execute_WhenRepoSaveFails_ReturnsError(t *testing.T) {
	clock := &FakeClock{fixedTime: time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	repo.saveError = errors.New("disk full")
	poster := &FakePoster{}
	schedule := domain.NewDailySchedule(12, 0)
	useCase := usecases.NewSchedulePostUseCase(clock, repo, poster, domain.NewPostGuard())

	err := useCase.Execute("test-schedule", schedule, "Hello World")

	require.Error(t, err)
	assert.True(t, poster.postCalled)
}

func TestSchedulePostUseCase_ShouldExecuteNow_WhenTimeMatches_ReturnsTrue(t *testing.T) {
	clock := &FakeClock{fixedTime: time.Date(2026, 2, 1, 12, 0, 30, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{}
	schedule := domain.NewDailySchedule(12, 0)
	useCase := usecases.NewSchedulePostUseCase(clock, repo, poster, domain.NewPostGuard())

	shouldExecute := useCase.ShouldExecuteNow("test-schedule", schedule, time.Minute)

	assert.True(t, shouldExecute)
}

func TestSchedulePostUseCase_ShouldExecuteNow_WhenTimeNotMatches_ReturnsFalse(t *testing.T) {
	clock := &FakeClock{fixedTime: time.Date(2026, 2, 1, 11, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{}
	schedule := domain.NewDailySchedule(12, 0)
	useCase := usecases.NewSchedulePostUseCase(clock, repo, poster, domain.NewPostGuard())

	shouldExecute := useCase.ShouldExecuteNow("test-schedule", schedule, time.Minute)

	assert.False(t, shouldExecute)
}

func TestSchedulePostUseCase_ShouldExecuteNow_WhenAlreadyPosted_ReturnsFalse(t *testing.T) {
	clock := &FakeClock{fixedTime: time.Date(2026, 2, 1, 12, 0, 30, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	repo.records["test-schedule"] = domain.NewPostRecord("test-schedule", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC))
	poster := &FakePoster{}
	schedule := domain.NewDailySchedule(12, 0)
	useCase := usecases.NewSchedulePostUseCase(clock, repo, poster, domain.NewPostGuard())

	shouldExecute := useCase.ShouldExecuteNow("test-schedule", schedule, time.Minute)

	assert.False(t, shouldExecute)
}
