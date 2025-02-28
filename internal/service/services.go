package service

import (
	"habit-tracker/internal/repository"
)

type Services struct {
	Users  UserService
	Habits HabitService
}

func NewServices(usersRepo repository.UserRepository, habitsRepo repository.HabitRepository) *Services {
	return &Services{
		Users:  NewUserService(usersRepo),
		Habits: NewHabitService(habitsRepo),
	}
}
