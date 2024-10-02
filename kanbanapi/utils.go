package main

import (
	"encoding/json"
	"errors"
	"net/http"
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
