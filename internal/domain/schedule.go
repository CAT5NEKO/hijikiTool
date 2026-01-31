package domain

import "time"

type PeriodType int

const (
	PeriodDaily PeriodType = iota
	PeriodWeekly
	PeriodMonthly
	PeriodYearly
)

type Schedule interface {
	NextTime(now time.Time) time.Time
	DurationUntil(now time.Time) time.Duration
	Period() PeriodType
}

type DailySchedule struct {
	hour   int
	minute int
}

func NewDailySchedule(hour, minute int) *DailySchedule {
	return &DailySchedule{hour: hour, minute: minute}
}

func (s *DailySchedule) NextTime(now time.Time) time.Time {
	candidate := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, 0, 0, now.Location())
	if now.Before(candidate) {
		return candidate
	}
	return candidate.AddDate(0, 0, 1)
}

func (s *DailySchedule) DurationUntil(now time.Time) time.Duration {
	return s.NextTime(now).Sub(now)
}

func (s *DailySchedule) Period() PeriodType {
	return PeriodDaily
}

type WeeklySchedule struct {
	dayOfWeek time.Weekday
	hour      int
	minute    int
}

func NewWeeklySchedule(dayOfWeek time.Weekday, hour, minute int) *WeeklySchedule {
	return &WeeklySchedule{dayOfWeek: dayOfWeek, hour: hour, minute: minute}
}

func (s *WeeklySchedule) NextTime(now time.Time) time.Time {
	daysUntilTarget := s.daysUntilTargetWeekday(now.Weekday())
	candidate := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, 0, 0, now.Location())
	candidate = candidate.AddDate(0, 0, daysUntilTarget)

	if daysUntilTarget == 0 && !now.Before(candidate) {
		return candidate.AddDate(0, 0, 7)
	}
	return candidate
}

func (s *WeeklySchedule) daysUntilTargetWeekday(current time.Weekday) int {
	days := int(s.dayOfWeek) - int(current)
	if days < 0 {
		days += 7
	}
	return days
}

func (s *WeeklySchedule) DurationUntil(now time.Time) time.Duration {
	return s.NextTime(now).Sub(now)
}

func (s *WeeklySchedule) Period() PeriodType {
	return PeriodWeekly
}

type MonthlySchedule struct {
	dayOfMonth int
	hour       int
	minute     int
}

func NewMonthlySchedule(dayOfMonth, hour, minute int) *MonthlySchedule {
	return &MonthlySchedule{dayOfMonth: dayOfMonth, hour: hour, minute: minute}
}

func (s *MonthlySchedule) NextTime(now time.Time) time.Time {
	candidate := s.clampedDateInMonth(now.Year(), now.Month(), now.Location())
	if now.Before(candidate) {
		return candidate
	}
	return s.clampedDateInMonth(now.Year(), now.Month()+1, now.Location())
}

func (s *MonthlySchedule) clampedDateInMonth(year int, month time.Month, loc *time.Location) time.Time {
	if month > 12 {
		year++
		month = 1
	}
	lastDayOfMonth := s.lastDayOfMonth(year, month)
	day := s.dayOfMonth
	if day > lastDayOfMonth {
		day = lastDayOfMonth
	}
	return time.Date(year, month, day, s.hour, s.minute, 0, 0, loc)
}

func (s *MonthlySchedule) lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func (s *MonthlySchedule) DurationUntil(now time.Time) time.Duration {
	return s.NextTime(now).Sub(now)
}

func (s *MonthlySchedule) Period() PeriodType {
	return PeriodMonthly
}

type YearlySchedule struct {
	month  time.Month
	day    int
	hour   int
	minute int
}

func NewYearlySchedule(month time.Month, day, hour, minute int) *YearlySchedule {
	return &YearlySchedule{month: month, day: day, hour: hour, minute: minute}
}

func (s *YearlySchedule) NextTime(now time.Time) time.Time {
	candidate := s.clampedDateInYear(now.Year(), now.Location())
	if now.Before(candidate) {
		return candidate
	}
	return s.clampedDateInYear(now.Year()+1, now.Location())
}

func (s *YearlySchedule) clampedDateInYear(year int, loc *time.Location) time.Time {
	lastDayOfMonth := time.Date(year, s.month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	day := s.day
	if day > lastDayOfMonth {
		day = lastDayOfMonth
	}
	return time.Date(year, s.month, day, s.hour, s.minute, 0, 0, loc)
}

func (s *YearlySchedule) DurationUntil(now time.Time) time.Duration {
	return s.NextTime(now).Sub(now)
}

func (s *YearlySchedule) Period() PeriodType {
	return PeriodYearly
}
