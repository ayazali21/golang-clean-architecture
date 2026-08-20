// cmd/api/main.go

// @title           Task API
// @version         1.0
// @description     A small task management REST API built for learning Clean Architecture in Go.
// @host            localhost:8080
// @BasePath        /api/v1
package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/az/task-api/internal/config"
	"github.com/az/task-api/internal/handler"
	"github.com/az/task-api/internal/infrastructure/database"
	"github.com/az/task-api/internal/repository/postgres"
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
	slog.Info("config loaded", "env", cfg.AppEnv, "port", cfg.HTTPPort)

	//Database
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Composition root: dependency graph flows top to bottom.
	taskRepo := postgres.NewTaskRepository(db)
	taskUsecase := usecase.NewTaskUsecase(taskRepo)
	taskHandler := handler.NewTaskHandler(taskUsecase)

	router := handler.NewRouter(taskHandler)

	slog.Info("starting server", "port", cfg.HTTPPort)

	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
