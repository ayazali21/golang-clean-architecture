package main

import (
	"log/slog"
	"os"

	"github.com/az/task-api/internal/config"
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

}
