package student

import (
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/common/validation"
	"github.com/google/uuid"
)

// Student represents a student entity
// @Description Student information
type Student struct {
	UUID      string    `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid"`
	FirstName string    `json:"first_name" example:"Max" validate:"required,min=1,max=100"`
	LastName  string    `json:"last_name" example:"Mustermann" validate:"required,min=1,max=100"`
	Email     string    `json:"email" example:"max@example.com" validate:"required,email"`
	TeacherID *string   `json:"teacher_id,omitempty" example:"teacher-uuid-123"` // Optional - can be NULL
	CreatedAt time.Time `json:"created_at" example:"2025-09-30T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-09-30T12:00:00Z"`
}

// CreateStudentRequest represents the request to create a new student
// @Description Request body for creating a new student
type CreateStudentRequest struct {
	FirstName string `json:"first_name" example:"Max" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" example:"Mustermann" validate:"required,min=1,max=100"`
	Email     string `json:"email" example:"max@example.com" validate:"required,email"`
	TeacherID string `json:"teacher_id,omitempty" example:"teacher-uuid-123" validate:"omitempty,uuid"` // Optional - can be assigned later
}

// UpdateStudentRequest represents the request to update a student
// @Description Request body for updating a student
type UpdateStudentRequest struct {
	FirstName string  `json:"first_name,omitempty" example:"Moritz" validate:"omitempty,min=1,max=100"`
	LastName  string  `json:"last_name,omitempty" example:"Schmidt" validate:"omitempty,min=1,max=100"`
	Email     string  `json:"email,omitempty" example:"moritz@example.com" validate:"omitempty,email"`
	TeacherID *string `json:"teacher_id,omitempty" example:"teacher-uuid-456" validate:"omitempty,uuid"`
}

// ToStudent converts CreateStudentRequest to Student entity
// Sanitizes input to prevent XSS and injection attacks
// UUID is always generated server-side for security
func (r *CreateStudentRequest) ToStudent() *Student {
	now := time.Now()
	studentUUID := uuid.New().String() // Server always generates UUID

	var teacherID *string
	if r.TeacherID != "" {
		teacherID = &r.TeacherID
	}

	return &Student{
		UUID:      studentUUID,
		FirstName: validation.SanitizeName(r.FirstName),
		LastName:  validation.SanitizeName(r.LastName),
		Email:     validation.SanitizeEmail(r.Email),
		TeacherID: teacherID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates student fields from UpdateStudentRequest
// Sanitizes input to prevent XSS attacks before updating
func (s *Student) Update(req *UpdateStudentRequest) {
	if req.FirstName != "" {
		s.FirstName = validation.SanitizeName(req.FirstName)
	}
	if req.LastName != "" {
		s.LastName = validation.SanitizeName(req.LastName)
	}
	if req.Email != "" {
		s.Email = validation.SanitizeEmail(req.Email)
	}
	if req.TeacherID != nil {
		s.TeacherID = req.TeacherID
	}
	s.UpdatedAt = time.Now()
}
