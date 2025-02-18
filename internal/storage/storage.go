package storage

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user exists")
	ErrHabitExists        = errors.New("habit exists")
	ErrHabitAlreadyMarked = errors.New("habit already marked")
)
