package infrastructure_test

import (
	"os"
	"testing"
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/domain"
	"github.com/CAT5NEKO/hijikiTool/internal/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduleConfigLoader_Load_DailySchedule(t *testing.T) {
	configJSON := `{
		"schedules": [
			{
				"id": "morning-post",
				"type": "daily",
				"hour": 8,
				"minute": 30,
				"content": "おはようございます！"
			}
		]
	}`
	filePath := createTempConfigFile(t, configJSON)
	defer os.Remove(filePath)

	loader := infrastructure.NewScheduleConfigLoader(filePath)
	configs, err := loader.Load()

	require.NoError(t, err)
	require.Len(t, configs, 1)
	assert.Equal(t, "morning-post", configs[0].ID)
	assert.Equal(t, "おはようございます！", configs[0].Content)
	assert.Equal(t, domain.PeriodDaily, configs[0].Schedule.Period())
}

func TestScheduleConfigLoader_Load_WeeklySchedule(t *testing.T) {
	configJSON := `{
		"schedules": [
			{
				"id": "weekly-report",
				"type": "weekly",
				"dayOfWeek": 1,
				"hour": 9,
				"minute": 0,
				"content": "週報です"
			}
		]
	}`
	filePath := createTempConfigFile(t, configJSON)
	defer os.Remove(filePath)

	loader := infrastructure.NewScheduleConfigLoader(filePath)
	configs, err := loader.Load()

	require.NoError(t, err)
	require.Len(t, configs, 1)
	assert.Equal(t, domain.PeriodWeekly, configs[0].Schedule.Period())
}

func TestScheduleConfigLoader_Load_MonthlySchedule(t *testing.T) {
	configJSON := `{
		"schedules": [
			{
				"id": "monthly-post",
				"type": "monthly",
				"dayOfMonth": 1,
				"hour": 12,
				"minute": 0,
				"content": "月初投稿"
			}
		]
	}`
	filePath := createTempConfigFile(t, configJSON)
	defer os.Remove(filePath)

	loader := infrastructure.NewScheduleConfigLoader(filePath)
	configs, err := loader.Load()

	require.NoError(t, err)
	require.Len(t, configs, 1)
	assert.Equal(t, domain.PeriodMonthly, configs[0].Schedule.Period())
}

func TestScheduleConfigLoader_Load_YearlySchedule(t *testing.T) {
	configJSON := `{
		"schedules": [
			{
				"id": "new-year",
				"type": "yearly",
				"month": 1,
				"dayOfMonth": 1,
				"hour": 0,
				"minute": 0,
				"content": "あけましておめでとうございます！"
			}
		]
	}`
	filePath := createTempConfigFile(t, configJSON)
	defer os.Remove(filePath)

	loader := infrastructure.NewScheduleConfigLoader(filePath)
	configs, err := loader.Load()

	require.NoError(t, err)
	require.Len(t, configs, 1)
	assert.Equal(t, domain.PeriodYearly, configs[0].Schedule.Period())
}

func TestScheduleConfigLoader_Load_MultipleSchedules(t *testing.T) {
	configJSON := `{
		"schedules": [
			{
				"id": "morning",
				"type": "daily",
				"hour": 8,
				"minute": 0,
				"content": "おはよう"
			},
			{
				"id": "night",
				"type": "daily",
				"hour": 22,
				"minute": 0,
				"content": "おやすみ"
			}
		]
	}`
	filePath := createTempConfigFile(t, configJSON)
	defer os.Remove(filePath)

	loader := infrastructure.NewScheduleConfigLoader(filePath)
	configs, err := loader.Load()

	require.NoError(t, err)
	assert.Len(t, configs, 2)
}

func TestScheduleConfigLoader_Load_InvalidType_ReturnsError(t *testing.T) {
	configJSON := `{
		"schedules": [
			{
				"id": "invalid",
				"type": "hourly",
				"hour": 8,
				"minute": 0,
				"content": "test"
			}
		]
	}`
	filePath := createTempConfigFile(t, configJSON)
	defer os.Remove(filePath)

	loader := infrastructure.NewScheduleConfigLoader(filePath)
	_, err := loader.Load()

	require.Error(t, err)
}

func TestScheduleConfigLoader_Load_FileNotFound_ReturnsError(t *testing.T) {
	loader := infrastructure.NewScheduleConfigLoader("nonexistent.json")
	_, err := loader.Load()

	require.Error(t, err)
}

func TestScheduleConfigLoader_Load_InvalidJSON_ReturnsError(t *testing.T) {
	filePath := createTempConfigFile(t, "{ invalid json }")
	defer os.Remove(filePath)

	loader := infrastructure.NewScheduleConfigLoader(filePath)
	_, err := loader.Load()

	require.Error(t, err)
}

func TestScheduleConfigLoader_Load_DailySchedule_VerifyNextTime(t *testing.T) {
	configJSON := `{
		"schedules": [
			{
				"id": "test",
				"type": "daily",
				"hour": 15,
				"minute": 30,
				"content": "test"
			}
		]
	}`
	filePath := createTempConfigFile(t, configJSON)
	defer os.Remove(filePath)

	loader := infrastructure.NewScheduleConfigLoader(filePath)
	configs, err := loader.Load()

	require.NoError(t, err)
	now := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	nextTime := configs[0].Schedule.NextTime(now)

	assert.Equal(t, time.Date(2026, 2, 1, 15, 30, 0, 0, time.UTC), nextTime)
}

func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	file, err := os.CreateTemp("", "config-*.json")
	require.NoError(t, err)
	_, err = file.WriteString(content)
	require.NoError(t, err)
	file.Close()
	return file.Name()
}
