//go:build integration

// internal/repository/postgres/task_repository_test.go
package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/az/task-api/internal/domain"
	"github.com/az/task-api/internal/repository"
	"github.com/az/task-api/internal/repository/postgres"
)

func TestTaskRepository_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewTaskRepository(db)
	ctx := context.Background()

	task := domain.NewTask("Integration test task", "desc", time.Now().Add(24*time.Hour))

	if err := repo.Create(ctx, task); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if task.ID == "" {
		t.Fatal("expected DB to populate ID via RETURNING id")
	}

	fetched, err := repo.GetByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if fetched.Title != task.Title {
		t.Errorf("expected title %q, got %q", task.Title, fetched.Title)
	}
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewTaskRepository(db)

	_, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")

	if err != repository.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskRepository_FindOverdue(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewTaskRepository(db)
	ctx := context.Background()

	pastDue := domain.NewTask("Overdue task", "", time.Now().Add(-1*time.Hour))
	if err := repo.Create(ctx, pastDue); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	overdue, err := repo.FindOverdue(ctx)
	if err != nil {
		t.Fatalf("find overdue failed: %v", err)
	}
	if len(overdue) != 1 {
		t.Fatalf("expected 1 overdue task, got %d", len(overdue))
	}
}
