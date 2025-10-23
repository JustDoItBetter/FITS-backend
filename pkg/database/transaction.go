package database

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager provides transaction management utilities
// This enables explicit transaction boundaries in the service layer
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
//
// Usage:
//
//	err := txMgr.WithTransaction(ctx, func(tx *gorm.DB) error {
//	    // All database operations here are within the transaction
//	    if err := repo.WithDB(tx).Create(ctx, entity); err != nil {
//	        return err // Triggers rollback
//	    }
//	    return nil // Triggers commit
//	})
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return tm.db.WithContext(ctx).Transaction(fn)
}

// WithTransactionValue executes a function within a transaction and returns a value
// This is useful when you need to return a created/updated entity from the transaction
//
// Usage:
//
//	student, err := WithTransactionValue(ctx, txMgr, func(tx *gorm.DB) (*Student, error) {
//	    student, err := repo.WithDB(tx).Create(ctx, req)
//	    if err != nil {
//	        return nil, err
//	    }
//	    return student, nil
//	})
func WithTransactionValue[T any](ctx context.Context, tm *TransactionManager, fn func(tx *gorm.DB) (T, error)) (T, error) {
	var result T
	var returnErr error

	err := tm.WithTransaction(ctx, func(tx *gorm.DB) error {
		var err error
		result, err = fn(tx)
		if err != nil {
			returnErr = err
			return err
		}
		return nil
	})

	if err != nil {
		return result, err
	}
	if returnErr != nil {
		return result, returnErr
	}

	return result, nil
}

// Transactional is an interface for repositories that support transaction contexts
// Repositories should implement this to enable transaction-aware operations
type Transactional interface {
	// WithDB returns a new repository instance using the provided database connection
	// This allows the repository to participate in a transaction
	WithDB(db *gorm.DB) interface{}
}

// GetDB returns the underlying database connection
// This can be used for advanced use cases, but WithTransaction is preferred
func (tm *TransactionManager) GetDB() *gorm.DB {
	return tm.db
}
