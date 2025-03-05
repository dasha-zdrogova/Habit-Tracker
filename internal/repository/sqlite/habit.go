package sqlite

import (
	"database/sql"
	"fmt"
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"
	"time"

	"github.com/mattn/go-sqlite3"
)

type SqliteHabitRepository struct {
	db *sql.DB
}

func NewSqliteHabitRepository(db *sql.DB) *SqliteHabitRepository {
	return &SqliteHabitRepository{db: db}
}

func (r *SqliteHabitRepository) Create(userID int, name string, description string) error {
	const op = "storage.sqlite.CreateHabit"
	createHabit := `INSERT INTO habits (user_id, name, description) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(createHabit, userID, name, description)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s, %w", op, repository.ErrHabitExists)
		}
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

// TODO: добавить обработку ошибки, когда не существует привычки с нужным ID
func (r *SqliteHabitRepository) Mark(habitID int, completedDate time.Time) error {
	const op = "storage.sqlite.MarkHabit"
	markHabit := `INSERT INTO habit_logs (habit_id, completed_date) VALUES ($1, $2)`

	_, err := r.db.Exec(markHabit, habitID, completedDate)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s, %w", op, repository.ErrHabitMarked)
		}
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *SqliteHabitRepository) GetInfo(habitID int) ([]*models.HabitLogs, error) {
	const op = "storage.sqlite.GetHabitInfo"
	getHabit := `SELECT * FROM habit_logs WHERE habit_id = $1`

	notes, err := r.db.Query(getHabit, habitID)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer notes.Close()

	habits, err := r.scanHabit(notes)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if len(habits) == 0 {
		return nil, repository.ErrHabitNotFound
	}
	return habits, nil
}

func (r *SqliteHabitRepository) Delete(habitID int) error {
	const op = "storage.sqlite.DeleteHabit"
	deleteHabit := `DELETE FROM habits WHERE id = $1`

	_, err := r.db.Exec(deleteHabit, habitID)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *SqliteHabitRepository) BelongsToUser(habitID int, userID int) error {
	const op = "storage.sqlite.BelongsToUser"
	isBelong := `SELECT COUNT(*) FROM habits WHERE id = $1 AND user_id = $2`

	var count int
	err := r.db.QueryRow(isBelong, habitID, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	if count == 0 {
		return fmt.Errorf("%s, %w", op, repository.ErrHabitNotBelongToUser)
	}
	return nil
}

func (r *SqliteHabitRepository) scanHabit(rows *sql.Rows) ([]*models.HabitLogs, error) {
	const op = "storage.sqlite.scanHabit"

	var notes []*models.HabitLogs
	for rows.Next() {
		var note models.HabitLogs
		err := rows.Scan(&note.ID, &note.HabitID, &note.CompletedAt)

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
