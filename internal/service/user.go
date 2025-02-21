package service

import (
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository/sqlite"
)

type UserService interface {
	Register(username string, password string) error
	Validate(username string, password string) (int, error)
	GetHabits(userID int) ([]*models.Habit, error)
	// TODO: добавить удаление пользователей
}

type UserServiceImpl struct {
	repo *sqlite.UserRepository
}

func NewUserService(repo *sqlite.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Register(username string, password string) error {
	return s.repo.Create(username, password)
}

func (s *UserServiceImpl) Validate(username string, password string) (int, error) {
	return s.repo.ValidatePassword(username, password)
}

func (s *UserServiceImpl) GetHabits(userID int) ([]*models.Habit, error) {
	return s.repo.GetHabits(userID)
}
