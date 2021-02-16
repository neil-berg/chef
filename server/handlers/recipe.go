package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

// GetRecipe reads one recipe by its ID
func (handler *Handler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeIDStr := vars["recipeID"]
	recipeID, err := strconv.ParseUint(recipeIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Unable to get recipe", http.StatusInternalServerError)
		return
	}

	recipe := models.Recipe{ID: uint32(recipeID)}
	result := handler.db.First(&recipe)
	if result.Error != nil {
		http.Error(w, "Unable to get recipe", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(recipe)
	if err != nil {
		http.Error(w, "Unable to get recipe", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
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

// UpdateRecipe updates one recipe based on the sent updated fields
func (handler *Handler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeIDStr := vars["recipeID"]
	recipeID, err := strconv.ParseUint(recipeIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Incoming recipe in the request body
	incomingRecipe := &models.Recipe{}
	err = incomingRecipe.ParseBody(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = incomingRecipe.Validate()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Read the existing recipe
	recipe := models.Recipe{ID: uint32(recipeID)}
	result := handler.db.Find(&recipe)
	if result.Error != nil {
		http.Error(w, "Unable to update recipe", http.StatusInternalServerError)
		return
	}

	// Apply updates and save
	recipe.Title = incomingRecipe.Title
	recipe.Ingredients = incomingRecipe.Ingredients
	recipe.Instructions = incomingRecipe.Instructions

	result = handler.db.Save(&recipe)
	if result.Error != nil {
		http.Error(w, "Unable to update recipe", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(recipe)
	if err != nil {
		http.Error(w, "Unable to update recipe", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// DeleteRecipe deletes one recipe by its ID
func (handler *Handler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeIDStr := vars["recipeID"]
	recipeID, err := strconv.ParseUint(recipeIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Unable to delete recipe", http.StatusInternalServerError)
		return
	}

	recipe := models.Recipe{ID: uint32(recipeID)}
	result := handler.db.Unscoped().Delete(&recipe)
	if result.Error != nil {
		http.Error(w, "Unable to delete recipe", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}
