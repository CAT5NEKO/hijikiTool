package infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/domain"
)

type ScheduleConfig struct {
	ID       string
	Schedule domain.Schedule
	Content  string
}

type ScheduleConfigLoader struct {
	filePath string
}

type scheduleConfigFile struct {
	Schedules []scheduleConfigEntry `json:"schedules"`
}

type scheduleConfigEntry struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
	DayOfWeek  int    `json:"dayOfWeek"`
	DayOfMonth int    `json:"dayOfMonth"`
	Month      int    `json:"month"`
	Content    string `json:"content"`
}

func NewScheduleConfigLoader(filePath string) *ScheduleConfigLoader {
	return &ScheduleConfigLoader{filePath: filePath}
}

func (l *ScheduleConfigLoader) Load() ([]ScheduleConfig, error) {
	data, err := os.ReadFile(l.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configFile scheduleConfigFile
	if err := json.Unmarshal(data, &configFile); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return l.convertToScheduleConfigs(configFile.Schedules)
}

func (l *ScheduleConfigLoader) convertToScheduleConfigs(entries []scheduleConfigEntry) ([]ScheduleConfig, error) {
	configs := make([]ScheduleConfig, 0, len(entries))

	for _, entry := range entries {
		schedule, err := l.createSchedule(entry)
		if err != nil {
			return nil, err
		}

		configs = append(configs, ScheduleConfig{
			ID:       entry.ID,
			Schedule: schedule,
			Content:  entry.Content,
		})
	}

	return configs, nil
}

func (l *ScheduleConfigLoader) createSchedule(entry scheduleConfigEntry) (domain.Schedule, error) {
	switch entry.Type {
	case "daily":
		return domain.NewDailySchedule(entry.Hour, entry.Minute), nil
	case "weekly":
		return domain.NewWeeklySchedule(time.Weekday(entry.DayOfWeek), entry.Hour, entry.Minute), nil
	case "monthly":
		return domain.NewMonthlySchedule(entry.DayOfMonth, entry.Hour, entry.Minute), nil
	case "yearly":
		return domain.NewYearlySchedule(time.Month(entry.Month), entry.DayOfMonth, entry.Hour, entry.Minute), nil
	default:
		return nil, fmt.Errorf("unknown schedule type: %s", entry.Type)
	}
}
