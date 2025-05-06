package repository

import (
	"database/sql"
	"errors"
	"github.com/siddharthray/task-manager-api/internal/model"
)

type TaskRepository interface {
	GetAll() ([]model.Task, error)
	GetByID(id int64) (*model.Task, error)
	Create(t *model.Task) (int64, error)
	UpdateTask(t *model.Task) (*model.Task, error)
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
	rows, err := r.DB.Query(`
        SELECT
          id,
          text,
          completed,
          created_at,
          completed_at,
          reopened_at,
          updated_at
        FROM tasks
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(
			&t.ID,
			&t.Text,
			&t.Completed,
			&t.CreatedAt,
			&t.CompletedAt,
			&t.ReopenedAt,
			&t.UpdatedAt,
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
	err := r.DB.QueryRow(`
        SELECT
          id,
          text,
          completed,
          created_at,
          completed_at,
          reopened_at,
          updated_at
        FROM tasks
        WHERE id = ?
    `, id).Scan(
		&t.ID,
		&t.Text,
		&t.Completed,
		&t.CreatedAt,
		&t.CompletedAt,
		&t.ReopenedAt,
		&t.UpdatedAt,
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
	res, err := r.DB.Exec(`
        INSERT INTO tasks (text)
        VALUES (?)
    `, t.Text)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateTask Update implements TaskRepository
func (r *mysqlTaskRepo) UpdateTask(t *model.Task) (*model.Task, error) {
	_, err := r.DB.Exec(`
        UPDATE tasks
           SET text      = ?,
               completed = ?
         WHERE id = ?
    `, t.Text, t.Completed, t.ID)
	if err != nil {
		return nil, err
	}

	// re‑select the full row
	var out model.Task
	err = r.DB.QueryRow(`
        SELECT id, text, completed,
               created_at, completed_at, reopened_at, updated_at
          FROM tasks
         WHERE id = ?
    `, t.ID).Scan(
		&out.ID,
		&out.Text,
		&out.Completed,
		&out.CreatedAt,
		&out.CompletedAt,
		&out.ReopenedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete implements TaskRepository
func (r *mysqlTaskRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}
