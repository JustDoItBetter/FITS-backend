package models

import "time"

type SignRequest struct {
	RequestID   string `parquet:"name=request_id, type=BYTE_ARRAY, convertedtype=UTF8" json:"request_id"`
	StudentUUID string `parquet:"name=student_uuid, type=BYTE_ARRAY, convertedtype=UTF8" json:"student_uuid"`
	WeekNumber  int    `parquet:"name=week_number, type=INT32" json:"week_number"`
	Year        int    `parquet:"name=year, type=INT32" json:"year"`
	Description string `parquet:"name=description, type=BYTE_ARRAY, convertedtype=UTF8" json:"description"`
	CreatedAt   int64  `parquet:"name=created_at, type=INT64, convertedtype=TIMESTAMP_MILLIS" json:"created_at"`
	Status      string `parquet:"name=status, type=BYTE_ARRAY, convertedtype=UTF8" json:"status"`
}

type SignedRequest struct {
	RequestID   string `parquet:"name=request_id, type=BYTE_ARRAY, convertedtype=UTF8" json:"request_id"`
	StudentUUID string `parquet:"name=student_uuid, type=BYTE_ARRAY, convertedtype=UTF8" json:"student_uuid"`
	Signed      bool   `parquet:"name=signed, type=BOOLEAN" json:"signed"`
	Signature   string `parquet:"name=signature, type=BYTE_ARRAY, convertedtype=UTF8" json:"signature"`
	SignedAt    int64  `parquet:"name=signed_at, type=INT64, convertedtype=TIMESTAMP_MILLIS" json:"signed_at"`
	Reason      string `parquet:"name=reason, type=BYTE_ARRAY, convertedtype=UTF8" json:"reason"`
}

type UploadRecord struct {
	UploadID    string `parquet:"name=upload_id, type=BYTE_ARRAY, convertedtype=UTF8" json:"upload_id"`
	StudentUUID string `parquet:"name=student_uuid, type=BYTE_ARRAY, convertedtype=UTF8" json:"student_uuid"`
	FileName    string `parquet:"name=file_name, type=BYTE_ARRAY, convertedtype=UTF8" json:"file_name"`
	FileSize    int64  `parquet:"name=file_size, type=INT64" json:"file_size"`
	UploadedAt  int64  `parquet:"name=uploaded_at, type=INT64, convertedtype=TIMESTAMP_MILLIS" json:"uploaded_at"`
	ContentHash string `parquet:"name=content_hash, type=BYTE_ARRAY, convertedtype=UTF8" json:"content_hash"`
}

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

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
