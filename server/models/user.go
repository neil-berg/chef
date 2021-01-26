package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// User defines the user model
type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	Email     string         `json:"email" validate:"required,email"`
	Password  string         `json:"password" validate:"required"`
	Token     string         `json:"token"`
}

// ParseBody decodes the JSON-encoded body of the request into a User
func (user *User) ParseBody(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(user)
}

// Validate checks for required fields on the user
func (user *User) Validate() error {
	validate := validator.New()
	return validate.Struct(user)
}
