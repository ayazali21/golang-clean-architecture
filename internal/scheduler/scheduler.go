// internal/scheduler/scheduler.go
package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/az/task-api/internal/usecase"
)

type TaskScheduler struct {
	usecase  *usecase.TaskUsecase
	interval time.Duration
}

func NewTaskScheduler(u *usecase.TaskUsecase, interval time.Duration) *TaskScheduler {
	return &TaskScheduler{usecase: u, interval: interval}
}

// Run blocks until ctx is canceled. Call it in its own goroutine from main.go.
func (s *TaskScheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop() // always stop tickers — leaked tickers leak goroutines/timers

	slog.Info("scheduler started", "interval", s.interval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("scheduler stopping")
			return
		case <-ticker.C:
			s.runOnce(ctx)
		}
	}
}

// runOnce is synchronous by design within a single tick: the next tick
// won't fire logic again until this returns, since we're driven by the
// ticker channel in a single select loop — this alone prevents overlap.
func (s *TaskScheduler) runOnce(ctx context.Context) {
	// Guard against a slow run blocking past its own interval indefinitely.
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	count, err := s.usecase.ProcessOverdueTasks(ctx)
	if err != nil {
		slog.Error("overdue task processing failed", "error", err)
		return
	}
	slog.Info("overdue tasks processed", "count", count)
}
