package storage

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserExists     = errors.New("user exists")
	ErrHabitsNotFound = errors.New("habits not found")
	ErrHabitNotFound  = errors.New("habit not found")
	ErrHabitExists    = errors.New("habit already exists")
	ErrHabitMarked    = errors.New("habit already marked")
)
