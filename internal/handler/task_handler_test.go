// internal/handler/task_handler_test.go
package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/az/task-api/internal/handler"
	"github.com/az/task-api/internal/repository/mocks"
	"github.com/az/task-api/internal/usecase"
	"github.com/go-chi/chi/v5"
)

// reuse the same mock — export it or duplicate a minimal version in this package.
// simplest: define a tiny local mock here since usecase_test's is unexported to that package.

func TestCreateTaskHandler_Success(t *testing.T) {
	repo := mocks.NewTaskRepository() // local mock, same shape as Step 13's usecase mock
	u := usecase.NewTaskUsecase(repo)
	h := handler.NewTaskHandler(u)

	body, _ := json.Marshal(map[string]any{
		"title":    "Test task",
		"due_date": time.Now().Add(24 * time.Hour),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestGetTaskHandler_NotFound(t *testing.T) {
	repo := mocks.NewTaskRepository()
	u := usecase.NewTaskUsecase(repo)
	h := handler.NewTaskHandler(u)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/missing-id", nil)
	rec := httptest.NewRecorder()

	// chi URL params aren't parsed from raw httptest requests automatically —
	// inject the route context manually so chi.URLParam(r, "id") works.
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "missing-id")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.Get(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}
