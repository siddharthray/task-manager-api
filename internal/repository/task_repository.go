package repository

import (
	"database/sql"
	"errors"
	"github.com/siddharthray/task-manager-api/internal/model"
	"log"
)

type TaskRepository interface {
	GetAll() ([]model.Task, error)
	GetByID(id int64) (*model.Task, error)
	Create(t *model.Task) (int64, error)
	Update(t *model.Task) error
	Delete(id int64) error
}

type mysqlTaskRepo struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &mysqlTaskRepo{DB: db}
}

// GetAll implements TaskRepository
func (r *mysqlTaskRepo) GetAll() ([]model.Task, error) {
	rows, err := r.DB.Query(`SELECT id, text, completed, created_at, completed_at, reopened_at, user_id FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("warning: rows.Close() failed: %v", closeErr)
		}
	}()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(
			&t.ID, &t.Text, &t.Completed,
			&t.CreatedAt, &t.CompletedAt, &t.ReopenedAt,
			&t.UserID,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// GetByID implements TaskRepository
func (r *mysqlTaskRepo) GetByID(id int64) (*model.Task, error) {
	var t model.Task
	err := r.DB.QueryRow(
		`SELECT id, text, completed, created_at, completed_at, reopened_at, user_id
         FROM tasks WHERE id = ?`,
		id,
	).Scan(
		&t.ID, &t.Text, &t.Completed,
		&t.CreatedAt, &t.CompletedAt, &t.ReopenedAt,
		&t.UserID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Create implements TaskRepository
func (r *mysqlTaskRepo) Create(t *model.Task) (int64, error) {
	res, err := r.DB.Exec(
		`INSERT INTO tasks (text, completed, created_at, user_id)
         VALUES (?, ?, ?, ?)`,
		t.Text, t.Completed, t.CreatedAt, t.UserID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update implements TaskRepository
func (r *mysqlTaskRepo) Update(t *model.Task) error {
	_, err := r.DB.Exec(
		`UPDATE tasks
         SET text = ?, completed = ?, completed_at = ?, reopened_at = ?
         WHERE id = ?`,
		t.Text, t.Completed, t.CompletedAt, t.ReopenedAt, t.ID,
	)
	return err
}

// Delete implements TaskRepository
func (r *mysqlTaskRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}
