package student

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, student *Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockRepository) GetByUUID(ctx context.Context, uuid string) (*Student, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Student), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, student *Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context) ([]*Student, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Student), args.Error(1)
}

func (m *MockRepository) ListPaginated(ctx context.Context, params pagination.Params) ([]*Student, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*Student), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) WithDB(db *gorm.DB) Repository {
	return m
}

// Helper function to create a valid student
func createValidStudent() *Student {
	teacherID := "550e8400-e29b-41d4-a716-446655440001"
	return &Student{
		UUID:      "550e8400-e29b-41d4-a716-446655440000",
		FirstName: "Max",
		LastName:  "Mustermann",
		Email:     "max@example.com",
		TeacherID: &teacherID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// TestCreate tests the Create method of the service
func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		request     *CreateStudentRequest
		setupMock   func(*MockRepository)
		expectError bool
		errorType   *apperrors.AppError
	}{
		{
			name: "successful creation",
			request: &CreateStudentRequest{
				FirstName: "Max",
				LastName:  "Mustermann",
				Email:     "max@example.com",
			},
			setupMock: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "successful creation with teacher ID",
			request: &CreateStudentRequest{
				FirstName: "Lisa",
				LastName:  "Schmidt",
				Email:     "lisa@example.com",
				TeacherID: "550e8400-e29b-41d4-a716-446655440001",
			},
			setupMock: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "validation error - missing first name",
			request: &CreateStudentRequest{
				FirstName: "",
				LastName:  "Mustermann",
				Email:     "max@example.com",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - missing last name",
			request: &CreateStudentRequest{
				FirstName: "Max",
				LastName:  "",
				Email:     "max@example.com",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - invalid email",
			request: &CreateStudentRequest{
				FirstName: "Max",
				LastName:  "Mustermann",
				Email:     "invalid-email",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - first name too long",
			request: &CreateStudentRequest{
				FirstName: string(make([]byte, 101)), // 101 characters
				LastName:  "Mustermann",
				Email:     "max@example.com",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "repository error - duplicate email",
			request: &CreateStudentRequest{
				FirstName: "Max",
				LastName:  "Mustermann",
				Email:     "max@example.com",
			},
			setupMock: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.Anything).
					Return(apperrors.Conflict("student with this email already exists"))
			},
			expectError: true,
		},
		{
			name: "repository error - internal error",
			request: &CreateStudentRequest{
				FirstName: "Max",
				LastName:  "Mustermann",
				Email:     "max@example.com",
			},
			setupMock: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.Anything).
					Return(errors.New("database connection failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			student, err := service.Create(context.Background(), tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, student)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, student)
				assert.Equal(t, tt.request.FirstName, student.FirstName)
				assert.Equal(t, tt.request.LastName, student.LastName)
				assert.Equal(t, tt.request.Email, student.Email)
				assert.NotEmpty(t, student.UUID)
				assert.False(t, student.CreatedAt.IsZero())
				assert.False(t, student.UpdatedAt.IsZero())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGetByUUID tests the GetByUUID method
func TestGetByUUID(t *testing.T) {
	tests := []struct {
		name        string
		uuid        string
		setupMock   func(*MockRepository)
		expectError bool
	}{
		{
			name: "successful get",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *MockRepository) {
				student := createValidStudent()
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(student, nil)
			},
			expectError: false,
		},
		{
			name: "student not found",
			uuid: "nonexistent-uuid",
			setupMock: func(m *MockRepository) {
				m.On("GetByUUID", mock.Anything, "nonexistent-uuid").
					Return(nil, apperrors.NotFound("student not found"))
			},
			expectError: true,
		},
		{
			name: "repository error",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *MockRepository) {
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").
					Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			student, err := service.GetByUUID(context.Background(), tt.uuid)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, student)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, student)
				assert.Equal(t, tt.uuid, student.UUID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUpdate tests the Update method
func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		uuid        string
		request     *UpdateStudentRequest
		setupMock   func(*MockRepository)
		expectError bool
	}{
		{
			name: "successful update - all fields",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			request: &UpdateStudentRequest{
				FirstName: "Moritz",
				LastName:  "Schmidt",
				Email:     "moritz@example.com",
				TeacherID: stringPtr("550e8400-e29b-41d4-a716-446655440002"),
			},
			setupMock: func(m *MockRepository) {
				student := createValidStudent()
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(student, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "successful update - partial fields",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			request: &UpdateStudentRequest{
				FirstName: "Moritz",
			},
			setupMock: func(m *MockRepository) {
				student := createValidStudent()
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(student, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "validation error - invalid email",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			request: &UpdateStudentRequest{
				Email: "invalid-email",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - first name too long",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			request: &UpdateStudentRequest{
				FirstName: string(make([]byte, 101)),
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "student not found",
			uuid: "nonexistent-uuid",
			request: &UpdateStudentRequest{
				FirstName: "Moritz",
			},
			setupMock: func(m *MockRepository) {
				m.On("GetByUUID", mock.Anything, "nonexistent-uuid").
					Return(nil, apperrors.NotFound("student not found"))
			},
			expectError: true,
		},
		{
			name: "repository update error - duplicate email",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			request: &UpdateStudentRequest{
				Email: "duplicate@example.com",
			},
			setupMock: func(m *MockRepository) {
				student := createValidStudent()
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(student, nil)
				m.On("Update", mock.Anything, mock.Anything).
					Return(apperrors.Conflict("student with this email already exists"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			student, err := service.Update(context.Background(), tt.uuid, tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, student)
				// Verify updated fields
				if tt.request.FirstName != "" {
					assert.Equal(t, tt.request.FirstName, student.FirstName)
				}
				if tt.request.LastName != "" {
					assert.Equal(t, tt.request.LastName, student.LastName)
				}
				if tt.request.Email != "" {
					assert.Equal(t, tt.request.Email, student.Email)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestDelete tests the Delete method
func TestDelete(t *testing.T) {
	tests := []struct {
		name        string
		uuid        string
		setupMock   func(*MockRepository)
		expectError bool
	}{
		{
			name: "successful delete",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(nil)
			},
			expectError: false,
		},
		{
			name: "student not found",
			uuid: "nonexistent-uuid",
			setupMock: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "nonexistent-uuid").
					Return(apperrors.NotFound("student not found"))
			},
			expectError: true,
		},
		{
			name: "repository error",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").
					Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			err := service.Delete(context.Background(), tt.uuid)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestList tests the List method
func TestList(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockRepository)
		expectError bool
		expectCount int
	}{
		{
			name: "successful list - multiple students",
			setupMock: func(m *MockRepository) {
				students := []*Student{
					createValidStudent(),
					{
						UUID:      "550e8400-e29b-41d4-a716-446655440003",
						FirstName: "Lisa",
						LastName:  "Schmidt",
						Email:     "lisa@example.com",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				m.On("List", mock.Anything).Return(students, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "successful list - empty",
			setupMock: func(m *MockRepository) {
				m.On("List", mock.Anything).Return([]*Student{}, nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "repository error",
			setupMock: func(m *MockRepository) {
				m.On("List", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectError: true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			students, err := service.List(context.Background())

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, students)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, students)
				assert.Len(t, students, tt.expectCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
