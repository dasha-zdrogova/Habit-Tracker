package sqlite

import (
	"database/sql"
	"fmt"
	"habit-tracker/internal/models"
	"habit-tracker/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	createUsers := `
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL
	)`

	createHabits := `CREATE TABLE IF NOT EXISTS habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(50) NOT NULL,
		description TEXT,
		UNIQUE (user_id, name)
	)`

	createHabitLogs := `CREATE TABLE IF NOT EXISTS habit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		habit_id INTEGER REFERENCES habits(habit_id) ON DELETE CASCADE,
		completed_date DATE,
		UNIQUE (habit_id, completed_date)
	)`

	_, err = db.Exec(createUsers)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	_, err = db.Exec(createHabits)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	_, err = db.Exec(createHabitLogs)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return &Storage{db: db}, nil
}

// TODO: разделить на папки users, habits, habit_logs и добавить новые методы
func (s *Storage) CreateUser(user *models.User) error {
	const op = "storage.sqlite.CreateUser"
	createUser := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`

	_, err := s.db.Exec(createUser, user.Username, user.PasswordHash)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s, %w", op, storage.ErrUserExists)
		}
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (s *Storage) GetUserHabits(user *models.User) ([]*models.Habit, error) {
	const op = "storage.sqlite.GetUserHabits"
	getHabits := `SELECT * FROM habits WHERE user_id = $1`

	notes, err := s.db.Query(getHabits, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer notes.Close()

	habits, err := s.scanHabits(notes)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if len(habits) == 0 {
		return nil, storage.ErrHabitsNotFound
	}
	return habits, nil
}

func (s *Storage) CreateHabit(habit *models.Habit) error {
	const op = "storage.sqlite.CreateHabit"
	createHabit := `INSERT INTO habits (user_id, name, description) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(createHabit, habit.UserID, habit.Name, habit.Description)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s, %w", op, storage.ErrHabitExists)
		}
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (s *Storage) MarkHabit(habitLogs *models.HabitLogs) error {
	const op = "storage.sqlite.MarkHabit"
	markHabit := `INSERT INTO habit_logs (habit_id, completed_date) VALUES ($1, $2)`

	_, err := s.db.Exec(markHabit, habitLogs.HabitID, habitLogs.CompletedDate)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s, %w", op, storage.ErrHabitMarked)
		}
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (s *Storage) GetHabitInfo(habit *models.Habit) ([]*models.HabitLogs, error) {
	const op = "storage.sqlite.GetHabitInfo"
	getHabit := `SELECT * FROM habit_logs WHERE habit_id = $1`

	notes, err := s.db.Query(getHabit, habit.ID)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer notes.Close()

	habits, err := s.scanHabit(notes)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if len(habits) == 0 {
		return nil, storage.ErrHabitNotFound
	}
	return habits, nil
}

func (s *Storage) DeleteHabit(habit *models.Habit) error {
	const op = "storage.sqlite.DeleteHabit"
	deleteHabit := `DELETE FROM habit WHERE id = $1`

	_, err := s.db.Exec(deleteHabit, habit.ID)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (s *Storage) scanHabit(rows *sql.Rows) ([]*models.HabitLogs, error) {
	const op = "storage.sqlite.scanHabit"

	var notes []*models.HabitLogs
	for rows.Next() {
		var note models.HabitLogs
		err := rows.Scan(&note.ID, &note.HabitID, &note.CompletedDate)

		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return notes, nil
}

// TODO: попробовать переписать на дженерики
func (s *Storage) scanHabits(rows *sql.Rows) ([]*models.Habit, error) {
	const op = "storage.sqlite.scanHabits"

	var notes []*models.Habit
	for rows.Next() {
		var note models.Habit
		err := rows.Scan(&note.ID, &note.UserID, &note.Name, &note.Description)

		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return notes, nil
}
