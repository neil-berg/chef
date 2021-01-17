package database

import (
	"fmt"

	"github.com/neil-berg/chef/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connects to the database
func Connect(host, port, dbname, user, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}

// Migrate automatically migrates schemas on the DB
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	return err
}
