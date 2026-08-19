// internal/usecase/task_usecase.go
package usecase

import (
	"context"
	"time"

	"github.com/az/task-api/internal/domain"
	"github.com/az/task-api/internal/repository"
)

type TaskUsecase struct {
	repo repository.TaskRepository
}

func NewTaskUsecase(repo repository.TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

func (u *TaskUsecase) CreateTask(ctx context.Context, title, description string, dueDate time.Time) (*domain.Task, error) {
	task := domain.NewTask(title, description, dueDate)
	if err := u.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *TaskUsecase) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *TaskUsecase) ListTasks(ctx context.Context) ([]*domain.Task, error) {
	return u.repo.List(ctx)
}

func (u *TaskUsecase) UpdateTask(ctx context.Context, id, title, description string, dueDate time.Time) (*domain.Task, error) {
	task, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Title = title
	task.Description = description
	task.DueDate = dueDate
	task.UpdatedAt = time.Now().UTC()

	if err := u.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *TaskUsecase) DeleteTask(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *TaskUsecase) CompleteTask(ctx context.Context, id string) (*domain.Task, error) {
	task, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.MarkCompleted() // domain rule lives on the entity, usecase just orchestrates

	if err := u.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}
