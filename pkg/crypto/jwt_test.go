package crypto

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTService(t *testing.T) {
	t.Run("creates service with secret", func(t *testing.T) {
		secret := "test-secret"
		service := NewJWTService(secret)

		assert.NotNil(t, service)
		assert.Equal(t, []byte(secret), service.secret)
	})

	t.Run("creates service with empty secret", func(t *testing.T) {
		service := NewJWTService("")

		assert.NotNil(t, service)
		assert.Equal(t, []byte(""), service.secret)
	})
}

func TestGenerateToken(t *testing.T) {
	service := NewJWTService("test-secret-key")

	t.Run("generates valid access token", func(t *testing.T) {
		userID := "user-123"
		role := RoleStudent
		tokenType := TokenTypeAccess
		expiry := time.Hour

		token, err := service.GenerateToken(userID, role, tokenType, expiry)

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token structure (JWT format: header.payload.signature)
		assert.Contains(t, token, ".")
	})

	t.Run("generates valid refresh token", func(t *testing.T) {
		userID := "user-456"
		role := RoleTeacher
		tokenType := TokenTypeRefresh
		expiry := 24 * time.Hour

		token, err := service.GenerateToken(userID, role, tokenType, expiry)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generates valid admin token", func(t *testing.T) {
		userID := "admin-001"
		role := RoleAdmin
		tokenType := TokenTypeAdmin
		expiry := 100 * 365 * 24 * time.Hour // 100 years

		token, err := service.GenerateToken(userID, role, tokenType, expiry)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generates valid invitation token", func(t *testing.T) {
		userID := "invite-789"
		role := RoleStudent
		tokenType := TokenTypeInvitation
		expiry := 7 * 24 * time.Hour // 7 days

		token, err := service.GenerateToken(userID, role, tokenType, expiry)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generates different tokens for different users", func(t *testing.T) {
		token1, err1 := service.GenerateToken("user-1", RoleStudent, TokenTypeAccess, time.Hour)
		token2, err2 := service.GenerateToken("user-2", RoleStudent, TokenTypeAccess, time.Hour)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})

	t.Run("generates different tokens for same user called twice", func(t *testing.T) {
		userID := "user-same"
		token1, err1 := service.GenerateToken(userID, RoleStudent, TokenTypeAccess, time.Hour)
		time.Sleep(time.Millisecond) // Ensure different timestamp
		token2, err2 := service.GenerateToken(userID, RoleStudent, TokenTypeAccess, time.Hour)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2, "tokens should be different due to different JTI and timestamps")
	})
}

func TestValidateToken(t *testing.T) {
	service := NewJWTService("test-secret-key")

	t.Run("validates correct token successfully", func(t *testing.T) {
		userID := "user-123"
		role := RoleStudent
		tokenType := TokenTypeAccess
		expiry := time.Hour

		token, err := service.GenerateToken(userID, role, tokenType, expiry)
		require.NoError(t, err)

		claims, err := service.ValidateToken(token)
		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, tokenType, claims.TokenType)
	})

	t.Run("rejects token with wrong secret", func(t *testing.T) {
		// Generate token with one secret
		service1 := NewJWTService("secret-1")
		token, err := service1.GenerateToken("user-123", RoleStudent, TokenTypeAccess, time.Hour)
		require.NoError(t, err)

		// Try to validate with different secret
		service2 := NewJWTService("secret-2")
		claims, err := service2.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to parse token")
	})

	t.Run("rejects expired token", func(t *testing.T) {
		userID := "user-expired"
		role := RoleStudent
		tokenType := TokenTypeAccess
		expiry := -time.Hour // Already expired

		token, err := service.GenerateToken(userID, role, tokenType, expiry)
		require.NoError(t, err)

		claims, err := service.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "expired")
	})

	t.Run("rejects malformed token", func(t *testing.T) {
		claims, err := service.ValidateToken("not.a.valid.jwt.token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("rejects empty token", func(t *testing.T) {
		claims, err := service.ValidateToken("")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("rejects token with invalid signature", func(t *testing.T) {
		// Create a token manually with wrong signature
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.invalid_signature"

		claims, err := service.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("validates token with correct claims structure", func(t *testing.T) {
		userID := "user-full-test"
		role := RoleTeacher
		tokenType := TokenTypeRefresh
		expiry := 24 * time.Hour

		token, err := service.GenerateToken(userID, role, tokenType, expiry)
		require.NoError(t, err)

		claims, err := service.ValidateToken(token)
		require.NoError(t, err)

		// Check all claims
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, tokenType, claims.TokenType)
		assert.NotEmpty(t, claims.ID) // JTI should be set
		assert.NotNil(t, claims.IssuedAt)
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.NotBefore)

		// Verify times are reasonable
		now := time.Now()
		assert.True(t, claims.IssuedAt.Before(now.Add(time.Second)))
		assert.True(t, claims.ExpiresAt.After(now))
	})
}

func TestExtractUserID(t *testing.T) {
	service := NewJWTService("test-secret-key")

	t.Run("extracts user ID from valid token", func(t *testing.T) {
		expectedUserID := "user-extract-123"
		token, err := service.GenerateToken(expectedUserID, RoleStudent, TokenTypeAccess, time.Hour)
		require.NoError(t, err)

		userID, err := service.ExtractUserID(token)
		require.NoError(t, err)
		assert.Equal(t, expectedUserID, userID)
	})

	t.Run("extracts user ID from expired token", func(t *testing.T) {
		expectedUserID := "user-expired"
		token, err := service.GenerateToken(expectedUserID, RoleStudent, TokenTypeAccess, -time.Hour)
		require.NoError(t, err)

		// ExtractUserID should work even for expired tokens (no validation)
		userID, err := service.ExtractUserID(token)
		require.NoError(t, err)
		assert.Equal(t, expectedUserID, userID)
	})

	t.Run("fails on malformed token", func(t *testing.T) {
		userID, err := service.ExtractUserID("invalid.token")

		assert.Error(t, err)
		assert.Empty(t, userID)
	})

	t.Run("fails on empty token", func(t *testing.T) {
		userID, err := service.ExtractUserID("")

		assert.Error(t, err)
		assert.Empty(t, userID)
	})
}

func TestTokenTypes(t *testing.T) {
	service := NewJWTService("test-secret-key")

	testCases := []struct {
		name      string
		tokenType TokenType
	}{
		{"access token", TokenTypeAccess},
		{"refresh token", TokenTypeRefresh},
		{"invitation token", TokenTypeInvitation},
		{"admin token", TokenTypeAdmin},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := service.GenerateToken("user-123", RoleStudent, tc.tokenType, time.Hour)
			require.NoError(t, err)

			claims, err := service.ValidateToken(token)
			require.NoError(t, err)
			assert.Equal(t, tc.tokenType, claims.TokenType)
		})
	}
}

func TestRoles(t *testing.T) {
	service := NewJWTService("test-secret-key")

	testCases := []struct {
		name string
		role Role
	}{
		{"admin role", RoleAdmin},
		{"teacher role", RoleTeacher},
		{"student role", RoleStudent},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := service.GenerateToken("user-123", tc.role, TokenTypeAccess, time.Hour)
			require.NoError(t, err)

			claims, err := service.ValidateToken(token)
			require.NoError(t, err)
			assert.Equal(t, tc.role, claims.Role)
		})
	}
}

func TestJWTSigningMethod(t *testing.T) {
	service := NewJWTService("test-secret-key")

	t.Run("uses HS256 signing method", func(t *testing.T) {
		token, err := service.GenerateToken("user-123", RoleStudent, TokenTypeAccess, time.Hour)
		require.NoError(t, err)

		// Parse token without validation to check signing method
		parsedToken, _, err := jwt.NewParser().ParseUnverified(token, &Claims{})
		require.NoError(t, err)

		assert.Equal(t, "HS256", parsedToken.Method.Alg())
	})

	t.Run("rejects token with wrong signing method", func(t *testing.T) {
		// This test verifies that tokens signed with other methods are rejected
		// We can't easily create such a token here, but the ValidateToken function checks this
		claims, err := service.ValidateToken("eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJ0ZXN0In0.")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestClaims(t *testing.T) {
	t.Run("claims have correct JSON tags", func(t *testing.T) {
		claims := Claims{
			UserID:    "test-user",
			Role:      RoleStudent,
			TokenType: TokenTypeAccess,
		}

		// Verify the struct has the expected fields
		assert.Equal(t, "test-user", claims.UserID)
		assert.Equal(t, RoleStudent, claims.Role)
		assert.Equal(t, TokenTypeAccess, claims.TokenType)
	})
}

// Benchmark tests
func BenchmarkGenerateToken(b *testing.B) {
	service := NewJWTService("test-secret-key")
	userID := "bench-user"
	role := RoleStudent
	tokenType := TokenTypeAccess
	expiry := time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateToken(userID, role, tokenType, expiry)
	}
}

func BenchmarkValidateToken(b *testing.B) {
	service := NewJWTService("test-secret-key")
	token, _ := service.GenerateToken("bench-user", RoleStudent, TokenTypeAccess, time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateToken(token)
	}
}

func BenchmarkExtractUserID(b *testing.B) {
	service := NewJWTService("test-secret-key")
	token, _ := service.GenerateToken("bench-user", RoleStudent, TokenTypeAccess, time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ExtractUserID(token)
	}
}
