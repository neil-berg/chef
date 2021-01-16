package handlers

import (
	"log"

	"gorm.io/gorm"
)

// Handler defines the resuable handler shape
type Handler struct {
	logger *log.Logger
	db     *gorm.DB
}

// CreateHandler returns a pointer to a reusable handler
func CreateHandler(logger *log.Logger, db *gorm.DB) *Handler {
	return &Handler{logger, db}
}
