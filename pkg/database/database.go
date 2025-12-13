// Package database provides database connection and management utilities for PostgreSQL.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the gorm database instance
type DB struct {
	*gorm.DB
}

// New creates a new database connection with automatic initialization
// It performs the following steps:
// 1. Ensures the database exists (creates if needed)
// 2. Establishes connection with connection pooling
// 3. Runs all pending migrations
// 4. Verifies the connection health
//
// This function is idempotent - safe to call multiple times
func New(cfg *config.DatabaseConfig) (*DB, error) {
	log.Println("Initializing database connection...")

	// Ensure database exists
	if err := ensureDatabaseExists(cfg); err != nil {
		return nil, fmt.Errorf("failed to ensure database exists: %w", err)
	}

	// Build DSN for application database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database: %w", err)
	}

	// Configure connection pool
	if cfg.MaxConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxConns)
	}
	if cfg.MinConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MinConns)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	// Wrap in DB struct
	wrappedDB := &DB{DB: db}

	// Run migrations automatically
	if err := wrappedDB.RunMigrations(ctx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database initialization complete")

	return wrappedDB, nil
}

// ensureDatabaseExists checks if the database exists and creates it if needed
// It connects to the 'postgres' database to perform administrative tasks
func ensureDatabaseExists(cfg *config.DatabaseConfig) error {
	// Connect to 'postgres' database for administrative operations
	adminDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", adminDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer func() { _ = db.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if database exists
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	err = db.QueryRowContext(ctx, query, cfg.Database).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if exists {
		log.Printf("Database '%s' already exists", cfg.Database)
		return nil
	}

	// Create database
	log.Printf("Database '%s' does not exist, creating...", cfg.Database)
	createQuery := fmt.Sprintf(`CREATE DATABASE %s`, cfg.Database)
	_, err = db.ExecContext(ctx, createQuery)
	if err != nil {
		return fmt.Errorf("failed to create database '%s': %w", cfg.Database, err)
	}

	log.Printf("Database '%s' created successfully", cfg.Database)
	return nil
}

// RunMigrations executes all pending database migrations
// This is automatically called by New() but can be called manually if needed
func (db *DB) RunMigrations(ctx context.Context) error {
	log.Println("Running database migrations...")

	migrator := NewMigrator(db.DB)
	if err := migrator.RunMigrations(ctx); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Get migration status
	applied, pending, err := migrator.GetMigrationStatus()
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	log.Printf("Migrations complete: %d applied, %d pending", applied, pending)
	return nil
}

// VerifySchema checks if all expected tables exist in the database
// Returns an error if any critical tables are missing
func (db *DB) VerifySchema(ctx context.Context) error {
	requiredTables := []string{
		"users",
		"refresh_tokens",
		"invitations",
		"students",
		"teachers",
		"reports",
		"signatures",
		"teacher_keys",
		"schema_migrations",
	}

	var missingTables []string
	for _, table := range requiredTables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)
		`
		if err := db.WithContext(ctx).Raw(query, table).Scan(&exists).Error; err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}

		if !exists {
			missingTables = append(missingTables, table)
		}
	}

	if len(missingTables) > 0 {
		return fmt.Errorf("missing required tables: %v", missingTables)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database: %w", err)
	}
	return sqlDB.Close()
}

// Health checks if the database connection is alive
func (db *DB) Health(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// WithTx executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise it is committed
func (db *DB) WithTx(ctx context.Context, fn func(*gorm.DB) error) error {
	return db.DB.WithContext(ctx).Transaction(fn)
}
