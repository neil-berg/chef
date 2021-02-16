package handlers

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/neil-berg/chef/models"
	"github.com/neil-berg/chef/utils"
)

// CheckToken checks for a token cookie, verifies the token, and if so, attaches
// the user in context. Otherwise it returns an unauthorized error.
func (handler *Handler) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &utils.Claims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(handler.config.JWTSecret), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieve this user from the DB based on their ID found in the token
		user := models.User{ID: claims.UserID}
		result := handler.db.First(&user)
		if result.Error != nil {
			http.Error(w, "Unauthoirzed", http.StatusUnauthorized)
			return
		}

		ctxKey := models.UserContextKey("user")
		ctx := context.WithValue(r.Context(), ctxKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
