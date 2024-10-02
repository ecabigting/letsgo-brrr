package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting server..")
	// defining router for http request
	router := mux.NewRouter()

	// defining the endpoints
	router.HandleFunc("/", rootRoute).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/projects", createProject).Methods("POST")
	router.HandleFunc("/projects/{id}", updateProject).Methods("PUT")
	router.HandleFunc("/projects", getProjects).Methods("GET")
	router.HandleFunc("/projects/{id}", getProject).Methods("GET")
	router.HandleFunc("/projects/{id}", deleteProject).Methods("DELETE")
	// setup the http server
	// log any errors that occurs
	port := "6969"
	log.Println("Server running at http://localhost:6969 nice!")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
