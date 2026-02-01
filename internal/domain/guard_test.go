package domain_test

import (
	"testing"
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPostGuard_CanPost_WithNoRecord_ReturnsTrue(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)
	emptyRecord := domain.PostRecord{}

	canPost := guard.CanPost(schedule, emptyRecord, now)

	assert.True(t, canPost)
}

func TestPostGuard_CanPost_DailySchedule_SameDay_ReturnsFalse(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 2, 1, 13, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("daily-12", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.False(t, canPost)
}

func TestPostGuard_CanPost_DailySchedule_NextDay_ReturnsTrue(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 2, 2, 12, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("daily-12", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.True(t, canPost)
}

func TestPostGuard_CanPost_WeeklySchedule_SameWeek_ReturnsFalse(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewWeeklySchedule(time.Monday, 9, 0)
	now := time.Date(2026, 2, 4, 10, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("weekly-mon", time.Date(2026, 2, 2, 9, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.False(t, canPost)
}

func TestPostGuard_CanPost_WeeklySchedule_NextWeek_ReturnsTrue(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewWeeklySchedule(time.Monday, 9, 0)
	now := time.Date(2026, 2, 9, 9, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("weekly-mon", time.Date(2026, 2, 2, 9, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.True(t, canPost)
}

func TestPostGuard_CanPost_MonthlySchedule_SameMonth_ReturnsFalse(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewMonthlySchedule(1, 12, 0)
	now := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("monthly-1", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.False(t, canPost)
}

func TestPostGuard_CanPost_MonthlySchedule_NextMonth_ReturnsTrue(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewMonthlySchedule(1, 12, 0)
	now := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("monthly-1", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.True(t, canPost)
}

func TestPostGuard_CanPost_YearlySchedule_SameYear_ReturnsFalse(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewYearlySchedule(time.January, 1, 0, 0)
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("yearly-jan1", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.False(t, canPost)
}

func TestPostGuard_CanPost_YearlySchedule_NextYear_ReturnsTrue(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewYearlySchedule(time.January, 1, 0, 0)
	now := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("yearly-jan1", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.True(t, canPost)
}

func TestPostGuard_CanPost_DailySchedule_CrossingYearBoundary_ReturnsTrue(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewDailySchedule(0, 0)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("daily-0", time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.True(t, canPost)
}

func TestPostGuard_CanPost_WeeklySchedule_CrossingYearBoundary_ReturnsTrue(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewWeeklySchedule(time.Thursday, 9, 0)
	now := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("weekly-thu", time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.True(t, canPost)
}

func TestPostGuard_CanPost_AfterServerRestart_WithOldRecord_ReturnsFalse(t *testing.T) {
	guard := domain.NewPostGuard()
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 2, 1, 14, 0, 0, 0, time.UTC)
	record := domain.NewPostRecord("daily-12", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC))

	canPost := guard.CanPost(schedule, record, now)

	assert.False(t, canPost)
}
