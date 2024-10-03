package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code) // return bad request
	json.NewEncoder(w).Encode(ErrorResponse{Message: msg})
}

func (lc *LoginCredentials) Validate() error {
	if lc.Username == "" {
		return errors.New("Username is required")
	}

	if lc.Password == "" {
		return errors.New("Password is required.")
	}

	return nil
}

// Generate Token function
// takes the user name and id
// returns a jwt token and/or error
func (app *App) generateToken(uname, id string) (string, error) {
	// get an expiration time in the future using the time function
	expTime := time.Now().Add(3 * time.Hour)

	// create the claims for the jwt
	claim := &Claims{
		Username: uname,
		ID:       id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(app.TokenKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
