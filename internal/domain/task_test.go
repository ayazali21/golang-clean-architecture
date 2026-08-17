package domain

import (
	"testing"
	"time"
)

func TestMarkCompleted(t *testing.T) {
	task := NewTask("test task", "some description", time.Now().Add(24*time.Hour))
	task.MarkCompleted()
	if task.Status != StatusCompleted {
		t.Errorf("expected status %s, got %s", StatusCompleted, task.Status)
	}
}
