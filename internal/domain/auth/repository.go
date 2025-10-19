package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"gorm.io/gorm"
)

// Repository defines the interface for auth data access
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	UpdateLastLogin(ctx context.Context, userID string) error

	// Refresh token operations
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
	DeleteUserRefreshTokens(ctx context.Context, userID string) error

	// Invitation operations
	CreateInvitation(ctx context.Context, invitation *Invitation) error
	GetInvitationByToken(ctx context.Context, token string) (*Invitation, error)
	MarkInvitationAsUsed(ctx context.Context, token string) error
	DeleteExpiredInvitations(ctx context.Context) error

	// Student/Teacher operations (for invitation completion)
	CreateStudent(ctx context.Context, student *StudentRecord) error
	CreateTeacher(ctx context.Context, teacher *TeacherRecord) error

	// Transaction support
	// ExecuteInTransaction runs the given function within a database transaction
	// The function receives a Repository instance that uses the transaction
	// If the function returns an error, the transaction is rolled back
	ExecuteInTransaction(ctx context.Context, fn func(Repository) error) error
}

// StudentRecord represents a student record in the database
type StudentRecord struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	TeacherID *string
}

// TeacherRecord represents a teacher record in the database
type TeacherRecord struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	Department string
}

// GormRepository implements Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM repository
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// User operations

func (r *GormRepository) CreateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *GormRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("user")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *GormRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("user")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *GormRepository) UpdateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *GormRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("last_login", now).Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// Refresh token operations

func (r *GormRepository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

func (r *GormRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("refresh token")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return &refreshToken, nil
}

func (r *GormRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

func (r *GormRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}
	return nil
}

func (r *GormRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete user refresh tokens: %w", err)
	}
	return nil
}

// Invitation operations

func (r *GormRepository) CreateInvitation(ctx context.Context, invitation *Invitation) error {
	if err := r.db.WithContext(ctx).Create(invitation).Error; err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}
	return nil
}

func (r *GormRepository) GetInvitationByToken(ctx context.Context, token string) (*Invitation, error) {
	var invitation Invitation
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("invitation")
		}
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}
	return &invitation, nil
}

func (r *GormRepository) MarkInvitationAsUsed(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).Model(&Invitation{}).Where("token = ?", token).Update("used", true).Error; err != nil {
		return fmt.Errorf("failed to mark invitation as used: %w", err)
	}
	return nil
}

func (r *GormRepository) DeleteExpiredInvitations(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&Invitation{}).Error; err != nil {
		return fmt.Errorf("failed to delete expired invitations: %w", err)
	}
	return nil
}

// Student/Teacher operations

func (r *GormRepository) CreateStudent(ctx context.Context, student *StudentRecord) error {
	result := r.db.WithContext(ctx).Exec(`
		INSERT INTO students (id, first_name, last_name, email, teacher_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, student.ID, student.FirstName, student.LastName, student.Email, student.TeacherID)

	if result.Error != nil {
		return fmt.Errorf("failed to create student: %w", result.Error)
	}
	return nil
}

func (r *GormRepository) CreateTeacher(ctx context.Context, teacher *TeacherRecord) error {
	result := r.db.WithContext(ctx).Exec(`
		INSERT INTO teachers (id, first_name, last_name, email, department, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, teacher.ID, teacher.FirstName, teacher.LastName, teacher.Email, teacher.Department)

	if result.Error != nil {
		return fmt.Errorf("failed to create teacher: %w", result.Error)
	}
	return nil
}

// ExecuteInTransaction runs the given function within a database transaction
// If the function returns an error, the transaction is automatically rolled back
// Otherwise, the transaction is committed
func (r *GormRepository) ExecuteInTransaction(ctx context.Context, fn func(Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new repository instance with the transaction
		txRepo := &GormRepository{db: tx}

		// Execute the function with the transaction repository
		// If an error is returned, GORM will automatically roll back
		return fn(txRepo)
	})
}
