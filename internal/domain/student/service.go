package student

import (
	"context"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"github.com/JustDoItBetter/FITS-backend/pkg/database"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Service handles business logic for students
type Service struct {
	repo     Repository
	txMgr    *database.TransactionManager
	validate *validator.Validate
}

// NewService creates a new student service
func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		validate: validator.New(),
	}
}

// NewServiceWithTx creates a new student service with transaction support
func NewServiceWithTx(repo Repository, txMgr *database.TransactionManager) *Service {
	return &Service{
		repo:     repo,
		txMgr:    txMgr,
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
// Uses transactions to prevent race conditions between read and write operations
func (s *Service) Update(ctx context.Context, uuid string, req *UpdateStudentRequest) (*Student, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(err.Error())
	}

	// If transaction manager is available, use transaction for atomic read-update
	// This prevents race conditions where another request modifies the student
	// between our read and write operations
	if s.txMgr != nil {
		return database.WithTransactionValue(ctx, s.txMgr, func(tx *gorm.DB) (*Student, error) {
			// Get existing student within transaction
			txRepo := s.repo.WithDB(tx)
			student, err := txRepo.GetByUUID(ctx, uuid)
			if err != nil {
				return nil, err
			}

			// Update student fields
			student.Update(req)

			// Save updated student within same transaction
			if err := txRepo.Update(ctx, student); err != nil {
				return nil, err
			}

			return student, nil
		})
	}

	// Fallback to non-transactional update (backward compatibility)
	student, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	student.Update(req)

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
