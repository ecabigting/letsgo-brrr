package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	log.Println("Starting server..")
	// defining router for http request
	router := mux.NewRouter()

	// defining the endpoints
	// adding alice as middleware handler for all request
	// check the `loggingMiddleware` function
	router.Handle("/", alice.New(loggingMiddleware).ThenFunc(rootRoute)).Methods("GET")
	router.Handle("/register", alice.New(loggingMiddleware).ThenFunc(register)).Methods("POST")
	router.Handle("/login", alice.New(loggingMiddleware).ThenFunc(login)).Methods("POST")
	router.Handle("/projects", alice.New(loggingMiddleware).ThenFunc(createProject)).Methods("POST")
	router.Handle("/projects/{id}", alice.New(loggingMiddleware).ThenFunc(updateProject)).Methods("PUT")
	router.Handle("/projects", alice.New(loggingMiddleware).ThenFunc(getProjects)).Methods("GET")
	router.Handle("/projects/{id}", alice.New(loggingMiddleware).ThenFunc(getProject)).Methods("GET")
	router.Handle("/projects/{id}", alice.New(loggingMiddleware).ThenFunc(deleteProject)).Methods("DELETE")
	// setup the http server
	// log any errors that occurs
	port := "6969"
	log.Println("Server running at http://localhost:6969 nice!")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
