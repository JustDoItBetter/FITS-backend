package teacher

import (
	"context"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"github.com/JustDoItBetter/FITS-backend/pkg/database"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Service handles business logic for teachers
type Service struct {
	repo     Repository
	txMgr    *database.TransactionManager
	validate *validator.Validate
}

// NewService creates a new teacher service
func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		validate: validator.New(),
	}
}

// NewServiceWithTx creates a new teacher service with transaction support
func NewServiceWithTx(repo Repository, txMgr *database.TransactionManager) *Service {
	return &Service{
		repo:     repo,
		txMgr:    txMgr,
		validate: validator.New(),
	}
}

// Create creates a new teacher
func (s *Service) Create(ctx context.Context, req *CreateTeacherRequest) (*Teacher, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(err.Error())
	}

	// Convert request to teacher entity
	teacher := req.ToTeacher()

	// Create teacher in repository
	if err := s.repo.Create(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// GetByUUID retrieves a teacher by UUID
func (s *Service) GetByUUID(ctx context.Context, uuid string) (*Teacher, error) {
	return s.repo.GetByUUID(ctx, uuid)
}

// Update updates an existing teacher
// Uses transactions to prevent race conditions between read and write operations
func (s *Service) Update(ctx context.Context, uuid string, req *UpdateTeacherRequest) (*Teacher, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(err.Error())
	}

	// If transaction manager is available, use transaction for atomic read-update
	// This prevents race conditions where another request modifies the teacher
	// between our read and write operations
	if s.txMgr != nil {
		return database.WithTransactionValue(ctx, s.txMgr, func(tx *gorm.DB) (*Teacher, error) {
			// Get existing teacher within transaction
			txRepo := s.repo.WithDB(tx)
			teacher, err := txRepo.GetByUUID(ctx, uuid)
			if err != nil {
				return nil, err
			}

			// Update teacher fields
			teacher.Update(req)

			// Save updated teacher within same transaction
			if err := txRepo.Update(ctx, teacher); err != nil {
				return nil, err
			}

			return teacher, nil
		})
	}

	// Fallback to non-transactional update (backward compatibility)
	teacher, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	teacher.Update(req)

	if err := s.repo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// Delete deletes a teacher by UUID
func (s *Service) Delete(ctx context.Context, uuid string) error {
	return s.repo.Delete(ctx, uuid)
}

// List retrieves all teachers (deprecated: use ListPaginated)
func (s *Service) List(ctx context.Context) ([]*Teacher, error) {
	return s.repo.List(ctx)
}

// ListPaginated retrieves teachers with pagination
// Returns teachers slice, total count, and error
func (s *Service) ListPaginated(ctx context.Context, params pagination.Params) ([]*Teacher, int64, error) {
	return s.repo.ListPaginated(ctx, params)
}
