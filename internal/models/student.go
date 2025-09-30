package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	UUID      string    `toml:"uuid" json:"uuid" validate:"required,uuid"`
	FirstName string    `toml:"first_name" json:"first_name" validate:"required,min=1,max=100"`
	LastName  string    `toml:"last_name" json:"last_name" validate:"required,min=1,max=100"`
	Email     string    `toml:"email" json:"email" validate:"required,email"`
	TeacherID string    `toml:"teacher_id" json:"teacher_id" validate:"required,uuid"`
	CreatedAt time.Time `toml:"created_at" json:"created_at"`
	UpdatedAt time.Time `toml:"updated_at" json:"updated_at"`
}

type StudentConfig struct {
	UUID      string `toml:"uuid"`
	FirstName string `toml:"first_name"`
	LastName  string `toml:"last_name"`
	Email     string `toml:"email"`
	TeacherID string `toml:"teacher_id"`
}

func NewStudent(config StudentConfig) *Student {
	now := time.Now()

	studentUUID := config.UUID
	if studentUUID == "" {
		studentUUID = uuid.New().String()
	}

	return &Student{
		UUID:      studentUUID,
		FirstName: config.FirstName,
		LastName:  config.LastName,
		Email:     config.Email,
		TeacherID: config.TeacherID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *Student) Update(config StudentConfig) {
	if config.FirstName != "" {
		s.FirstName = config.FirstName
	}
	if config.LastName != "" {
		s.LastName = config.LastName
	}
	if config.Email != "" {
		s.Email = config.Email
	}
	if config.TeacherID != "" {
		s.TeacherID = config.TeacherID
	}
	s.UpdatedAt = time.Now()
}
