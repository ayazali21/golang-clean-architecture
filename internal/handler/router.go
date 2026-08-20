// internal/handler/router.go
package handler

import (
	_ "github.com/az/task-api/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(taskHandler *TaskHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // replaced with slog-based logging in Step 11
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Post("/", taskHandler.Create)
		r.Get("/", taskHandler.List)
		r.Get("/{id}", taskHandler.Get)
		r.Put("/{id}", taskHandler.Update)
		r.Delete("/{id}", taskHandler.Delete)
		r.Post("/{id}/complete", taskHandler.Complete)
	})

	return r
}
