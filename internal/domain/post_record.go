package domain

import "time"

type PostRecord struct {
	ScheduleID   string
	LastPostedAt time.Time
}

func NewPostRecord(scheduleID string, lastPostedAt time.Time) PostRecord {
	return PostRecord{
		ScheduleID:   scheduleID,
		LastPostedAt: lastPostedAt,
	}
}

func (r PostRecord) IsZero() bool {
	return r.LastPostedAt.IsZero()
}
