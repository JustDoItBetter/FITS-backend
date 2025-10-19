package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		want string
	}{
		{
			name: "error with message and details",
			err:  &AppError{Message: "test error", Details: "some details", Code: 400},
			want: "test error: some details",
		},
		{
			name: "error with message only",
			err:  &AppError{Message: "test error", Code: 400},
			want: "test error",
		},
		{
			name: "error with empty message",
			err:  &AppError{Message: "", Code: 500},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func TestNewAppError(t *testing.T) {
	err := NewAppError(404, "not found", "resource not found")

	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "not found", err.Message)
	assert.Equal(t, "resource not found", err.Details)
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("invalid input")

	assert.Equal(t, 400, err.Code)
	assert.Equal(t, "Bad Request", err.Message)
	assert.Equal(t, "invalid input", err.Details)
}

func TestNotFound(t *testing.T) {
	err := NotFound("user")

	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "Not Found", err.Message)
	assert.Equal(t, "user not found", err.Details)
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("invalid token")

	assert.Equal(t, 401, err.Code)
	assert.Equal(t, "Unauthorized", err.Message)
	assert.Equal(t, "invalid token", err.Details)
}

func TestInternalServerError(t *testing.T) {
	err := InternalServerError("database error")

	assert.Equal(t, 500, err.Code)
	assert.Equal(t, "Internal Server Error", err.Message)
	assert.Equal(t, "database error", err.Details)
}

func TestValidationError(t *testing.T) {
	err := ValidationError("email is required")

	assert.Equal(t, 422, err.Code)
	assert.Equal(t, "Validation Error", err.Message)
	assert.Equal(t, "email is required", err.Details)
}

func TestConflict(t *testing.T) {
	err := Conflict("email already exists")

	assert.Equal(t, 409, err.Code)
	assert.Equal(t, "Conflict", err.Message)
	assert.Equal(t, "email already exists", err.Details)
}

func TestInternal(t *testing.T) {
	err := Internal("unexpected error")

	assert.Equal(t, 500, err.Code)
	assert.Equal(t, "Internal Server Error", err.Message)
	assert.Equal(t, "unexpected error", err.Details)
}

func TestIsUniqueViolation(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "unique violation - duplicate key",
			err:  errors.New("duplicate key value violates unique constraint"),
			want: true,
		},
		{
			name: "unique violation - UNIQUE constraint",
			err:  errors.New("UNIQUE constraint failed"),
			want: true,
		},
		{
			name: "not a unique violation",
			err:  errors.New("some other error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsUniqueViolation(tt.err))
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{
			name:   "substring found in middle",
			s:      "this is a test string",
			substr: "test",
			want:   true,
		},
		{
			name:   "substring at beginning",
			s:      "hello world",
			substr: "hello",
			want:   true,
		},
		{
			name:   "substring at end",
			s:      "hello world",
			substr: "world",
			want:   true,
		},
		{
			name:   "substring not found",
			s:      "hello world",
			substr: "missing",
			want:   false,
		},
		{
			name:   "exact match",
			s:      "test",
			substr: "test",
			want:   true,
		},
		{
			name:   "empty substring",
			s:      "test",
			substr: "",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, contains(tt.s, tt.substr))
		})
	}
}

func TestFindInString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{
			name:   "pattern found",
			s:      "this is a test string",
			substr: "test",
			want:   true,
		},
		{
			name:   "pattern not found",
			s:      "this is a test string",
			substr: "missing",
			want:   false,
		},
		{
			name:   "empty pattern",
			s:      "test",
			substr: "",
			want:   true,
		},
		{
			name:   "pattern at start",
			s:      "test string",
			substr: "test",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, findInString(tt.s, tt.substr))
		})
	}
}
