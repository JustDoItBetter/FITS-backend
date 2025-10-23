package teacher

import (
	"context"
	"sync"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"gorm.io/gorm"
)

// Repository defines the interface for teacher data access
type Repository interface {
	Create(ctx context.Context, teacher *Teacher) error
	GetByUUID(ctx context.Context, uuid string) (*Teacher, error)
	Update(ctx context.Context, teacher *Teacher) error
	Delete(ctx context.Context, uuid string) error
	// ListPaginated returns paginated teachers with total count for metadata
	ListPaginated(ctx context.Context, params pagination.Params) ([]*Teacher, int64, error)
	// List retrieves all teachers (deprecated: use ListPaginated for better performance)
	List(ctx context.Context) ([]*Teacher, error)
	// WithDB returns a new repository instance using the provided database connection
	// This enables the repository to participate in transactions
	WithDB(db *gorm.DB) Repository
}

// InMemoryRepository is a simple in-memory implementation of Repository
type InMemoryRepository struct {
	teachers map[string]*Teacher
	mu       sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		teachers: make(map[string]*Teacher),
	}
}

// Create adds a new teacher to the repository
func (r *InMemoryRepository) Create(ctx context.Context, teacher *Teacher) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teachers[teacher.UUID]; exists {
		return errors.Conflict("teacher with this UUID already exists")
	}

	r.teachers[teacher.UUID] = teacher
	return nil
}

// GetByUUID retrieves a teacher by UUID
func (r *InMemoryRepository) GetByUUID(ctx context.Context, uuid string) (*Teacher, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	teacher, exists := r.teachers[uuid]
	if !exists {
		return nil, errors.NotFound("teacher")
	}

	return teacher, nil
}

// Update updates an existing teacher
func (r *InMemoryRepository) Update(ctx context.Context, teacher *Teacher) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teachers[teacher.UUID]; !exists {
		return errors.NotFound("teacher")
	}

	r.teachers[teacher.UUID] = teacher
	return nil
}

// Delete removes a teacher from the repository
func (r *InMemoryRepository) Delete(ctx context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teachers[uuid]; !exists {
		return errors.NotFound("teacher")
	}

	delete(r.teachers, uuid)
	return nil
}

// List retrieves all teachers (deprecated: use ListPaginated)
func (r *InMemoryRepository) List(ctx context.Context) ([]*Teacher, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	teachers := make([]*Teacher, 0, len(r.teachers))
	for _, teacher := range r.teachers {
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

// ListPaginated retrieves teachers with pagination support
// Returns slice of teachers, total count, and error
func (r *InMemoryRepository) ListPaginated(ctx context.Context, params pagination.Params) ([]*Teacher, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all teachers first
	allTeachers := make([]*Teacher, 0, len(r.teachers))
	for _, teacher := range r.teachers {
		allTeachers = append(allTeachers, teacher)
	}

	totalCount := int64(len(allTeachers))

	// Apply pagination
	start := params.Offset()
	end := start + params.Limit

	// Handle edge cases
	if start >= len(allTeachers) {
		return []*Teacher{}, totalCount, nil
	}
	if end > len(allTeachers) {
		end = len(allTeachers)
	}

	return allTeachers[start:end], totalCount, nil
}

// WithDB returns the same repository instance (in-memory doesn't use database connections)
// This is a no-op implementation to satisfy the Repository interface
func (r *InMemoryRepository) WithDB(db *gorm.DB) Repository {
	return r
}
