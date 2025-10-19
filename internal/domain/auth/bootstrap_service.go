package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/config"
	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
)

// BootstrapService handles admin initialization
type BootstrapService struct {
	repo      Repository
	jwtConfig *config.JWTConfig
}

// NewBootstrapService creates a new bootstrap service
func NewBootstrapService(repo Repository, jwtConfig *config.JWTConfig) *BootstrapService {
	return &BootstrapService{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

// InitializeAdmin initializes the admin certificate and returns an admin token
// This should only be called once during installation
func (s *BootstrapService) InitializeAdmin(ctx context.Context) (*BootstrapResponse, error) {
	// 1. Check if admin already exists
	existingAdmin, _ := s.repo.GetUserByUsername(ctx, "admin")
	if existingAdmin != nil {
		return nil, fmt.Errorf("admin already initialized")
	}

	// 2. Create keys directory if not exists
	keyDir := "configs/keys"
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create keys directory: %w", err)
	}

	// 3. Generate RSA Keypair
	keypair, err := crypto.GenerateRSAKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	// 4. Save keys to files
	if err := crypto.SavePrivateKeyToFile(keypair.PrivateKey, s.jwtConfig.AdminKeyPath); err != nil {
		return nil, fmt.Errorf("failed to save private key: %w", err)
	}

	if err := crypto.SavePublicKeyToFile(keypair.PublicKey, s.jwtConfig.AdminPubKeyPath); err != nil {
		return nil, fmt.Errorf("failed to save public key: %w", err)
	}

	// 5. Create admin user in database
	adminUser := &User{
		Username:     "admin",
		PasswordHash: "not-used",
		Role:         crypto.RoleAdmin,
	}

	if err := s.repo.CreateUser(ctx, adminUser); err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	// 6. Generate admin JWT (very long expiry - 100 years)
	jwtService := crypto.NewJWTService(s.jwtConfig.Secret)
	adminToken, err := jwtService.GenerateToken(
		adminUser.ID,
		crypto.RoleAdmin,
		crypto.TokenTypeAdmin,
		100*365*24*time.Hour,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate admin token: %w", err)
	}

	return &BootstrapResponse{
		AdminToken:    adminToken,
		Message:       "Admin certificate generated successfully. Store this token securely - it cannot be recovered!",
		PublicKeyPath: s.jwtConfig.AdminPubKeyPath,
	}, nil
}

// IsBootstrapped checks if the system has been bootstrapped
func (s *BootstrapService) IsBootstrapped(ctx context.Context) (bool, error) {
	admin, err := s.repo.GetUserByUsername(ctx, "admin")
	if err != nil {
		return false, nil
	}
	return admin != nil, nil
}
