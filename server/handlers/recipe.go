package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/neil-berg/chef/models"
)

// AddRecipe adds a recipe to the DB
func (handler *Handler) AddRecipe(w http.ResponseWriter, r *http.Request) {
	ctxKey := models.UserContextKey("user")
	user := r.Context().Value(ctxKey).(models.User)

	recipe := &models.Recipe{}

	err := recipe.ParseBody(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = recipe.Validate()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	recipe.UserID = user.ID
	result := handler.db.Create(recipe)
	if result.Error != nil {
		http.Error(w, "Unable to add recipe", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}

// GetRecipes reads recipes from a user
func (handler *Handler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	ctxKey := models.UserContextKey("user")
	user := r.Context().Value(ctxKey).(models.User)

	recipes := []models.Recipe{}
	result := handler.db.Find(&recipes, "user_id = ?", user.ID)
	if result.Error != nil {
		http.Error(w, "Unable to get recipes", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(recipes)
	if err != nil {
		http.Error(w, "Unable to get recipes", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
