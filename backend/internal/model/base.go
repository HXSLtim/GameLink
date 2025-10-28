package model

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common fields for all persistent entities.
type Base struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}
