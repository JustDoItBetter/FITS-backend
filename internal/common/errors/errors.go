package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"invalid request"`
	Details string `json:"details,omitempty" example:"field 'email' is required"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error constructors
func BadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, "Bad Request", message)
}

func NotFound(resource string) *AppError {
	return NewAppError(http.StatusNotFound, "Not Found", fmt.Sprintf("%s not found", resource))
}

func Unauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, "Unauthorized", message)
}

func InternalServerError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, "Internal Server Error", message)
}

func ValidationError(message string) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, "Validation Error", message)
}

func Conflict(message string) *AppError {
	return NewAppError(http.StatusConflict, "Conflict", message)
}

func Internal(message string) *AppError {
	return InternalServerError(message)
}

// IsUniqueViolation checks if the error is a unique constraint violation
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// PostgreSQL unique violation error code: 23505
	// GORM wraps this in the error message
	errMsg := err.Error()
	return contains(errMsg, "duplicate key") ||
		contains(errMsg, "unique constraint") ||
		contains(errMsg, "UNIQUE constraint failed") ||
		contains(errMsg, "23505")
}

// contains checks if s contains substr (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
