package signing

import (
	"context"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
)

// WARNING: EXPERIMENTAL - Signing domain is not yet implemented
// This is a stub implementation that returns 501 Not Implemented errors
// Planned features:
//   - Parquet file upload and parsing for report submission
//   - Digital signature generation and verification with RSA keys
//   - Signature request tracking and management
// Status: Coming in v1.1

// Service handles business logic for signing operations
type Service struct {
	// TODO: Add parquet file handling when ready
	// TODO: Implement RSA signature verification
	// TODO: Add database persistence for signatures
}

// NewService creates a new signing service
func NewService() *Service {
	return &Service{}
}

// HandleUpload processes a parquet file upload
func (s *Service) HandleUpload(ctx context.Context, filename string, data []byte) (*UploadRecord, error) {
	// TODO: Implement parquet file parsing and validation
	return nil, errors.NewAppError(501, "Not Implemented", "upload handling not yet implemented")
}

// GetPendingSignRequests retrieves all pending sign requests as parquet
func (s *Service) GetPendingSignRequests(ctx context.Context) ([]byte, error) {
	// TODO: Implement parquet file generation
	return nil, errors.NewAppError(501, "Not Implemented", "sign requests retrieval not yet implemented")
}

// HandleSignedUploads processes signed request uploads
func (s *Service) HandleSignedUploads(ctx context.Context, filename string, data []byte) error {
	// TODO: Implement signed request processing
	return errors.NewAppError(501, "Not Implemented", "signed uploads handling not yet implemented")
}
