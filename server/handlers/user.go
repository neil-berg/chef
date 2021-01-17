package handlers

import (
	"fmt"
	"net/http"

	"github.com/neil-berg/chef/models"
)

// CreateUser adds a user to the DB
func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := user.ParseBody(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	result := handler.db.Create(user)
	if result.Error != nil {
		fmt.Println("error creating mock user")
	}

	// bytes, err := json.Marshal(mockUser)
	// if err != nil {
	// 	http.Error(w, "bad!", http.StatusBadRequest)
	// }
	w.Write([]byte("OK"))
}
