package models

import (
	"time"

	"github.com/google/uuid"
)

// Time string equivalent to Date.now().toISOString in js
func GenerateISOString() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
}

type Base struct {
	ID uuid `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (base *Base) BeforeCreate() error {
	base.ID = uuid.New()

	currentTime := GenerateISOString()

	base.CreatedAt, base.UpdatedAt = currentTime, currentTime

	return nil
}

func (base *Base) AfterUpdate() error {
	base.UpdatedAt = GenerateISOString()
	return nil
}