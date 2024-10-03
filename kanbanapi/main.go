package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
)

func main() {
	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("..Error loading .env file.")
		os.Exit(1)
	}

	connStr := os.Getenv("XATA_PSQL_URL")
	if len(connStr) == 0 {
		log.Fatalf("..XATA_PSQL_URL not set")
	}

	tokenKey := os.Getenv("JWT_TOKEN_SECRET")
	if len(tokenKey) == 0 {
		log.Fatalf("..JWT_TOKEN_SECRET not set")
	}

	// open database connection
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// close the db connection
	// on program close
	defer DB.Close()

	// adding the db connection
	// as an app state
	app := &App{DB: DB, TokenKey: []byte(tokenKey)}

	log.Println("Starting server..")
	// defining router for http request
	router := mux.NewRouter()

	// defining the endpoints
	// adding alice as middleware handler for all request
	// check the `loggingMiddleware` function
	router.Handle("/", alice.New(loggingMiddleware).ThenFunc(rootRoute)).Methods("GET")
	router.Handle("/register", alice.New(loggingMiddleware).ThenFunc(app.register)).Methods("POST")
	router.Handle("/login", alice.New(loggingMiddleware).ThenFunc(app.login)).Methods("POST")
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
