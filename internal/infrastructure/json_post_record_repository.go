package infrastructure

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/application/ports"
	"github.com/CAT5NEKO/hijikiTool/internal/domain"
)

type JSONPostRecordRepository struct {
	filePath string
	mutex    sync.RWMutex
}

type jsonRecordStore struct {
	Records map[string]jsonRecord `json:"records"`
}

type jsonRecord struct {
	ScheduleID   string    `json:"schedule_id"`
	LastPostedAt time.Time `json:"last_posted_at"`
}

func NewJSONPostRecordRepository(filePath string) ports.PostRecordRepository {
	return &JSONPostRecordRepository{filePath: filePath}
}

func (r *JSONPostRecordRepository) Find(scheduleID string) (domain.PostRecord, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	store, err := r.loadStore()
	if err != nil {
		if os.IsNotExist(err) {
			return domain.PostRecord{}, nil
		}
		return domain.PostRecord{}, err
	}

	record, exists := store.Records[scheduleID]
	if !exists {
		return domain.PostRecord{}, nil
	}

	return domain.NewPostRecord(record.ScheduleID, record.LastPostedAt), nil
}

func (r *JSONPostRecordRepository) Save(record domain.PostRecord) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	store, err := r.loadStore()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if store.Records == nil {
		store.Records = make(map[string]jsonRecord)
	}

	store.Records[record.ScheduleID] = jsonRecord{
		ScheduleID:   record.ScheduleID,
		LastPostedAt: record.LastPostedAt,
	}

	return r.saveStore(store)
}

func (r *JSONPostRecordRepository) loadStore() (jsonRecordStore, error) {
	var store jsonRecordStore

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return jsonRecordStore{Records: make(map[string]jsonRecord)}, err
	}

	if err := json.Unmarshal(data, &store); err != nil {
		return jsonRecordStore{Records: make(map[string]jsonRecord)}, err
	}

	return store, nil
}

func (r *JSONPostRecordRepository) saveStore(store jsonRecordStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}
