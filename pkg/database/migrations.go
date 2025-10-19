package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for GORM
func (Migration) TableName() string {
	return "schema_migrations"
}

// MigrationFunc is a function that performs a migration
type MigrationFunc func(db *gorm.DB) error

// MigrationDefinition defines a migration with version, name and up function
type MigrationDefinition struct {
	Version string
	Name    string
	Up      MigrationFunc
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []MigrationDefinition
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: getAllMigrations(),
	}
}

// RunMigrations executes all pending migrations
// This function ensures idempotency - it's safe to run multiple times
func (m *Migrator) RunMigrations(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	if err := m.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get already applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Track applied versions
	appliedVersions := make(map[string]bool)
	for _, applied := range appliedMigrations {
		appliedVersions[applied.Version] = true
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		if appliedVersions[migration.Version] {
			fmt.Printf("Migration %s (%s) already applied, skipping\n", migration.Version, migration.Name)
			continue
		}

		fmt.Printf("Applying migration %s (%s)...\n", migration.Version, migration.Name)

		// Run migration in transaction
		err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Execute migration
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("migration failed: %w", err)
			}

			// Record migration
			record := &Migration{
				Version: migration.Version,
				Name:    migration.Name,
			}
			if err := tx.Create(record).Error; err != nil {
				return fmt.Errorf("failed to record migration: %w", err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		fmt.Printf("Migration %s completed successfully\n", migration.Version)
	}

	return nil
}

// ensureMigrationsTable creates the migrations tracking table if it doesn't exist
func (m *Migrator) ensureMigrationsTable() error {
	return m.db.AutoMigrate(&Migration{})
}

// getAppliedMigrations retrieves all applied migrations from the database
func (m *Migrator) getAppliedMigrations() ([]Migration, error) {
	var migrations []Migration
	err := m.db.Order("version ASC").Find(&migrations).Error
	return migrations, err
}

// GetMigrationStatus returns the current migration status
func (m *Migrator) GetMigrationStatus() (applied int, pending int, err error) {
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return 0, 0, err
	}

	applied = len(appliedMigrations)
	pending = len(m.migrations) - applied
	return applied, pending, nil
}

// getAllMigrations returns all migration definitions in order
func getAllMigrations() []MigrationDefinition {
	return []MigrationDefinition{
		{
			Version: "001",
			Name:    "initial_schema",
			Up:      migration001InitialSchema,
		},
		{
			Version: "002",
			Name:    "increase_token_field_size",
			Up:      migration002IncreaseTokenFieldSize,
		},
		{
			Version: "003",
			Name:    "add_soft_delete_columns",
			Up:      migration003AddSoftDeleteColumns,
		},
		{
			Version: "004",
			Name:    "add_teacher_to_invitations_and_make_required",
			Up:      migration004AddTeacherToInvitations,
		},
		// Add future migrations here
	}
}

// migration001InitialSchema creates the initial database schema
func migration001InitialSchema(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return fmt.Errorf("failed to create uuid extension: %w", err)
	}

	// Create students table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			teacher_id UUID,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create students table: %w", err)
	}

	// Create indexes for students
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_students_email ON students(email)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_students_teacher_id ON students(teacher_id)`)

	// Create teachers table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS teachers (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			department VARCHAR(100) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create teachers table: %w", err)
	}

	// Create indexes for teachers
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_teachers_email ON teachers(email)`)

	// Add foreign key constraint for students
	db.Exec(`
		ALTER TABLE students
		DROP CONSTRAINT IF EXISTS fk_students_teacher
	`)
	db.Exec(`
		ALTER TABLE students
		ADD CONSTRAINT fk_students_teacher
		FOREIGN KEY (teacher_id)
		REFERENCES teachers(id)
		ON DELETE SET NULL
	`)

	// Create users table (for authentication)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'teacher', 'student')),
			user_uuid UUID,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			last_login TIMESTAMP WITH TIME ZONE
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create indexes for users
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_user_uuid ON users(user_uuid)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`)

	// Create refresh_tokens table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(500) UNIQUE NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create refresh_tokens table: %w", err)
	}

	// Create indexes for refresh_tokens
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at)`)

	// Create invitations table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS invitations (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			token VARCHAR(500) UNIQUE NOT NULL,
			email VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			role VARCHAR(20) NOT NULL CHECK (role IN ('teacher', 'student')),
			department VARCHAR(100),
			used BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create invitations table: %w", err)
	}

	// Create indexes for invitations
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_invitations_token ON invitations(token)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_invitations_email ON invitations(email)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_invitations_used ON invitations(used)`)

	// Create teacher_keys table (for digital signatures)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS teacher_keys (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			teacher_uuid UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
			public_key TEXT NOT NULL,
			private_key_encrypted TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create teacher_keys table: %w", err)
	}

	// Create indexes for teacher_keys
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_teacher_keys_teacher_uuid ON teacher_keys(teacher_uuid)`)

	// Create reports table (Berichtshefte)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS reports (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			student_uuid UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
			teacher_uuid UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
			week_number INTEGER NOT NULL,
			year INTEGER NOT NULL,
			description TEXT,
			file_path VARCHAR(500) NOT NULL,
			file_hash VARCHAR(64) NOT NULL,
			file_size BIGINT NOT NULL,
			status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'signed', 'rejected')),
			uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			signed_at TIMESTAMP WITH TIME ZONE,
			rejection_reason TEXT
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create reports table: %w", err)
	}

	// Create indexes for reports
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_reports_student_uuid ON reports(student_uuid)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_reports_teacher_uuid ON reports(teacher_uuid)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_reports_week_year ON reports(week_number, year)`)

	// Create signatures table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS signatures (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
			teacher_uuid UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
			signature TEXT NOT NULL,
			public_key_id UUID NOT NULL REFERENCES teacher_keys(id) ON DELETE CASCADE,
			signed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create signatures table: %w", err)
	}

	// Create indexes for signatures
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_signatures_report_id ON signatures(report_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_signatures_teacher_uuid ON signatures(teacher_uuid)`)

	return nil
}

// migration002IncreaseTokenFieldSize increases token field size for invitations and refresh_tokens
// JWT tokens can be longer than 255 characters, especially with custom claims
func migration002IncreaseTokenFieldSize(db *gorm.DB) error {
	// Increase invitations.token from VARCHAR(255) to VARCHAR(500)
	if err := db.Exec(`
		ALTER TABLE invitations
		ALTER COLUMN token TYPE VARCHAR(500)
	`).Error; err != nil {
		return fmt.Errorf("failed to alter invitations.token: %w", err)
	}

	// Ensure refresh_tokens.token is also VARCHAR(500) (should already be, but just in case)
	if err := db.Exec(`
		ALTER TABLE refresh_tokens
		ALTER COLUMN token TYPE VARCHAR(500)
	`).Error; err != nil {
		return fmt.Errorf("failed to alter refresh_tokens.token: %w", err)
	}

	return nil
}

// migration003AddSoftDeleteColumns adds deleted_at columns for GORM soft delete support
// GORM models use gorm.DeletedAt which requires this column to exist
func migration003AddSoftDeleteColumns(db *gorm.DB) error {
	// Add deleted_at column to students table
	if err := db.Exec(`
		ALTER TABLE students
		ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE
	`).Error; err != nil {
		return fmt.Errorf("failed to add deleted_at to students: %w", err)
	}

	// Add deleted_at column to teachers table
	if err := db.Exec(`
		ALTER TABLE teachers
		ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE
	`).Error; err != nil {
		return fmt.Errorf("failed to add deleted_at to teachers: %w", err)
	}

	// Create indexes for soft delete queries (improves performance)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_students_deleted_at ON students(deleted_at)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_teachers_deleted_at ON teachers(deleted_at)`)

	return nil
}

// migration004AddTeacherToInvitations adds teacher_uuid to invitations table
// This makes the student-teacher relationship mandatory during invitation
func migration004AddTeacherToInvitations(db *gorm.DB) error {
	// Add teacher_uuid column to invitations table
	// For students, this is required. For teachers, it remains NULL
	if err := db.Exec(`
		ALTER TABLE invitations
		ADD COLUMN IF NOT EXISTS teacher_uuid UUID
	`).Error; err != nil {
		return fmt.Errorf("failed to add teacher_uuid to invitations: %w", err)
	}

	// Add foreign key constraint
	db.Exec(`
		ALTER TABLE invitations
		DROP CONSTRAINT IF EXISTS fk_invitations_teacher
	`)

	if err := db.Exec(`
		ALTER TABLE invitations
		ADD CONSTRAINT fk_invitations_teacher
		FOREIGN KEY (teacher_uuid)
		REFERENCES teachers(id)
		ON DELETE SET NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to add foreign key constraint: %w", err)
	}

	// Create index for teacher_uuid lookups
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_invitations_teacher_uuid ON invitations(teacher_uuid)`)

	return nil
}
