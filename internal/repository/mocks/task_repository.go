package mocks

import (
	"context"

	"github.com/az/task-api/internal/domain"
	"github.com/az/task-api/internal/repository"
)

type TaskRepository struct {
	Tasks       map[string]*domain.Task
	ErrToReturn error
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{Tasks: make(map[string]*domain.Task)}
}

func (m *TaskRepository) Create(ctx context.Context, t *domain.Task) error {
	if m.ErrToReturn != nil {
		return m.ErrToReturn
	}
	t.ID = "mock-id-1"
	m.Tasks[t.ID] = t
	return nil
}

func (m *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	if m.ErrToReturn != nil {
		return nil, m.ErrToReturn
	}
	t, ok := m.Tasks[id]
	if !ok {
		return nil, repository.ErrTaskNotFound
	}
	return t, nil
}

func (m *TaskRepository) List(ctx context.Context) ([]*domain.Task, error) {
	var out []*domain.Task
	for _, t := range m.Tasks {
		out = append(out, t)
	}
	return out, m.ErrToReturn
}

func (m *TaskRepository) Update(ctx context.Context, t *domain.Task) error {
	if m.ErrToReturn != nil {
		return m.ErrToReturn
	}
	if _, ok := m.Tasks[t.ID]; !ok {
		return repository.ErrTaskNotFound
	}
	m.Tasks[t.ID] = t
	return nil
}

func (m *TaskRepository) Delete(ctx context.Context, id string) error {
	if m.ErrToReturn != nil {
		return m.ErrToReturn
	}
	if _, ok := m.Tasks[id]; !ok {
		return repository.ErrTaskNotFound
	}
	delete(m.Tasks, id)
	return nil
}

func (m *TaskRepository) FindOverdue(ctx context.Context) ([]*domain.Task, error) {
	return nil, m.ErrToReturn
}
