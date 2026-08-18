// internal/repository/postgres/task_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/az/task-api/internal/domain"
	"github.com/az/task-api/internal/repository"
)

type taskRepository struct {
	db *sql.DB
}

// NewTaskRepository returns the interface type, not the struct. Callers
// (usecases) code against repository.TaskRepository, never this concrete type.
func NewTaskRepository(db *sql.DB) repository.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, t *domain.Task) error {
	query := `
		INSERT INTO tasks (title, description, status, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		t.Title, t.Description, t.Status, t.DueDate, t.CreatedAt, t.UpdatedAt,
	).Scan(&t.ID)
}

func (r *taskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, due_date, created_at, updated_at
		FROM tasks WHERE id = $1`

	var t domain.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrTaskNotFound // defined next in Step 10 (errors)
	}
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return &t, nil
}

func (r *taskRepository) List(ctx context.Context) ([]*domain.Task, error) {
	query := `SELECT id, title, description, status, due_date, created_at, updated_at FROM tasks ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close() // critical: leaks connections from the pool if forgotten

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err() // always check after the loop — network errors mid-scan won't surface otherwise
}

func (r *taskRepository) Update(ctx context.Context, t *domain.Task) error {
	query := `
		UPDATE tasks SET title=$1, description=$2, status=$3, due_date=$4, updated_at=$5
		WHERE id=$6`
	res, err := r.db.ExecContext(ctx, query, t.Title, t.Description, t.Status, t.DueDate, t.UpdatedAt, t.ID)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrTaskNotFound
	}
	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrTaskNotFound
	}
	return nil
}

func (r *taskRepository) FindOverdue(ctx context.Context) ([]*domain.Task, error) {
	query := `SELECT id, title, description, status, due_date, created_at, updated_at
		FROM tasks WHERE status = $1 AND due_date < now()`

	rows, err := r.db.QueryContext(ctx, query, domain.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("find overdue: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}
