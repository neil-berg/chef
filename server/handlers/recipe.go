package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neil-berg/chef/models"
)

// TODO: AddRecipe, DeleteRecipe

// GetRecipes reads recipes from a user
func (handler *Handler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	ctxKey := models.UserContextKey("user")
	user := r.Context().Value(ctxKey).(models.User)

	recipes := []models.Recipe{}
	result := handler.db.Find(&recipes, "user_id = ?", user.ID)
	if result.Error != nil {
		http.Error(w, "Unable to get recipes", http.StatusInternalServerError)
	}

	bytes, err := json.Marshal(recipes)
	if err != nil {
		http.Error(w, "Unable to get recipes", http.StatusInternalServerError)
	}

	fmt.Println(result)

	w.Write(bytes)
}
