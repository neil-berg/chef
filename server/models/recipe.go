package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Recipe is a basic data shape of a recipe
type Recipe struct {
	ID           uint32         `gorm:"primaryKey; autoIncrement:true" json:"id"`
	UserID       uuid.UUID      `json:"userId"` // Foreign key is UserID
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt"`
	Title        string         `json:"title" validate:"required"`
	Ingredients  pq.StringArray `gorm:"type:text[]" json:"ingredients" validate:"required"`
	Instructions pq.StringArray `gorm:"type:text[]" json:"instructions" validate:"required"`
}

// ParseBody decodes the JSON-encoded body of the request into a Recipe
func (recipe *Recipe) ParseBody(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(recipe)
}

// Validate checks for required fields on the recipe
func (recipe *Recipe) Validate() error {
	validate := validator.New()
	return validate.Struct(recipe)
}
