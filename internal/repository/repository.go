package repository

import (
	"errors"
	"habit-tracker/internal/models"
	"time"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserExists     = errors.New("user exists")
	ErrHabitsNotFound = errors.New("habits not found")
	ErrHabitNotFound  = errors.New("habit not found")
	ErrHabitExists    = errors.New("habit already exists")
	ErrHabitMarked    = errors.New("habit already marked")
)

type UserRepository interface {
	Create(username string, password string) error
	Login(username string, password string) (int, error)
	GetHabits(userID int) ([]*models.Habit, error)
}

type HabitRepository interface {
	Create(userID int, name string, description string) error
	Mark(habitID int, completedDate time.Time) error
	GetInfo(habitID int) ([]*models.HabitLogs, error)
	Delete(habitID int) error
}
