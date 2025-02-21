package sqlite

import (
	"database/sql"
	"fmt"
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"

	"github.com/mattn/go-sqlite3"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(username string, password string) error {
	const op = "repository.sqlite.user.CreateUser"
	createUser := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`

	_, err := r.db.Exec(createUser, username, password)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s, %w", op, repository.ErrUserExists)
		}
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *UserRepository) ValidatePassword(username string, password string) (int, error) {
	const op = "repository.sqlite.user. ValidatePassword"
	//TODO: password_hash -> password
	//TODO: придумать что-то с хэшированием пароля
	//TODO: улучшить код обработки ошибок
	validatePassword := `SELECT id FROM users WHERE username = $1 AND password_hash = $2`

	var userID int
	err := r.db.QueryRow(validatePassword, username, password).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("%s, %w", op, repository.ErrUserNotFound)
		}
		return -1, fmt.Errorf("%s, %w", op, err)
	}
	return userID, nil
}

func (r *UserRepository) GetHabits(userID int) ([]*models.Habit, error) {
	const op = "repository.sqlite.GetUserHabits"
	getHabits := `SELECT * FROM habits WHERE user_id = $1`

	notes, err := r.db.Query(getHabits, userID)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer notes.Close()

	habits, err := r.scanHabits(notes)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if len(habits) == 0 {
		return nil, repository.ErrHabitsNotFound
	}
	return habits, nil
}

// func (r *UserRepository) GetByUsername(username string) (int, error) {
// 	const op = "repository.sqlite.GetByUsername"
// 	getByUsername := `SELECT * FROM users WHERE username = $1`

// 	var userID int
// 	err := r.db.QueryRow(getByUsername, username).Scan(&userID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return -1, fmt.Errorf("%s, %w", op, repository.ErrUserNotFound)
// 		}
// 		return -1, fmt.Errorf("%s, %w", op, err)
// 	}
// 	return userID, nil
// }

func (r *UserRepository) scanHabits(rows *sql.Rows) ([]*models.Habit, error) {
	const op = "repository.sqlite.scanHabits"

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
