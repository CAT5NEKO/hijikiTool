package scheduler_test

import (
	"sync"
	"testing"
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/application/ports"
	"github.com/CAT5NEKO/hijikiTool/internal/domain"
	"github.com/CAT5NEKO/hijikiTool/internal/scheduler"
	"github.com/stretchr/testify/assert"
)

type FakeClock struct {
	currentTime time.Time
	mutex       sync.RWMutex
}

func (c *FakeClock) Now() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.currentTime
}

func (c *FakeClock) Advance(d time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.currentTime = c.currentTime.Add(d)
}

type FakePostRecordRepository struct {
	records map[string]domain.PostRecord
	mutex   sync.RWMutex
}

func NewFakePostRecordRepository() *FakePostRecordRepository {
	return &FakePostRecordRepository{records: make(map[string]domain.PostRecord)}
}

func (r *FakePostRecordRepository) Find(scheduleID string) (domain.PostRecord, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	record, exists := r.records[scheduleID]
	if !exists {
		return domain.PostRecord{}, nil
	}
	return record, nil
}

func (r *FakePostRecordRepository) Save(record domain.PostRecord) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.records[record.ScheduleID] = record
	return nil
}

type FakePoster struct {
	postCount int
	mutex     sync.Mutex
}

func (p *FakePoster) Post(content string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.postCount++
	return nil
}

func (p *FakePoster) GetPostCount() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.postCount
}

func TestScheduler_RunOnce_ExecutesScheduledPost(t *testing.T) {
	clock := &FakeClock{currentTime: time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{}

	job := scheduler.Job{
		ID:       "daily-post",
		Schedule: domain.NewDailySchedule(12, 0),
		Content:  "Test post",
	}

	s := scheduler.New(clock, repo, poster, []scheduler.Job{job})
	s.RunOnce()

	assert.Equal(t, 1, poster.GetPostCount())
}

func TestScheduler_RunOnce_DoesNotExecuteWhenNotTime(t *testing.T) {
	clock := &FakeClock{currentTime: time.Date(2026, 2, 1, 11, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{}

	job := scheduler.Job{
		ID:       "daily-post",
		Schedule: domain.NewDailySchedule(12, 0),
		Content:  "Test post",
	}

	s := scheduler.New(clock, repo, poster, []scheduler.Job{job})
	s.RunOnce()

	assert.Equal(t, 0, poster.GetPostCount())
}

func TestScheduler_RunOnce_DoesNotPostTwiceInSamePeriod(t *testing.T) {
	clock := &FakeClock{currentTime: time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{}

	job := scheduler.Job{
		ID:       "daily-post",
		Schedule: domain.NewDailySchedule(12, 0),
		Content:  "Test post",
	}

	s := scheduler.New(clock, repo, poster, []scheduler.Job{job})
	s.RunOnce()
	s.RunOnce()

	assert.Equal(t, 1, poster.GetPostCount())
}

func TestScheduler_NextWakeUpDuration_ReturnsOptimalSleepTime(t *testing.T) {
	clock := &FakeClock{currentTime: time.Date(2026, 2, 1, 11, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{}

	job := scheduler.Job{
		ID:       "daily-post",
		Schedule: domain.NewDailySchedule(12, 0),
		Content:  "Test post",
	}

	s := scheduler.New(clock, repo, poster, []scheduler.Job{job})
	duration := s.NextWakeUpDuration()

	assert.Equal(t, time.Hour, duration)
}

func TestScheduler_NextWakeUpDuration_WithMultipleJobs_ReturnsEarliest(t *testing.T) {
	clock := &FakeClock{currentTime: time.Date(2026, 2, 1, 11, 0, 0, 0, time.UTC)}
	repo := NewFakePostRecordRepository()
	poster := &FakePoster{}

	jobs := []scheduler.Job{
		{ID: "job1", Schedule: domain.NewDailySchedule(14, 0), Content: "Later post"},
		{ID: "job2", Schedule: domain.NewDailySchedule(11, 30), Content: "Sooner post"},
	}

	s := scheduler.New(clock, repo, poster, jobs)
	duration := s.NextWakeUpDuration()

	assert.Equal(t, 30*time.Minute, duration)
}

func TestScheduler_Config_FromConfigSource(t *testing.T) {
	config := ports.Config{
		MisskeyHost:  "example.com",
		MisskeyToken: "test-token",
		Visibility:   "home",
	}

	assert.Equal(t, "example.com", config.MisskeyHost)
	assert.Equal(t, "test-token", config.MisskeyToken)
}
