package handlers

import (
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
	}

	err = user.Validate()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, "Unable to create ID", http.StatusInternalServerError)
	}
	user.ID = uuid.String()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
	}
	user.Password = string(hash)

	token, err := utils.CreateJWT(user.ID, handler.config.JWTSecret)
	if err != nil {
		http.Error(w, "Unable to create JWT", http.StatusInternalServerError)
	}
	user.Token = token

	result := handler.db.Create(user)
	if result.Error != nil {
		http.Error(w, "Unable to store user", http.StatusInternalServerError)
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
	}

	err = user.Validate()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	// Save the plain-text password in the request before over-writing this field
	// in the struct during the DB query below
	password := user.Password

	result := handler.db.Where(&models.User{Email: user.Email}).Find(user)
	if result.Error != nil {
		http.Error(w, "Not found", http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Not found", http.StatusInternalServerError)
	}

	token, err := utils.CreateJWT(user.ID, handler.config.JWTSecret)
	if err != nil {
		http.Error(w, "Unable to create token", http.StatusInternalServerError)
	}

	result = handler.db.Model(&user).Update("token", token)
	if result.Error != nil {
		http.Error(w, "Unable to update token", http.StatusInternalServerError)
	}
	handler.logger.Printf("Updated token for userId %s", user.ID)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: utils.OneWeekExpiry(),
	})
	w.Write([]byte("OK"))
}
