package service

import (
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"
)

type UserService interface {
	Register(username string, password string) error
	Login(username string, password string) (int, error)
	GetHabits(userID int) ([]*models.Habit, error)
	// TODO: добавить удаление пользователей
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Register(username string, password string) error {
	return s.repo.Create(username, password)
}

func (s *UserServiceImpl) Login(username string, password string) (int, error) {
	return s.repo.Login(username, password)
}

func (s *UserServiceImpl) GetHabits(userID int) ([]*models.Habit, error) {
	return s.repo.GetHabits(userID)
}
