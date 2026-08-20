// internal/handler/error.go
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/az/task-api/internal/apperror"
)

type errorResponseBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		// Unknown error type reached the handler — treat as internal,
		// never leak the raw message to the client.
		appErr = apperror.Internal("unexpected error", err)
	}

	resp := errorResponseBody{}
	resp.Error.Code = string(appErr.Code)
	resp.Error.Message = appErr.Message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus())
	json.NewEncoder(w).Encode(resp)
}
