package student

import (
	"context"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"github.com/go-playground/validator/v10"
)

// Service handles business logic for students
type Service struct {
	repo     Repository
	validate *validator.Validate
}

// NewService creates a new student service
func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		validate: validator.New(),
	}
}

// Create creates a new student
func (s *Service) Create(ctx context.Context, req *CreateStudentRequest) (*Student, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(err.Error())
	}

	// Convert request to student entity
	student := req.ToStudent()

	// Create student in repository
	if err := s.repo.Create(ctx, student); err != nil {
		return nil, err
	}

	return student, nil
}

// GetByUUID retrieves a student by UUID
func (s *Service) GetByUUID(ctx context.Context, uuid string) (*Student, error) {
	return s.repo.GetByUUID(ctx, uuid)
}

// Update updates an existing student
func (s *Service) Update(ctx context.Context, uuid string, req *UpdateStudentRequest) (*Student, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(err.Error())
	}

	// Get existing student
	student, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Update student fields
	student.Update(req)

	// Save updated student
	if err := s.repo.Update(ctx, student); err != nil {
		return nil, err
	}

	return student, nil
}

// Delete deletes a student by UUID
func (s *Service) Delete(ctx context.Context, uuid string) error {
	return s.repo.Delete(ctx, uuid)
}

// List retrieves all students (deprecated: use ListPaginated)
func (s *Service) List(ctx context.Context) ([]*Student, error) {
	return s.repo.List(ctx)
}

// ListPaginated retrieves students with pagination
// Returns students slice, total count, and error
func (s *Service) ListPaginated(ctx context.Context, params pagination.Params) ([]*Student, int64, error) {
	return s.repo.ListPaginated(ctx, params)
}
