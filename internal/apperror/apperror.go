// internal/apperror/apperror.go
package apperror

import "net/http"

type Code string

const (
	CodeNotFound   Code = "NOT_FOUND"
	CodeValidation Code = "VALIDATION_ERROR"
	CodeConflict   Code = "CONFLICT"
	CodeInternal   Code = "INTERNAL_ERROR"
)

type AppError struct {
	Code    Code
	Message string
	Err     error // wrapped original error, never exposed to client
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func NotFound(msg string) *AppError {
	return &AppError{Code: CodeNotFound, Message: msg}
}

func Validation(msg string) *AppError {
	return &AppError{Code: CodeValidation, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: CodeConflict, Message: msg}
}

func Internal(msg string, err error) *AppError {
	return &AppError{Code: CodeInternal, Message: msg, Err: err}
}

// HTTPStatus is the single source of truth for code → status mapping.
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeValidation:
		return http.StatusBadRequest
	case CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
