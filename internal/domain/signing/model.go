package signing

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// SignRequest represents a signature request
// @Description Signature request information
type SignRequest struct {
	RequestID   string `parquet:"name=request_id, type=BYTE_ARRAY, convertedtype=UTF8" json:"request_id" example:"20250930120000-abc123"`
	StudentUUID string `parquet:"name=student_uuid, type=BYTE_ARRAY, convertedtype=UTF8" json:"student_uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	WeekNumber  int    `parquet:"name=week_number, type=INT32" json:"week_number" example:"40"`
	Year        int    `parquet:"name=year, type=INT32" json:"year" example:"2025"`
	Description string `parquet:"name=description, type=BYTE_ARRAY, convertedtype=UTF8" json:"description" example:"Weekly report for week 40"`
	CreatedAt   int64  `parquet:"name=created_at, type=INT64, convertedtype=TIMESTAMP_MILLIS" json:"created_at" example:"1727697600000"`
	Status      string `parquet:"name=status, type=BYTE_ARRAY, convertedtype=UTF8" json:"status" example:"pending"`
}

// SignedRequest represents a signed request
// @Description Signed request information
type SignedRequest struct {
	RequestID   string `parquet:"name=request_id, type=BYTE_ARRAY, convertedtype=UTF8" json:"request_id" example:"20250930120000-abc123"`
	StudentUUID string `parquet:"name=student_uuid, type=BYTE_ARRAY, convertedtype=UTF8" json:"student_uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Signed      bool   `parquet:"name=signed, type=BOOLEAN" json:"signed" example:"true"`
	Signature   string `parquet:"name=signature, type=BYTE_ARRAY, convertedtype=UTF8" json:"signature" example:"base64_signature_data"`
	SignedAt    int64  `parquet:"name=signed_at, type=INT64, convertedtype=TIMESTAMP_MILLIS" json:"signed_at" example:"1727697600000"`
	Reason      string `parquet:"name=reason, type=BYTE_ARRAY, convertedtype=UTF8" json:"reason,omitempty" example:"Approved by teacher"`
}

// UploadRecord represents a file upload record
// @Description File upload record information
type UploadRecord struct {
	UploadID    string `parquet:"name=upload_id, type=BYTE_ARRAY, convertedtype=UTF8" json:"upload_id" example:"upload-123-456"`
	StudentUUID string `parquet:"name=student_uuid, type=BYTE_ARRAY, convertedtype=UTF8" json:"student_uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	FileName    string `parquet:"name=file_name, type=BYTE_ARRAY, convertedtype=UTF8" json:"file_name" example:"report.parquet"`
	FileSize    int64  `parquet:"name=file_size, type=INT64" json:"file_size" example:"1024000"`
	UploadedAt  int64  `parquet:"name=uploaded_at, type=INT64, convertedtype=TIMESTAMP_MILLIS" json:"uploaded_at" example:"1727697600000"`
	ContentHash string `parquet:"name=content_hash, type=BYTE_ARRAY, convertedtype=UTF8" json:"content_hash" example:"sha256:abc123..."`
}

// NewSignRequest creates a new sign request
func NewSignRequest(studentUUID, description string, weekNumber, year int) *SignRequest {
	return &SignRequest{
		RequestID:   generateRequestID(),
		StudentUUID: studentUUID,
		WeekNumber:  weekNumber,
		Year:        year,
		Description: description,
		CreatedAt:   time.Now().UnixMilli(),
		Status:      "pending",
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a cryptographically secure random string of length n
// Uses crypto/rand instead of predictable time-based randomness
func randomString(n int) string {
	// Generate n random bytes (each byte gives us 2 hex characters)
	bytes := make([]byte, (n+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails (extremely unlikely)
		// This should never happen in practice, but ensures we always return a valid ID
		return time.Now().Format("150405.000")[:n]
	}

	// Convert to hex string and truncate to exactly n characters
	hexStr := hex.EncodeToString(bytes)
	if len(hexStr) > n {
		return hexStr[:n]
	}
	return hexStr
}
