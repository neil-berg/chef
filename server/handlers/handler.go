package handlers

import (
	"log"

	cfg "github.com/neil-berg/chef/config"
	"gorm.io/gorm"
)

// Handler defines the resuable handler shape
type Handler struct {
	logger *log.Logger
	db     *gorm.DB
	config *cfg.Config
}

// CreateHandler returns a pointer to a reusable handler
func CreateHandler(logger *log.Logger, db *gorm.DB, config *cfg.Config) *Handler {
	return &Handler{logger, db, config}
}
