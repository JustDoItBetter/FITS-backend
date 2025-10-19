package teacher

import (
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/common/validation"
	"github.com/google/uuid"
)

// Teacher represents a teacher entity
// @Description Teacher information
type Teacher struct {
	UUID       string    `json:"uuid" example:"teacher-uuid-123" validate:"required,uuid"`
	FirstName  string    `json:"first_name" example:"Anna" validate:"required,min=1,max=100"`
	LastName   string    `json:"last_name" example:"Schmidt" validate:"required,min=1,max=100"`
	Email      string    `json:"email" example:"anna@example.com" validate:"required,email"`
	Department string    `json:"department" example:"Computer Science" validate:"required,min=1,max=100"`
	CreatedAt  time.Time `json:"created_at" example:"2025-09-30T12:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2025-09-30T12:00:00Z"`
}

// CreateTeacherRequest represents the request to create a new teacher
// @Description Request body for creating a new teacher
type CreateTeacherRequest struct {
	FirstName  string `json:"first_name" example:"Anna" validate:"required,min=1,max=100"`
	LastName   string `json:"last_name" example:"Schmidt" validate:"required,min=1,max=100"`
	Email      string `json:"email" example:"anna@example.com" validate:"required,email"`
	Department string `json:"department" example:"Computer Science" validate:"required,min=1,max=100"`
}

// UpdateTeacherRequest represents the request to update a teacher
// @Description Request body for updating a teacher
type UpdateTeacherRequest struct {
	FirstName  string `json:"first_name,omitempty" example:"Maria" validate:"omitempty,min=1,max=100"`
	LastName   string `json:"last_name,omitempty" example:"Müller" validate:"omitempty,min=1,max=100"`
	Email      string `json:"email,omitempty" example:"maria@example.com" validate:"omitempty,email"`
	Department string `json:"department,omitempty" example:"Mathematics" validate:"omitempty,min=1,max=100"`
}

// ToTeacher converts CreateTeacherRequest to Teacher entity
// UUID is always generated server-side for security
func (r *CreateTeacherRequest) ToTeacher() *Teacher {
	now := time.Now()

	return &Teacher{
		UUID:       uuid.New().String(), // Server always generates UUID
		FirstName:  validation.SanitizeName(r.FirstName),
		LastName:   validation.SanitizeName(r.LastName),
		Email:      validation.SanitizeEmail(r.Email),
		Department: validation.SanitizeName(r.Department),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Update updates teacher fields from UpdateTeacherRequest
func (t *Teacher) Update(req *UpdateTeacherRequest) {
	if req.FirstName != "" {
		t.FirstName = validation.SanitizeName(req.FirstName)
	}
	if req.LastName != "" {
		t.LastName = validation.SanitizeName(req.LastName)
	}
	if req.Email != "" {
		t.Email = validation.SanitizeEmail(req.Email)
	}
	if req.Department != "" {
		t.Department = validation.SanitizeName(req.Department)
	}
	t.UpdatedAt = time.Now()
}
