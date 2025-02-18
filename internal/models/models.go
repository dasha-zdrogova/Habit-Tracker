package models

import "time"

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type Habit struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

type HabitLogs struct {
	ID int `json:"id"`
	HabitID int `json:"habit_id"`
	CompletedDate time.Time `json:"completed_date"`
}