package main

import (
	"habit-tracker/internal/config"
	"habit-tracker/internal/handler"
	sl "habit-tracker/internal/lib/logger"
	"habit-tracker/internal/repository/sqlite"
	"habit-tracker/internal/service"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
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

	userRepo := sqlite.NewSqliteUserRepository(storage.DB)
	habitRepo := sqlite.NewSqliteHabitRepository(storage.DB)

	service := service.NewServices(userRepo, habitRepo)
	// регистрация
	// err = service.Users.Register("dasha", "d")
	// if err != nil {
	// 	logger.Error("failed to create user", sl.Err(err))
	// 	os.Exit(1)
	// }

	// // авторизация
	// userID, err := service.Users.Login("dasha", "d")
	// if err != nil {
	// 	logger.Error("failed to login", sl.Err(err))
	// }

	// // создание привычки
	// err = service.Habits.Create(userID, "dance", "")
	// if err != nil {
	// 	logger.Error("failed to create habit", sl.Err(err))
	// 	os.Exit(1)
	// }

	// // отметка привычки
	// err = service.Habits.Mark(1, time.Now().UTC().Truncate(24*time.Hour))
	// if err != nil {
	// 	logger.Error("failed to mark habit", sl.Err(err))
	// 	os.Exit(1)
	// }

	// // получение всех привычек
	// notes, err := service.Users.GetHabits(userID)
	// if err != nil {
	// 	logger.Error("failed get notes", sl.Err(err))
	// 	os.Exit(1)
	// }
	// for _, note := range notes {
	// 	fmt.Println(note)
	// }

	handlers := handler.NewHandler(service)
	router := chi.NewRouter()
	handlers.Register(router)

	if err := http.ListenAndServe(config.Address, router); err != nil {
		logger.Error("server stoped", sl.Err(err))
	}
	// router.Use(middleware.RequestID)
	// //TODO: сделать middleware на логирование + добавить локальное логирование
	// router.Use(middleware.Recoverer)
	// router.Use(middleware.URLFormat)
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
