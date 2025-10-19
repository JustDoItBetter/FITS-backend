package teacher

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, teacher *Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

func (m *MockRepository) GetByUUID(ctx context.Context, uuid string) (*Teacher, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Teacher), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, teacher *Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context) ([]*Teacher, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Teacher), args.Error(1)
}

func (m *MockRepository) ListPaginated(ctx context.Context, params pagination.Params) ([]*Teacher, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*Teacher), args.Get(1).(int64), args.Error(2)
}

// Helper function to create a valid teacher
func createValidTeacher() *Teacher {
	return &Teacher{
		UUID:       "550e8400-e29b-41d4-a716-446655440010",
		FirstName:  "Anna",
		LastName:   "Schmidt",
		Email:      "anna@example.com",
		Department: "Computer Science",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// TestCreate tests the Create method of the service
func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		request     *CreateTeacherRequest
		setupMock   func(*MockRepository)
		expectError bool
		errorType   *apperrors.AppError
	}{
		{
			name: "successful creation",
			request: &CreateTeacherRequest{
				FirstName:  "Anna",
				LastName:   "Schmidt",
				Email:      "anna@example.com",
				Department: "Computer Science",
			},
			setupMock: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "successful creation with different department",
			request: &CreateTeacherRequest{
				FirstName:  "Thomas",
				LastName:   "Müller",
				Email:      "thomas@example.com",
				Department: "Mathematics",
			},
			setupMock: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "validation error - missing first name",
			request: &CreateTeacherRequest{
				FirstName:  "",
				LastName:   "Schmidt",
				Email:      "anna@example.com",
				Department: "Computer Science",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - missing last name",
			request: &CreateTeacherRequest{
				FirstName:  "Anna",
				LastName:   "",
				Email:      "anna@example.com",
				Department: "Computer Science",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - invalid email",
			request: &CreateTeacherRequest{
				FirstName:  "Anna",
				LastName:   "Schmidt",
				Email:      "invalid-email",
				Department: "Computer Science",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - missing department",
			request: &CreateTeacherRequest{
				FirstName:  "Anna",
				LastName:   "Schmidt",
				Email:      "anna@example.com",
				Department: "",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - first name too long",
			request: &CreateTeacherRequest{
				FirstName:  string(make([]byte, 101)), // 101 characters
				LastName:   "Schmidt",
				Email:      "anna@example.com",
				Department: "Computer Science",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - department too long",
			request: &CreateTeacherRequest{
				FirstName:  "Anna",
				LastName:   "Schmidt",
				Email:      "anna@example.com",
				Department: string(make([]byte, 101)), // 101 characters
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "repository error - duplicate email",
			request: &CreateTeacherRequest{
				FirstName:  "Anna",
				LastName:   "Schmidt",
				Email:      "anna@example.com",
				Department: "Computer Science",
			},
			setupMock: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.Anything).
					Return(apperrors.Conflict("teacher with this email already exists"))
			},
			expectError: true,
		},
		{
			name: "repository error - internal error",
			request: &CreateTeacherRequest{
				FirstName:  "Anna",
				LastName:   "Schmidt",
				Email:      "anna@example.com",
				Department: "Computer Science",
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
			teacher, err := service.Create(context.Background(), tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, teacher)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teacher)
				assert.Equal(t, tt.request.FirstName, teacher.FirstName)
				assert.Equal(t, tt.request.LastName, teacher.LastName)
				assert.Equal(t, tt.request.Email, teacher.Email)
				assert.Equal(t, tt.request.Department, teacher.Department)
				assert.NotEmpty(t, teacher.UUID)
				assert.False(t, teacher.CreatedAt.IsZero())
				assert.False(t, teacher.UpdatedAt.IsZero())
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
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			setupMock: func(m *MockRepository) {
				teacher := createValidTeacher()
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440010").Return(teacher, nil)
			},
			expectError: false,
		},
		{
			name: "teacher not found",
			uuid: "nonexistent-uuid",
			setupMock: func(m *MockRepository) {
				m.On("GetByUUID", mock.Anything, "nonexistent-uuid").
					Return(nil, apperrors.NotFound("teacher not found"))
			},
			expectError: true,
		},
		{
			name: "repository error",
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			setupMock: func(m *MockRepository) {
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440010").
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
			teacher, err := service.GetByUUID(context.Background(), tt.uuid)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, teacher)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teacher)
				assert.Equal(t, tt.uuid, teacher.UUID)
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
		request     *UpdateTeacherRequest
		setupMock   func(*MockRepository)
		expectError bool
	}{
		{
			name: "successful update - all fields",
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			request: &UpdateTeacherRequest{
				FirstName:  "Maria",
				LastName:   "Müller",
				Email:      "maria@example.com",
				Department: "Mathematics",
			},
			setupMock: func(m *MockRepository) {
				teacher := createValidTeacher()
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440010").Return(teacher, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "successful update - partial fields",
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			request: &UpdateTeacherRequest{
				Department: "Physics",
			},
			setupMock: func(m *MockRepository) {
				teacher := createValidTeacher()
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440010").Return(teacher, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "validation error - invalid email",
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			request: &UpdateTeacherRequest{
				Email: "invalid-email",
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - first name too long",
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			request: &UpdateTeacherRequest{
				FirstName: string(make([]byte, 101)),
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - department too long",
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			request: &UpdateTeacherRequest{
				Department: string(make([]byte, 101)),
			},
			setupMock:   func(m *MockRepository) {},
			expectError: true,
		},
		{
			name: "teacher not found",
			uuid: "nonexistent-uuid",
			request: &UpdateTeacherRequest{
				FirstName: "Maria",
			},
			setupMock: func(m *MockRepository) {
				m.On("GetByUUID", mock.Anything, "nonexistent-uuid").
					Return(nil, apperrors.NotFound("teacher not found"))
			},
			expectError: true,
		},
		{
			name: "repository update error - duplicate email",
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			request: &UpdateTeacherRequest{
				Email: "duplicate@example.com",
			},
			setupMock: func(m *MockRepository) {
				teacher := createValidTeacher()
				m.On("GetByUUID", mock.Anything, "550e8400-e29b-41d4-a716-446655440010").Return(teacher, nil)
				m.On("Update", mock.Anything, mock.Anything).
					Return(apperrors.Conflict("teacher with this email already exists"))
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
			teacher, err := service.Update(context.Background(), tt.uuid, tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teacher)
				// Verify updated fields
				if tt.request.FirstName != "" {
					assert.Equal(t, tt.request.FirstName, teacher.FirstName)
				}
				if tt.request.LastName != "" {
					assert.Equal(t, tt.request.LastName, teacher.LastName)
				}
				if tt.request.Email != "" {
					assert.Equal(t, tt.request.Email, teacher.Email)
				}
				if tt.request.Department != "" {
					assert.Equal(t, tt.request.Department, teacher.Department)
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
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			setupMock: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "550e8400-e29b-41d4-a716-446655440010").Return(nil)
			},
			expectError: false,
		},
		{
			name: "teacher not found",
			uuid: "nonexistent-uuid",
			setupMock: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "nonexistent-uuid").
					Return(apperrors.NotFound("teacher not found"))
			},
			expectError: true,
		},
		{
			name: "repository error",
			uuid: "550e8400-e29b-41d4-a716-446655440010",
			setupMock: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "550e8400-e29b-41d4-a716-446655440010").
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
			name: "successful list - multiple teachers",
			setupMock: func(m *MockRepository) {
				teachers := []*Teacher{
					createValidTeacher(),
					{
						UUID:       "550e8400-e29b-41d4-a716-446655440011",
						FirstName:  "Thomas",
						LastName:   "Müller",
						Email:      "thomas@example.com",
						Department: "Mathematics",
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					},
				}
				m.On("List", mock.Anything).Return(teachers, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "successful list - empty",
			setupMock: func(m *MockRepository) {
				m.On("List", mock.Anything).Return([]*Teacher{}, nil)
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
			teachers, err := service.List(context.Background())

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, teachers)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teachers)
				assert.Len(t, teachers, tt.expectCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
