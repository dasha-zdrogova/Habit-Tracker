package main

import (
	"fmt"
	"habit-tracker/internal/config"
	sl "habit-tracker/internal/lib/logger"
	"habit-tracker/internal/models"
	"habit-tracker/internal/storage/sqlite"
	"log/slog"
	"os"
	// "time"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	// config
	config := config.MustLoad()

	// logger
	logger := SetupLogger(config.Env)

	logger.Info("starting habit-tracker", slog.String("env", config.Env))
	logger.Debug("debug messages are enabled")

	// storage
	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	// err = storage.CreateUser(&models.User{
	// 	ID:           1,
	// 	Username:     "dasha",
	// 	PasswordHash: "d",
	// })
	// if err != nil {
	// 	logger.Error("failed to create user", sl.Err(err))
	// 	os.Exit(1)
	// }
	// err = storage.CreateHabit(&models.Habit{
	// 	ID:          1,
	// 	UserID:      1,
	// 	Name:        "drink",
	// 	Description: "",
	// })
	// if err != nil {
	// 	logger.Error("failed to create habit", sl.Err(err))
	// 	os.Exit(1)
	// }
	// err = storage.MarkHabit(&models.HabitLogs{
	// 	ID:            0,
	// 	HabitID:       1,
	// 	CompletedDate: time.Now().UTC().Truncate(24 * time.Hour),
	// })
	// if err != nil {
	// 	logger.Error("failed to mark habit", sl.Err(err))
	// 	os.Exit(1)
	// }
	notes, err := storage.GetUserHabits(&models.User{
		ID:           2,
		Username:     "d",
		PasswordHash: "",
	})
	if err != nil {
		logger.Error("failed get notes", sl.Err(err))
		os.Exit(1)
	}
	logger.Info("habits")
	for _, note := range notes {
		fmt.Println(note)
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
