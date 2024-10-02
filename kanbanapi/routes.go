package main

import (
	"encoding/json"
	"net/http"
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

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Login Endpoint hit! Nice!"})

}

func createProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Create Projects Endpoint hit! Nice!"})

}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Update Projects Endpoint hit! Nice!"})

}

func getProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Get Project by ID Endpoint hit! Nice!"})

}

func getProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Get ALL Projects Endpoint hit! Nice!"})

}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Delete Projects Endpoint hit! Nice!"})

}
