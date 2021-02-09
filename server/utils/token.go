package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// OneWeekMinutes is the number of minutes in a week
const OneWeekMinutes = 1 * 7 * 24 * 60

// Claims is the payload shape encoded in the JWT
type Claims struct {
	UserID uuid.UUID `json:"userId"`
	jwt.StandardClaims
}

// OneWeekExpiry returns the date of one week from noew
func OneWeekExpiry() time.Time {
	return time.Now().Add(OneWeekMinutes * time.Minute)
}

// CreateJWT generates a new JSON Web Token given a user ID and JWT secret
func CreateJWT(id uuid.UUID, secret string) (string, error) {
	claims := &Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: OneWeekExpiry().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
