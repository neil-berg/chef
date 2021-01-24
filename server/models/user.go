package models

import (
	"encoding/json"
	"io"
	"time"

	"gorm.io/gorm"
)

// User defines the user model
type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Token     string         `json:"token"`
}

// ParseBody decodes the JSON-encoded body of the request into a User
func (user *User) ParseBody(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(user)
}
