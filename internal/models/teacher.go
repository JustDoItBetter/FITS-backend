package models

import (
	"time"

	"github.com/google/uuid"
)

type Teacher struct {
	UUID       string    `toml:"uuid" json:"uuid" validate:"required,uuid"`
	FirstName  string    `toml:"first_name" json:"first_name" validate:"required,min=1,max=100"`
	LastName   string    `toml:"last_name" json:"last_name" validate:"required,min=1,max=100"`
	Email      string    `toml:"email" json:"email" validate:"required,email"`
	Department string    `toml:"department" json:"department" validate:"required,min=1,max=100"`
	CreatedAt  time.Time `toml:"created_at" json:"created_at"`
	UpdatedAt  time.Time `toml:"updated_at" json:"updated_at"`
}

type TeacherConfig struct {
	UUID       string `toml:"uuid"`
	FirstName  string `toml:"first_name"`
	LastName   string `toml:"last_name"`
	Email      string `toml:"email"`
	Department string `toml:"department"`
}

func NewTeacher(config TeacherConfig) *Teacher {
	now := time.Now()

	teacherUUID := config.UUID
	if teacherUUID == "" {
		teacherUUID = uuid.New().String()
	}

	return &Teacher{
		UUID:       teacherUUID,
		FirstName:  config.FirstName,
		LastName:   config.LastName,
		Email:      config.Email,
		Department: config.Department,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (t *Teacher) Update(config TeacherConfig) {
	if config.FirstName != "" {
		t.FirstName = config.FirstName
	}
	if config.LastName != "" {
		t.LastName = config.LastName
	}
	if config.Email != "" {
		t.Email = config.Email
	}
	if config.Department != "" {
		t.Department = config.Department
	}
	t.UpdatedAt = time.Now()
}
