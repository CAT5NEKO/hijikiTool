package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CAT5NEKO/hijikiTool/internal/infrastructure"
	"github.com/CAT5NEKO/hijikiTool/internal/scheduler"
)

func main() {
	logFile := setupLogger()
	defer logFile.Close()

	envConfigLoader := infrastructure.NewEnvConfigLoader(".env")
	envConfig, err := envConfigLoader.Load()
	if err != nil {
		log.Fatalf("Failed to load env config: %v", err)
	}

	scheduleConfigLoader := infrastructure.NewScheduleConfigLoader("config.json")
	scheduleConfigs, err := scheduleConfigLoader.Load()
	if err != nil {
		log.Fatalf("Failed to load schedule config: %v", err)
	}

	clock := infrastructure.NewRealClock()
	repository := infrastructure.NewJSONPostRecordRepository("post_records.json")
	poster := infrastructure.NewMisskeyPoster(envConfig)

	jobs := createJobsFromScheduleConfigs(scheduleConfigs)

	s := scheduler.New(clock, repository, poster, jobs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleShutdown(cancel)

	log.Printf("Scheduler started with %d job(s)", len(jobs))
	s.Run(ctx)
	log.Println("Scheduler stopped")
}

func setupLogger() *os.File {
	logFile, err := os.OpenFile("hijiki.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	log.SetOutput(logFile)
	return logFile
}

func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown signal received")
	cancel()
}

func createJobsFromScheduleConfigs(configs []infrastructure.ScheduleConfig) []scheduler.Job {
	jobs := make([]scheduler.Job, 0, len(configs))
	for _, config := range configs {
		jobs = append(jobs, scheduler.Job{
			ID:       config.ID,
			Schedule: config.Schedule,
			Content:  config.Content,
		})
	}
	return jobs
}
