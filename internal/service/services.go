package service

import "habit-tracker/internal/repository/sqlite"

type Services struct {
	Users UserService
	Habits HabitService
}

func NewServices(repos *sqlite.Repositories) *Services {
	return &Services{
		Users:  NewUserService(repos.Users),
		Habits: NewHabitService(repos.Habits),
	}
}