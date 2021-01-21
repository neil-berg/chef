package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/neil-berg/chef/models"
	"golang.org/x/crypto/bcrypt"
)

// OneWeekMinutes is the number of minutes in a week, used in JWT expiration.
const OneWeekMinutes = 1 * 7 * 24 * 60

// CreateUser adds a user to the DB
func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := user.ParseBody(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
	}

	// Create a JWT valid for one week from now
	expiresAt := time.Now().Add(OneWeekMinutes * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     user.Email,
		"expiresAt": expiresAt,
	})

	tokenString, err := token.SignedString([]byte(handler.config.JWTSecret))
	if err != nil {
		http.Error(w, "Unable to create JWT", http.StatusInternalServerError)
	}

	user.Password = string(hash)
	user.Token = tokenString

	result := handler.db.Create(user)
	if result.Error != nil {
		http.Error(w, "Unable to store user", http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expiresAt,
	})
	w.Write([]byte("OK"))
}
