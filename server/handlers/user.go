package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/neil-berg/chef/models"
	"github.com/neil-berg/chef/utils"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser adds a user to the DB
func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	err := user.ParseBody(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = user.Validate()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, "Unable to create ID", http.StatusInternalServerError)
		return
	}
	user.ID = uuid

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hash)

	token, err := utils.CreateJWT(user.ID, handler.config.JWTSecret)
	if err != nil {
		http.Error(w, "Unable to create JWT", http.StatusInternalServerError)
		return
	}
	user.Token = token

	result := handler.db.Create(user)
	if result.Error != nil {
		http.Error(w, "Unable to store user", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: utils.OneWeekExpiry(),
	})
	w.Write([]byte("OK"))
}

// SignInUser checks for a valid email/password combination and if so, refreshes
// the users JSON web token.
func (handler *Handler) SignInUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := user.ParseBody(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = user.Validate()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Save the plain-text password in the request before over-writing this field
	// in the struct during the DB query below
	password := user.Password

	result := handler.db.Where(&models.User{Email: user.Email}).Find(user)
	if result.Error != nil {
		http.Error(w, "Not found", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Not found", http.StatusInternalServerError)
		return
	}

	token, err := utils.CreateJWT(user.ID, handler.config.JWTSecret)
	if err != nil {
		http.Error(w, "Unable to create token", http.StatusInternalServerError)
		return
	}

	result = handler.db.Model(&user).Update("token", token)
	if result.Error != nil {
		http.Error(w, "Unable to update token", http.StatusInternalServerError)
		return
	}
	handler.logger.Printf("Updated token for userId %s", user.ID)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: utils.OneWeekExpiry(),
	})
	w.Write([]byte("OK"))
}

// DeleteMe deletes the authorized user and their recipes from the DB
func (handler *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	ctxKey := models.UserContextKey("user")
	user := r.Context().Value(ctxKey).(models.User)

	result1 := handler.db.Unscoped().Where("user_id = ?", user.ID).Delete(&models.Recipe{})
	result2 := handler.db.Unscoped().Delete(&user)
	if result1.Error != nil || result2.Error != nil {
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK!"))
}

// AuthMe determines if the user is authenticated or not based on a valif token
func (handler *Handler) AuthMe(w http.ResponseWriter, r *http.Request) {
	ctxKey := models.UserContextKey("user")
	user := r.Context().Value(ctxKey).(models.User)

	bytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Unable to marshal user", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
