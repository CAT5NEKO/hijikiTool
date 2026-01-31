package domain_test

import (
	"testing"
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDailySchedule_NextTime_WhenBeforeScheduledHour_ReturnsSameDay(t *testing.T) {
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 2, 1, 11, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestDailySchedule_NextTime_WhenAfterScheduledHour_ReturnsNextDay(t *testing.T) {
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 2, 1, 13, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 2, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestDailySchedule_NextTime_WhenExactlyAtScheduledTime_ReturnsNextDay(t *testing.T) {
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 2, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestDailySchedule_NextTime_CrossingMonthBoundary(t *testing.T) {
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 1, 31, 13, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestDailySchedule_NextTime_CrossingYearBoundary(t *testing.T) {
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2025, 12, 31, 13, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestWeeklySchedule_NextTime_WhenTargetDayIsLaterThisWeek(t *testing.T) {
	schedule := domain.NewWeeklySchedule(time.Wednesday, 10, 30)
	now := time.Date(2026, 2, 2, 9, 0, 0, 0, time.UTC) // Monday

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 4, 10, 30, 0, 0, time.UTC), nextTime) // Wednesday
}

func TestWeeklySchedule_NextTime_WhenTargetDayIsTomorrow(t *testing.T) {
	schedule := domain.NewWeeklySchedule(time.Monday, 9, 0)
	now := time.Date(2026, 2, 1, 18, 0, 0, 0, time.UTC) // Sunday

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 2, 9, 0, 0, 0, time.UTC), nextTime) // Monday
}

func TestWeeklySchedule_NextTime_WhenSameDayButAfterTime_ReturnsNextWeek(t *testing.T) {
	schedule := domain.NewWeeklySchedule(time.Monday, 9, 0)
	now := time.Date(2026, 2, 2, 10, 0, 0, 0, time.UTC) // Monday 10:00

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 9, 9, 0, 0, 0, time.UTC), nextTime) // Next Monday
}

func TestMonthlySchedule_NextTime_WhenBeforeTargetDay(t *testing.T) {
	schedule := domain.NewMonthlySchedule(15, 8, 0)
	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 15, 8, 0, 0, 0, time.UTC), nextTime)
}

func TestMonthlySchedule_NextTime_WhenAfterTargetDay_ReturnsNextMonth(t *testing.T) {
	schedule := domain.NewMonthlySchedule(1, 12, 0)
	now := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestMonthlySchedule_NextTime_WhenTargetDayExceedsMonthDays_ClampsToLastDay(t *testing.T) {
	schedule := domain.NewMonthlySchedule(31, 12, 0)
	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 28, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestMonthlySchedule_NextTime_LeapYearFebruary29th(t *testing.T) {
	schedule := domain.NewMonthlySchedule(29, 12, 0)
	now := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC) // 2024 is leap year

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestYearlySchedule_NextTime_WhenBeforeTargetDate(t *testing.T) {
	schedule := domain.NewYearlySchedule(time.July, 4, 0, 0)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC), nextTime)
}

func TestYearlySchedule_NextTime_WhenAfterTargetDate_ReturnsNextYear(t *testing.T) {
	schedule := domain.NewYearlySchedule(time.January, 1, 0, 0)
	now := time.Date(2025, 12, 31, 23, 59, 0, 0, time.UTC)

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), nextTime)
}

func TestYearlySchedule_NextTime_February29OnNonLeapYear_ClampsTo28th(t *testing.T) {
	schedule := domain.NewYearlySchedule(time.February, 29, 12, 0)
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // 2025 is not leap year

	nextTime := schedule.NextTime(now)

	assert.Equal(t, time.Date(2025, 2, 28, 12, 0, 0, 0, time.UTC), nextTime)
}

func TestSchedule_Period_ReturnsCorrectPeriodType(t *testing.T) {
	dailySchedule := domain.NewDailySchedule(12, 0)
	weeklySchedule := domain.NewWeeklySchedule(time.Monday, 9, 0)
	monthlySchedule := domain.NewMonthlySchedule(1, 12, 0)
	yearlySchedule := domain.NewYearlySchedule(time.January, 1, 0, 0)

	assert.Equal(t, domain.PeriodDaily, dailySchedule.Period())
	assert.Equal(t, domain.PeriodWeekly, weeklySchedule.Period())
	assert.Equal(t, domain.PeriodMonthly, monthlySchedule.Period())
	assert.Equal(t, domain.PeriodYearly, yearlySchedule.Period())
}

func TestSchedule_DurationUntil_ReturnsPositiveDuration(t *testing.T) {
	schedule := domain.NewDailySchedule(12, 0)
	now := time.Date(2026, 2, 1, 11, 0, 0, 0, time.UTC)

	duration := schedule.DurationUntil(now)

	require.True(t, duration > 0)
	assert.Equal(t, time.Hour, duration)
}
