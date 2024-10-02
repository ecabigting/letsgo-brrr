package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func rootRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "API Live! Nice!"})
}

// Register User
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Register Endpoint hit! Nice!"})

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
