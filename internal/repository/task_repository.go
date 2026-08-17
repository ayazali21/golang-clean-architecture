package repository

import (
	"context"

	"github.com/az/task-api/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	Update(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, ID string)
}
