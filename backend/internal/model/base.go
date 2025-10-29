package model

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common fields for all persistent entities.
type Base struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;index"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"column:deleted_at;index" swaggerignore:"true"`
}
