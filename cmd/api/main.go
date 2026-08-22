// cmd/api/main.go

// @title           Task API
// @version         1.0
// @description     A small task management REST API built for learning Clean Architecture in Go.
// @host            localhost:8080
// @BasePath        /api/v1
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/az/task-api/internal/config"
	"github.com/az/task-api/internal/handler"
	"github.com/az/task-api/internal/infrastructure/database"
	"github.com/az/task-api/internal/repository/postgres"
	"github.com/az/task-api/internal/scheduler"
	"github.com/az/task-api/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}

	taskRepo := postgres.NewTaskRepository(db)
	taskUsecase := usecase.NewTaskUsecase(taskRepo)
	taskHandler := handler.NewTaskHandler(taskUsecase)
	router := handler.NewRouter(taskHandler)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	// ctx is canceled the moment SIGINT/SIGTERM arrives — this single ctx
	// is what the scheduler listens on to know when to stop.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	taskScheduler := scheduler.NewTaskScheduler(taskUsecase, 1*time.Minute)
	schedulerDone := make(chan struct{})
	go func() {
		taskScheduler.Run(ctx)
		close(schedulerDone) // signals scheduler has fully stopped
	}()

	// Server runs in its own goroutine so we can block on ctx.Done() below.
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting server", "port", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-serverErr:
		slog.Error("server failed unexpectedly", "error", err)
	}

	// Bounded timeout: don't wait forever for in-flight requests. Graceful Shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	<-schedulerDone // wait for the current scheduler tick (if any) to finish
	slog.Info("scheduler stopped")

	if err := db.Close(); err != nil {
		slog.Error("db close error", "error", err)
	}

	slog.Info("shutdown complete")
}
