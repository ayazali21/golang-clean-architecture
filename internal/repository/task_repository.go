package repository

import (
	"context"

	"github.com/az/task-api/internal/apperror"
	"github.com/az/task-api/internal/domain"
)

var ErrTaskNotFound = apperror.NotFound("task not found")

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	List(ctx context.Context) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id string) error
	FindOverdue(ctx context.Context) ([]*domain.Task, error) // for the scheduler in Step 12
}
