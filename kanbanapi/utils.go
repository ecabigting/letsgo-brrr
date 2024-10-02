package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code) // return bad request
	json.NewEncoder(w).Encode(ErrorResponse{Message: msg})
}
