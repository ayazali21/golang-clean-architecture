// internal/usecase/task_usecase_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/az/task-api/internal/apperror"
	"github.com/az/task-api/internal/domain"
	"github.com/az/task-api/internal/repository"
	"github.com/az/task-api/internal/repository/mocks"
	"github.com/az/task-api/internal/usecase"
)

func TestCreateTask_Success(t *testing.T) {
	repo := mocks.NewTaskRepository()
	u := usecase.NewTaskUsecase(repo)

	task, err := u.CreateTask(context.Background(), "Write tests", "desc", time.Now().Add(24*time.Hour))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.ID == "" {
		t.Error("expected task to have an ID after creation")
	}
	if task.Status != domain.StatusPending {
		t.Errorf("expected status pending, got %s", task.Status)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	repo := mocks.NewTaskRepository()
	u := usecase.NewTaskUsecase(repo)

	_, err := u.GetTask(context.Background(), "nonexistent-id")

	if !errors.Is(err, repository.ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestCompleteTask_AlreadyCompleted(t *testing.T) {
	repo := mocks.NewTaskRepository()
	repo.Tasks["t1"] = &domain.Task{ID: "t1", Status: domain.StatusCompleted}
	u := usecase.NewTaskUsecase(repo)

	_, err := u.CompleteTask(context.Background(), "t1")

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Code != apperror.CodeConflict {
		t.Errorf("expected CONFLICT error, got %v", err)
	}
}

func TestCompleteTask_RepositoryError(t *testing.T) {
	repo := mocks.NewTaskRepository()
	repo.Tasks["t1"] = &domain.Task{ID: "t1", Status: domain.StatusPending}
	repo.ErrToReturn = errors.New("boom") // only affects calls made AFTER this point in real usage;
	// here it's set before GetByID, so it fails immediately
	u := usecase.NewTaskUsecase(repo)

	_, err := u.CompleteTask(context.Background(), "t1")

	if err == nil {
		t.Fatal("expected error from repository failure, got nil")
	}
}
