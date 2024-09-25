package utilities

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(userID, role, name string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

var ErrInvalidToken = errors.New("invalid token")

// ValidateJWT validates the JWT token and returns the claims if valid
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Check if the token is valid
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}