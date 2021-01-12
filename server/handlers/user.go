package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/neil-berg/chef/models"
)

// GetUsers is a test handler
func GetUsers(w http.ResponseWriter, r *http.Request) {
	mockUser := models.User{
		ID:       "abc",
		Email:    "test@example.com",
		Password: "password",
	}

	bytes, err := json.Marshal(mockUser)
	if err != nil {
		http.Error(w, "bad!", http.StatusBadRequest)
	}
	w.Write(bytes)
}
