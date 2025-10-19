package auth

import (
	"time"

	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
)

// User represents an authenticated user in the system
// @Description User authentication information
type User struct {
	ID           string      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username     string      `json:"username" gorm:"uniqueIndex;not null" example:"max.mustermann" validate:"required,min=3,max=100"`
	PasswordHash string      `json:"-" gorm:"not null"` // Never expose password hash in JSON
	Role         crypto.Role `json:"role" gorm:"not null" example:"student"`
	UserUUID     *string     `json:"user_uuid,omitempty" gorm:"type:uuid;index" example:"550e8400-e29b-41d4-a716-446655440000"` // References student or teacher (NULL for admin)
	CreatedAt    time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	LastLogin    *time.Time  `json:"last_login,omitempty"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// RefreshToken represents a refresh token for extending sessions
// @Description Refresh token for session management
type RefreshToken struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// Invitation represents an invitation for a user to register
// @Description Invitation token for user registration
type Invitation struct {
	ID          string      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Token       string      `json:"token" gorm:"uniqueIndex;not null"`
	Email       string      `json:"email" gorm:"not null"`
	FirstName   string      `json:"first_name" gorm:"not null"`
	LastName    string      `json:"last_name" gorm:"not null"`
	Role        crypto.Role `json:"role" gorm:"not null"`
	Department  *string     `json:"department,omitempty"`                          // Only for teachers
	TeacherUUID *string     `json:"teacher_uuid,omitempty" gorm:"type:uuid;index"` // Required for students, NULL for teachers
	Used        bool        `json:"used" gorm:"default:false;index"`
	CreatedAt   time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	ExpiresAt   time.Time   `json:"expires_at" gorm:"not null"`
}

// TableName specifies the table name for GORM
func (Invitation) TableName() string {
	return "invitations"
}

// IsExpired checks if the invitation has expired
func (i *Invitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// IsValid checks if the invitation is valid (not used and not expired)
func (i *Invitation) IsValid() bool {
	return !i.Used && !i.IsExpired()
}

// DTO Models for API requests/responses

// LoginRequest represents a login request
// @Description Login credentials
type LoginRequest struct {
	Username string `json:"username" example:"max.mustermann" validate:"required,min=3"`
	Password string `json:"password" example:"SecurePassword123!" validate:"required,min=8"`
}

// LoginResponse represents a successful login response
// @Description Login response with tokens
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGc..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGc..."`
	ExpiresIn    int64  `json:"expires_in" example:"3600"` // Seconds
	TokenType    string `json:"token_type" example:"Bearer"`
	Role         string `json:"role" example:"student"`
	UserID       string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// RefreshTokenRequest represents a token refresh request
// @Description Refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGc..." validate:"required"`
}

// CompleteInvitationRequest represents a request to complete an invitation
// @Description Complete invitation with credentials
type CompleteInvitationRequest struct {
	Username string `json:"username" example:"max.mustermann" validate:"required,min=3,max=100"`
	Password string `json:"password" example:"SecurePassword123!" validate:"required,min=8"`
}

// InvitationResponse represents invitation details
// @Description Invitation information
type InvitationResponse struct {
	Email       string  `json:"email" example:"max@example.com"`
	FirstName   string  `json:"first_name" example:"Max"`
	LastName    string  `json:"last_name" example:"Mustermann"`
	Role        string  `json:"role" example:"student"`
	Department  *string `json:"department,omitempty" example:"IT"`
	TeacherUUID *string `json:"teacher_uuid,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	ExpiresAt   string  `json:"expires_at" example:"2025-10-25T12:00:00Z"`
}

// CreateInvitationRequest represents a request to create an invitation
// @Description Create invitation for new user
type CreateInvitationRequest struct {
	Email       string  `json:"email" example:"max@example.com" validate:"required,email"`
	FirstName   string  `json:"first_name" example:"Max" validate:"required,min=1,max=100"`
	LastName    string  `json:"last_name" example:"Mustermann" validate:"required,min=1,max=100"`
	Role        string  `json:"role" example:"student" validate:"required,oneof=student teacher"`
	Department  *string `json:"department,omitempty" example:"IT" validate:"omitempty,min=1,max=100"`                            // Required for teachers
	TeacherUUID *string `json:"teacher_uuid,omitempty" example:"550e8400-e29b-41d4-a716-446655440000" validate:"omitempty,uuid"` // Required for students
}

// CreateInvitationResponse represents the created invitation
// @Description Created invitation with link
type CreateInvitationResponse struct {
	InvitationToken string `json:"invitation_token" example:"eyJhbGc..."`
	InvitationLink  string `json:"invitation_link" example:"https://fits.example.com/invite/eyJhbGc..."`
	ExpiresAt       string `json:"expires_at" example:"2025-10-25T12:00:00Z"`
}

// BootstrapResponse represents the bootstrap initialization response
// @Description Admin bootstrap response
type BootstrapResponse struct {
	AdminToken    string `json:"admin_token" example:"eyJhbGc..."`
	Message       string `json:"message" example:"Admin certificate generated successfully"`
	PublicKeyPath string `json:"public_key_path" example:"./configs/keys/admin.pub"`
}
