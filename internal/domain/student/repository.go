package student

import (
	"context"
	"sync"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
)

// Repository defines the interface for student data access
type Repository interface {
	Create(ctx context.Context, student *Student) error
	GetByUUID(ctx context.Context, uuid string) (*Student, error)
	Update(ctx context.Context, student *Student) error
	Delete(ctx context.Context, uuid string) error
	// ListPaginated returns paginated students with total count for metadata
	ListPaginated(ctx context.Context, params pagination.Params) ([]*Student, int64, error)
	// List retrieves all students (deprecated: use ListPaginated for better performance)
	List(ctx context.Context) ([]*Student, error)
}

// InMemoryRepository is a simple in-memory implementation of Repository
// This is temporary until we implement a real database
type InMemoryRepository struct {
	students map[string]*Student
	mu       sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		students: make(map[string]*Student),
	}
}

// Create adds a new student to the repository
func (r *InMemoryRepository) Create(ctx context.Context, student *Student) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.students[student.UUID]; exists {
		return errors.Conflict("student with this UUID already exists")
	}

	r.students[student.UUID] = student
	return nil
}

// GetByUUID retrieves a student by UUID
func (r *InMemoryRepository) GetByUUID(ctx context.Context, uuid string) (*Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	student, exists := r.students[uuid]
	if !exists {
		return nil, errors.NotFound("student")
	}

	return student, nil
}

// Update updates an existing student
func (r *InMemoryRepository) Update(ctx context.Context, student *Student) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.students[student.UUID]; !exists {
		return errors.NotFound("student")
	}

	r.students[student.UUID] = student
	return nil
}

// Delete removes a student from the repository
func (r *InMemoryRepository) Delete(ctx context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.students[uuid]; !exists {
		return errors.NotFound("student")
	}

	delete(r.students, uuid)
	return nil
}

// List retrieves all students (deprecated: use ListPaginated)
func (r *InMemoryRepository) List(ctx context.Context) ([]*Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	students := make([]*Student, 0, len(r.students))
	for _, student := range r.students {
		students = append(students, student)
	}

	return students, nil
}

// ListPaginated retrieves students with pagination support
// Returns slice of students, total count, and error
func (r *InMemoryRepository) ListPaginated(ctx context.Context, params pagination.Params) ([]*Student, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all students first
	allStudents := make([]*Student, 0, len(r.students))
	for _, student := range r.students {
		allStudents = append(allStudents, student)
	}

	totalCount := int64(len(allStudents))

	// Apply pagination
	start := params.Offset()
	end := start + params.Limit

	// Handle edge cases
	if start >= len(allStudents) {
		return []*Student{}, totalCount, nil
	}
	if end > len(allStudents) {
		end = len(allStudents)
	}

	return allStudents[start:end], totalCount, nil
}
