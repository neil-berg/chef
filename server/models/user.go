package models

import (
	"time"

	"gorm.io/gorm"
)

// User defines the user model
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
}
