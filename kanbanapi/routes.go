package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func rootRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "API Live! Nice!"})
}

// Register User function to handle user registraion API request
func (app *App) register(w http.ResponseWriter, r *http.Request) {
	// json decode the credentials
	// from the request body
	var loginCreds LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&loginCreds)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload.")
		return
	}

	// validate the credentials
	if err := loginCreds.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// hash user password
	// with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginCreds.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error!")
		return
	}
	var returnedXataID string
	err = app.DB.QueryRow(`INSERT INTO "users" (username,password) VALUES($1,$2) RETURNING xata_id`, loginCreds.Username, string(hashedPassword)).Scan(&returnedXataID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: " >> Register Endpoint hit! Nice!"})
}

// Login user
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Login Endpoint hit! Nice!"})
}

// Create Project
func createProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Create Projects Endpoint hit! Nice!"})
}

// Update Project by ID
func updateProject(w http.ResponseWriter, r *http.Request) {
	// Read the parameters from the request 'r'
	vars := mux.Vars(r)
	// Get the value of id from the parameter
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Update Projects Endpoint hit! Nice!", ID: id})
}

// Get Project by ID
func getProject(w http.ResponseWriter, r *http.Request) {
	// Read the parameters from the request 'r'
	vars := mux.Vars(r)
	// Get the value of id from the parameter
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Get Project by ID Endpoint hit! Nice!", ID: id})
}

// Get all Projects
func getProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Get ALL Projects Endpoint hit! Nice!"})
}

// Delete Project by ID
func deleteProject(w http.ResponseWriter, r *http.Request) {
	// Read the parameters from the request 'r'
	vars := mux.Vars(r)
	// Get the value of id from the parameter
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Delete Projects Endpoint hit! Nice!", ID: id})
}
