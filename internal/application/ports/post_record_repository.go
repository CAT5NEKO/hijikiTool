package ports

import "github.com/CAT5NEKO/hijikiTool/internal/domain"

type PostRecordRepository interface {
	Find(scheduleID string) (domain.PostRecord, error)
	Save(record domain.PostRecord) error
}
