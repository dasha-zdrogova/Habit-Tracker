package service

import (
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"
	"time"
)

type HabitService interface {
	Create(userID int, name string, description string) error
	Mark(habitID int, completedDate time.Time) error
	GetInfo(habitId int) ([]*models.HabitLogs, error)
	Delete(habitID int) error
	BelongsToUser(habitID int, userID int) error
}

type HabitServiceImpl struct {
	repo repository.HabitRepository
}

func NewHabitService(repo repository.HabitRepository) HabitService {
	return &HabitServiceImpl{repo: repo}
}

func (s *HabitServiceImpl) Create(userID int, name string, description string) error {
	return s.repo.Create(userID, name, description)
}

func (s *HabitServiceImpl) Mark(habitID int, completedDate time.Time) error {
	return s.repo.Mark(habitID, completedDate)
}

func (s *HabitServiceImpl) GetInfo(habitId int) ([]*models.HabitLogs, error) {
	return s.repo.GetInfo(habitId)
}

func (s *HabitServiceImpl) Delete(habitId int) error {
	return s.repo.Delete(habitId)
}

func (s *HabitServiceImpl) BelongsToUser(habitID int, userID int) error {
	return s.repo.BelongsToUser(habitID, userID)
}