package auth

import (
	"context"
	"testing"
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/config"
	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RefreshToken), args.Error(1)
}

func (m *MockRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) UpdateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRepository) CreateInvitation(ctx context.Context, invitation *Invitation) error {
	args := m.Called(ctx, invitation)
	return args.Error(0)
}

func (m *MockRepository) GetInvitationByToken(ctx context.Context, token string) (*Invitation, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Invitation), args.Error(1)
}

func (m *MockRepository) MarkInvitationAsUsed(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRepository) DeleteExpiredInvitations(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRepository) CreateStudent(ctx context.Context, student *StudentRecord) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockRepository) CreateTeacher(ctx context.Context, teacher *TeacherRecord) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

// Helper function to create test JWT config
func getTestJWTConfig() *config.JWTConfig {
	return &config.JWTConfig{
		Secret:             "test-secret-key-for-testing",
		AccessTokenExpiry:  "1h",
		RefreshTokenExpiry: "720h",
		InvitationExpiry:   "168h",
		AdminKeyPath:       "./test-admin.key",
		AdminPubKeyPath:    "./test-admin.pub",
	}
}

func TestNewAuthService(t *testing.T) {
	t.Run("creates service successfully", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()

		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		assert.NotNil(t, service)
		assert.Equal(t, mockRepo, service.repo)
		assert.Equal(t, jwtService, service.jwtService)
		assert.Equal(t, jwtConfig, service.jwtConfig)
	})
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()

	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		// Create test user with hashed password
		password := "testPassword123"
		hashedPassword, _ := crypto.HashPassword(password)
		user := &User{
			ID:           "user-123",
			Username:     "testuser",
			PasswordHash: hashedPassword,
			Role:         crypto.RoleStudent,
		}

		// Mock expectations
		mockRepo.On("GetUserByUsername", ctx, "testuser").Return(user, nil)
		mockRepo.On("CreateRefreshToken", ctx, mock.AnythingOfType("*auth.RefreshToken")).Return(nil)
		mockRepo.On("UpdateLastLogin", ctx, "user-123").Return(nil)

		// Execute login
		req := &LoginRequest{
			Username: "testuser",
			Password: password,
		}
		response, err := service.Login(ctx, req)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Equal(t, "Bearer", response.TokenType)
		assert.Equal(t, string(crypto.RoleStudent), response.Role)
		assert.Equal(t, "user-123", response.UserID)
		assert.Greater(t, response.ExpiresIn, int64(0))

		mockRepo.AssertExpectations(t)
	})

	t.Run("login with wrong password", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		hashedPassword, _ := crypto.HashPassword("correctPassword")
		user := &User{
			ID:           "user-123",
			Username:     "testuser",
			PasswordHash: hashedPassword,
			Role:         crypto.RoleStudent,
		}

		mockRepo.On("GetUserByUsername", ctx, "testuser").Return(user, nil)

		req := &LoginRequest{
			Username: "testuser",
			Password: "wrongPassword",
		}
		response, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid credentials")

		mockRepo.AssertExpectations(t)
	})

	t.Run("login with non-existent user", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		mockRepo.On("GetUserByUsername", ctx, "nonexistent").Return(nil, assert.AnError)

		req := &LoginRequest{
			Username: "nonexistent",
			Password: "anyPassword",
		}
		response, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid credentials")

		mockRepo.AssertExpectations(t)
	})

	t.Run("login with different roles", func(t *testing.T) {
		roles := []crypto.Role{crypto.RoleStudent, crypto.RoleTeacher, crypto.RoleAdmin}

		for _, role := range roles {
			mockRepo := new(MockRepository)
			jwtService := crypto.NewJWTService("test-secret")
			jwtConfig := getTestJWTConfig()
			service := NewAuthService(mockRepo, jwtService, jwtConfig)

			password := "testPassword123"
			hashedPassword, _ := crypto.HashPassword(password)
			user := &User{
				ID:           "user-123",
				Username:     "testuser",
				PasswordHash: hashedPassword,
				Role:         role,
			}

			mockRepo.On("GetUserByUsername", ctx, "testuser").Return(user, nil)
			mockRepo.On("CreateRefreshToken", ctx, mock.AnythingOfType("*auth.RefreshToken")).Return(nil)
			mockRepo.On("UpdateLastLogin", ctx, "user-123").Return(nil)

			req := &LoginRequest{
				Username: "testuser",
				Password: password,
			}
			response, err := service.Login(ctx, req)

			require.NoError(t, err)
			assert.Equal(t, string(role), response.Role)

			mockRepo.AssertExpectations(t)
		}
	})
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	ctx := context.Background()

	t.Run("successful token refresh", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		// Generate a valid refresh token
		userID := "user-123"
		refreshToken, _ := jwtService.GenerateToken(userID, crypto.RoleStudent, crypto.TokenTypeRefresh, 24*time.Hour)

		// Mock data
		refreshTokenRecord := &RefreshToken{
			UserID:    userID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		user := &User{
			ID:       userID,
			Username: "testuser",
			Role:     crypto.RoleStudent,
		}

		mockRepo.On("GetRefreshToken", ctx, refreshToken).Return(refreshTokenRecord, nil)
		mockRepo.On("GetUserByID", ctx, userID).Return(user, nil)

		// Execute refresh
		response, err := service.RefreshAccessToken(ctx, refreshToken)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.AccessToken)
		assert.Equal(t, refreshToken, response.RefreshToken) // Same refresh token
		assert.Equal(t, "Bearer", response.TokenType)
		assert.Equal(t, userID, response.UserID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("refresh with invalid token", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		response, err := service.RefreshAccessToken(ctx, "invalid-token")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid refresh token")

		mockRepo.AssertNotCalled(t, "GetRefreshToken")
	})

	t.Run("refresh with wrong token type", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		// Generate an access token instead of refresh token
		accessToken, _ := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, time.Hour)

		response, err := service.RefreshAccessToken(ctx, accessToken)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid token type")
	})

	t.Run("refresh with expired token", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		// Generate an expired refresh token
		userID := "user-123"
		refreshToken, _ := jwtService.GenerateToken(userID, crypto.RoleStudent, crypto.TokenTypeRefresh, 24*time.Hour)

		// Mock expired token in database
		refreshTokenRecord := &RefreshToken{
			UserID:    userID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
		}

		mockRepo.On("GetRefreshToken", ctx, refreshToken).Return(refreshTokenRecord, nil)
		mockRepo.On("DeleteRefreshToken", ctx, refreshToken).Return(nil)

		response, err := service.RefreshAccessToken(ctx, refreshToken)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "expired")

		mockRepo.AssertExpectations(t)
	})

	t.Run("refresh with token not in database", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		refreshToken, _ := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeRefresh, 24*time.Hour)

		mockRepo.On("GetRefreshToken", ctx, refreshToken).Return(nil, assert.AnError)

		response, err := service.RefreshAccessToken(ctx, refreshToken)

		assert.Error(t, err)
		assert.Nil(t, response)

		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Logout(t *testing.T) {
	ctx := context.Background()

	t.Run("successful logout", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		userID := "user-123"
		mockRepo.On("DeleteUserRefreshTokens", ctx, userID).Return(nil)

		err := service.Logout(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("logout with database error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		jwtService := crypto.NewJWTService("test-secret")
		jwtConfig := getTestJWTConfig()
		service := NewAuthService(mockRepo, jwtService, jwtConfig)

		userID := "user-123"
		mockRepo.On("DeleteUserRefreshTokens", ctx, userID).Return(assert.AnError)

		err := service.Logout(ctx, userID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	jwtService := crypto.NewJWTService("test-secret")
	jwtConfig := getTestJWTConfig()
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo, jwtService, jwtConfig)

	t.Run("validates access token successfully", func(t *testing.T) {
		token, _ := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, time.Hour)

		claims, err := service.ValidateToken(token)

		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "user-123", claims.UserID)
		assert.Equal(t, crypto.RoleStudent, claims.Role)
		assert.Equal(t, crypto.TokenTypeAccess, claims.TokenType)
	})

	t.Run("validates admin token successfully", func(t *testing.T) {
		token, _ := jwtService.GenerateToken("admin-001", crypto.RoleAdmin, crypto.TokenTypeAdmin, 100*365*24*time.Hour)

		claims, err := service.ValidateToken(token)

		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, crypto.TokenTypeAdmin, claims.TokenType)
	})

	t.Run("rejects refresh token", func(t *testing.T) {
		token, _ := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeRefresh, time.Hour)

		claims, err := service.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token type")
	})

	t.Run("rejects invitation token", func(t *testing.T) {
		token, _ := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeInvitation, time.Hour)

		claims, err := service.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token type")
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		claims, err := service.ValidateToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("rejects expired token", func(t *testing.T) {
		token, _ := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, -time.Hour)

		claims, err := service.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

// Benchmark tests
func BenchmarkAuthService_Login(b *testing.B) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	jwtService := crypto.NewJWTService("test-secret")
	jwtConfig := getTestJWTConfig()
	service := NewAuthService(mockRepo, jwtService, jwtConfig)

	password := "testPassword123"
	hashedPassword, _ := crypto.HashPassword(password)
	user := &User{
		ID:           "user-123",
		Username:     "testuser",
		PasswordHash: hashedPassword,
		Role:         crypto.RoleStudent,
	}

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil)
	mockRepo.On("CreateRefreshToken", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("UpdateLastLogin", mock.Anything, "user-123").Return(nil)

	req := &LoginRequest{
		Username: "testuser",
		Password: password,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Login(ctx, req)
	}
}
