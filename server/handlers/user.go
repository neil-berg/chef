package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/neil-berg/chef/models"
)

// CreateUser adds a mock user to the DB
func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	mockUser := models.User{
		ID:        2,
		Email:     "test2@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := handler.db.Create(&mockUser)
	if result.Error != nil {
		fmt.Println("error creating mock user")
	}

	bytes, err := json.Marshal(mockUser)
	if err != nil {
		http.Error(w, "bad!", http.StatusBadRequest)
	}
	w.Write(bytes)
}
