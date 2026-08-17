package domain

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
	StatusOverdue   TaskStatus = "overdue"
)

type Task struct {
	ID          string
	Title       string
	Description string
	Status      TaskStatus
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(title, desc string, dueDate time.Time) *Task {
	now := time.Now().UTC()
	return &Task{
		Title:       title,
		Description: desc,
		Status:      StatusPending,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) MarkCompleted() {
	t.Status = StatusCompleted
	t.UpdatedAt = time.Now().UTC()
}

func (t *Task) IsOverDue() bool {
	return time.Now().UTC().After(t.DueDate) && t.Status == StatusPending
}
