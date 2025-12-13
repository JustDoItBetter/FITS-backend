package crypto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType represents the type of JWT token
type TokenType string

const (
	// TokenTypeAccess represents an access token for API authentication.
	TokenTypeAccess TokenType = "access"
	// TokenTypeRefresh represents a refresh token for obtaining new access tokens.
	TokenTypeRefresh TokenType = "refresh"
	// TokenTypeInvitation represents an invitation token for user registration.
	TokenTypeInvitation TokenType = "invitation"
	// TokenTypeAdmin represents an admin token with elevated privileges.
	TokenTypeAdmin TokenType = "admin"
)

// Role represents a user role
type Role string

const (
	// RoleAdmin represents an administrator user with full system access.
	RoleAdmin Role = "admin"
	// RoleTeacher represents a teacher user with access to teacher resources.
	RoleTeacher Role = "teacher"
	// RoleStudent represents a student user with limited access.
	RoleStudent Role = "student"
)

// Claims represents JWT claims
type Claims struct {
	UserID    string    `json:"sub"`
	Role      Role      `json:"role"`
	TokenType TokenType `json:"type"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secret []byte
}

// NewJWTService creates a new JWT service
func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
	}
}

// GenerateToken generates a new JWT token
func (s *JWTService) GenerateToken(userID string, role Role, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID:    userID,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns its claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// ExtractUserID extracts the user ID from a token without full validation
// Useful for logging/debugging purposes
func (s *JWTService) ExtractUserID(tokenString string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, nil
}
