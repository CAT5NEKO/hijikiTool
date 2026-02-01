package domain

import "time"

type PostGuard struct{}

func NewPostGuard() *PostGuard {
	return &PostGuard{}
}

func (g *PostGuard) CanPost(schedule Schedule, record PostRecord, now time.Time) bool {
	if record.IsZero() {
		return true
	}
	return g.isInNewPeriod(schedule.Period(), record.LastPostedAt, now)
}

func (g *PostGuard) isInNewPeriod(period PeriodType, lastPosted, now time.Time) bool {
	switch period {
	case PeriodDaily:
		return g.isDifferentDay(lastPosted, now)
	case PeriodWeekly:
		return g.isDifferentWeek(lastPosted, now)
	case PeriodMonthly:
		return g.isDifferentMonth(lastPosted, now)
	case PeriodYearly:
		return g.isDifferentYear(lastPosted, now)
	default:
		return true
	}
}

func (g *PostGuard) isDifferentDay(lastPosted, now time.Time) bool {
	lastYear, lastMonth, lastDay := lastPosted.Date()
	nowYear, nowMonth, nowDay := now.Date()
	return lastYear != nowYear || lastMonth != nowMonth || lastDay != nowDay
}

func (g *PostGuard) isDifferentWeek(lastPosted, now time.Time) bool {
	lastYear, lastWeek := lastPosted.ISOWeek()
	nowYear, nowWeek := now.ISOWeek()
	return lastYear != nowYear || lastWeek != nowWeek
}

func (g *PostGuard) isDifferentMonth(lastPosted, now time.Time) bool {
	lastYear, lastMonth, _ := lastPosted.Date()
	nowYear, nowMonth, _ := now.Date()
	return lastYear != nowYear || lastMonth != nowMonth
}

func (g *PostGuard) isDifferentYear(lastPosted, now time.Time) bool {
	return lastPosted.Year() != now.Year()
}
