// Package auth provides authentication and authorization services for the FITS backend.
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/config"
	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/JustDoItBetter/FITS-backend/pkg/logger"
	"go.uber.org/zap"
)

// AuthService handles authentication operations
type AuthService struct {
	repo       Repository
	jwtService *crypto.JWTService
	jwtConfig  *config.JWTConfig
}

// NewAuthService creates a new auth service
func NewAuthService(repo Repository, jwtService *crypto.JWTService, jwtConfig *config.JWTConfig) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
		jwtConfig:  jwtConfig,
	}
}

// Login authenticates a user and returns access and refresh tokens
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get user by username
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.Unauthorized("invalid credentials")
	}

	// Verify password
	if err := crypto.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errors.Unauthorized("invalid credentials")
	}

	// Generate access token
	accessToken, err := s.jwtService.GenerateToken(
		user.ID,
		user.Role,
		crypto.TokenTypeAccess,
		s.jwtConfig.GetAccessTokenExpiry(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshTokenString, err := s.jwtService.GenerateToken(
		user.ID,
		user.Role,
		crypto.TokenTypeRefresh,
		s.jwtConfig.GetRefreshTokenExpiry(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token to DB
	refreshToken := &RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(s.jwtConfig.GetRefreshTokenExpiry()),
	}

	if err := s.repo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Update last login (non-critical, but log errors for monitoring)
	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		logger.Warn("Failed to update last login timestamp",
			zap.String("user_id", user.ID),
			zap.Error(err),
		)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.jwtConfig.GetAccessTokenExpiry().Seconds()),
		TokenType:    "Bearer",
		Role:         string(user.Role),
		UserID:       user.ID,
	}, nil
}

// RefreshAccessToken generates a new access token using a refresh token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (*LoginResponse, error) {
	// Validate refresh token JWT
	claims, err := s.jwtService.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, errors.Unauthorized("invalid refresh token")
	}

	// Check token type
	if claims.TokenType != crypto.TokenTypeRefresh {
		return nil, errors.Unauthorized("invalid token type")
	}

	// Check if refresh token exists in DB
	refreshToken, err := s.repo.GetRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, errors.Unauthorized("invalid refresh token")
	}

	// Check if expired
	if refreshToken.IsExpired() {
		// Attempt to delete expired token from database
		// Error is logged but doesn't block the rejection - security over cleanup
		if err := s.repo.DeleteRefreshToken(ctx, refreshTokenString); err != nil {
			logger.Error("Failed to delete expired refresh token",
				zap.String("user_id", refreshToken.UserID),
				zap.Error(err),
			)
		}
		return nil, errors.Unauthorized("refresh token expired")
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, err := s.jwtService.GenerateToken(
		user.ID,
		user.Role,
		crypto.TokenTypeAccess,
		s.jwtConfig.GetAccessTokenExpiry(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString, // Keep same refresh token
		ExpiresIn:    int64(s.jwtConfig.GetAccessTokenExpiry().Seconds()),
		TokenType:    "Bearer",
		Role:         string(user.Role),
		UserID:       user.ID,
	}, nil
}

// Logout logs out a user by deleting all their refresh tokens
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	// Delete all refresh tokens for user
	return s.repo.DeleteUserRefreshTokens(ctx, userID)
}

// ValidateToken validates an access token and returns the claims
func (s *AuthService) ValidateToken(token string) (*crypto.Claims, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, errors.Unauthorized("invalid token")
	}

	// Check token type (should be access token)
	if claims.TokenType != crypto.TokenTypeAccess && claims.TokenType != crypto.TokenTypeAdmin {
		return nil, errors.Unauthorized("invalid token type")
	}

	return claims, nil
}
