package sqlite

import "database/sql"

type Repositories struct {
	Users *UserRepository
	Habits *HabitRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users:  NewUserRepository(db),
		Habits: NewHabitRepository(db),
	}
}