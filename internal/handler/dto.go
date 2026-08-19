// internal/handler/dto.go
package handler

import (
	"time"

	"github.com/az/task-api/internal/domain"
)

type createTaskRequest struct {
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"max=2000"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

type updateTaskRequest struct {
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"max=2000"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

type taskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toTaskResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
