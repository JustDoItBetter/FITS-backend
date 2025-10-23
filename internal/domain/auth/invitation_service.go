package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/validation"
	"github.com/JustDoItBetter/FITS-backend/internal/config"
	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
)

// InvitationService handles user invitation operations
type InvitationService struct {
	repo       Repository
	jwtService *crypto.JWTService
	jwtConfig  *config.JWTConfig
}

// NewInvitationService creates a new invitation service
func NewInvitationService(repo Repository, jwtService *crypto.JWTService, jwtConfig *config.JWTConfig) *InvitationService {
	return &InvitationService{
		repo:       repo,
		jwtService: jwtService,
		jwtConfig:  jwtConfig,
	}
}

// CreateInvitation creates a new invitation for a user
func (s *InvitationService) CreateInvitation(ctx context.Context, req *CreateInvitationRequest) (*CreateInvitationResponse, error) {
	// Validate role
	var role crypto.Role
	if req.Role == "student" {
		role = crypto.RoleStudent
		// Validate that student has teacher_uuid
		if req.TeacherUUID == nil || *req.TeacherUUID == "" {
			return nil, errors.ValidationError("teacher_uuid is required for students")
		}
		// TODO: Validate that teacher exists in database
	} else if req.Role == "teacher" {
		role = crypto.RoleTeacher
		// Validate that teacher has department
		if req.Department == nil || *req.Department == "" {
			return nil, errors.ValidationError("department is required for teachers")
		}
	} else {
		return nil, errors.ValidationError("role must be 'student' or 'teacher'")
	}

	// Check if email already used in another invitation
	// (Optional: prevent duplicate invitations)

	// Generate invitation token
	invitationToken, err := s.jwtService.GenerateToken(
		req.Email, // Use email as subject
		role,
		crypto.TokenTypeInvitation,
		s.jwtConfig.GetInvitationExpiry(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invitation token: %w", err)
	}

	// Create invitation record with user data
	invitation := &Invitation{
		Token:       invitationToken,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Role:        role,
		Department:  req.Department,
		TeacherUUID: req.TeacherUUID,
		Used:        false,
		ExpiresAt:   time.Now().Add(s.jwtConfig.GetInvitationExpiry()),
	}

	if err := s.repo.CreateInvitation(ctx, invitation); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	// Generate invitation link
	// TODO: Get base URL from config
	invitationLink := fmt.Sprintf("https://fits.example.com/invite/%s", invitationToken)

	return &CreateInvitationResponse{
		InvitationToken: invitationToken,
		InvitationLink:  invitationLink,
		ExpiresAt:       invitation.ExpiresAt.Format(time.RFC3339),
	}, nil
}

// GetInvitationDetails retrieves invitation details by token
func (s *InvitationService) GetInvitationDetails(ctx context.Context, token string) (*InvitationResponse, error) {
	// Validate JWT token
	_, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, errors.Unauthorized("invalid invitation token")
	}

	// Get invitation from DB
	invitation, err := s.repo.GetInvitationByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Check if invitation is valid
	if !invitation.IsValid() {
		if invitation.Used {
			return nil, errors.BadRequest("invitation already used")
		}
		return nil, errors.BadRequest("invitation expired")
	}

	// Return invitation data
	return &InvitationResponse{
		Email:       invitation.Email,
		FirstName:   invitation.FirstName,
		LastName:    invitation.LastName,
		Role:        string(invitation.Role),
		Department:  invitation.Department,
		TeacherUUID: invitation.TeacherUUID,
		ExpiresAt:   invitation.ExpiresAt.Format(time.RFC3339),
	}, nil
}

// CompleteInvitation completes an invitation by creating a user account
// All operations are executed within a database transaction to ensure data consistency
func (s *InvitationService) CompleteInvitation(ctx context.Context, token string, req *CompleteInvitationRequest) error {
	// Get invitation from DB (outside transaction for read-only operation)
	invitation, err := s.repo.GetInvitationByToken(ctx, token)
	if err != nil {
		return err
	}

	// Validate invitation
	if !invitation.IsValid() {
		if invitation.Used {
			return errors.BadRequest("invitation already used")
		}
		return errors.BadRequest("invitation expired")
	}

	// Check if username already exists
	existingUser, _ := s.repo.GetUserByUsername(ctx, req.Username)
	if existingUser != nil {
		return errors.Conflict("username already exists")
	}

	// Validate password strength before hashing
	if err := validation.ValidatePasswordStrength(req.Password); err != nil {
		return err
	}

	// Check against common passwords
	if validation.IsCommonPassword(req.Password) {
		return errors.ValidationError("password is too common and easily guessable, please choose a stronger password")
	}

	// Hash password
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate UUID for Student/Teacher
	entityUUID := uuid.New().String()

	// Execute all write operations within a transaction
	// If any operation fails, all changes are rolled back automatically
	return s.repo.ExecuteInTransaction(ctx, func(txRepo Repository) error {
		// Create Student or Teacher record
		if invitation.Role == crypto.RoleStudent {
			student := &StudentRecord{
				ID:        entityUUID,
				FirstName: invitation.FirstName,
				LastName:  invitation.LastName,
				Email:     invitation.Email,
				TeacherID: invitation.TeacherUUID, // Assigned from invitation
			}
			if err := txRepo.CreateStudent(ctx, student); err != nil {
				return fmt.Errorf("failed to create student: %w", err)
			}
		} else if invitation.Role == crypto.RoleTeacher {
			if invitation.Department == nil {
				return errors.ValidationError("department is required for teacher")
			}
			teacher := &TeacherRecord{
				ID:         entityUUID,
				FirstName:  invitation.FirstName,
				LastName:   invitation.LastName,
				Email:      invitation.Email,
				Department: *invitation.Department,
			}
			if err := txRepo.CreateTeacher(ctx, teacher); err != nil {
				return fmt.Errorf("failed to create teacher: %w", err)
			}
		}

		// Create user with reference to Student/Teacher UUID
		user := &User{
			Username:     req.Username,
			PasswordHash: passwordHash,
			Role:         invitation.Role,
			UserUUID:     &entityUUID,
		}

		if err := txRepo.CreateUser(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Mark invitation as used
		if err := txRepo.MarkInvitationAsUsed(ctx, token); err != nil {
			return fmt.Errorf("failed to mark invitation as used: %w", err)
		}

		// TODO: If teacher, generate RSA keypair (implement in keypair domain first)
		// This will be added later when keypair domain is implemented

		// Transaction commits automatically if no error is returned
		return nil
	})
}
