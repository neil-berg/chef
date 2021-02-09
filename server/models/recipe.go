package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Recipe is a basic data shape of a recipe
type Recipe struct {
	ID        uint32         `gorm:"primaryKey; autoIncrement:true" json:"id"`
	UserID    uuid.UUID      `json:"userId"` // Foreign key is UserID
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	Title     string         `json:"title"`
}
