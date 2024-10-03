package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func rootRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "API Live! Nice!"})
}

// Register User function to handle user registraion API request
// and set it as function receiver for the App struct
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

	// insert the records to user table
	// and set the returning XataID to
	// returnedXataID, then check
	// if we encounter an error
	var returnedXataID string
	err = app.DB.QueryRow(`INSERT INTO "users" (username,password) VALUES($1,$2) RETURNING xata_id`, loginCreds.Username, string(hashedPassword)).Scan(&returnedXataID)
	if err != nil {
		respondWithError(w, http.StatusBadGateway, err.Error()+"\nError Creating User")
		return
	}

	// generate a JWT token
	tokenString, err := app.generateToken(loginCreds.Username, returnedXataID)
	// check if there was an
	// error generating the jwt
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating access token"+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{XataID: returnedXataID, Username: loginCreds.Username, Token: tokenString})
}

// Login user
func (app *App) login(w http.ResponseWriter, r *http.Request) {
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

	// check if user credentials
	// exist in the database
	var storedCreds LoginCredentials
	var returnedXataID string
	err = app.DB.QueryRow(`SELECT xata_id,username,password FROM "users" WHERE username=$1`, loginCreds.Username).Scan(&returnedXataID, &storedCreds.Username, &storedCreds.Password)
	// check if we encounter
	// any kind of errors
	if err != nil {
		// check if we found
		// any rows for the username
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error()+"\nError: Invalid request payload!")
		return
	}

	// check if password match
	// using bycrpt and converting
	// the stored password and
	// request password as byte slices
	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(loginCreds.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// generate a JWT token
	tokenString, err := app.generateToken(loginCreds.Username, returnedXataID)
	// check if there was an
	// error generating the jwt
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating access token"+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{XataID: returnedXataID, Username: loginCreds.Username, Token: tokenString})
}

// Create Project
func (app *App) createProject(w http.ResponseWriter, r *http.Request) {
	var project Project
	// decode the project from
	// the request body using json.NewDecoder
	// and then pass it to the project variable
	// above via &project pointer
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	// get the userID from
	// claims via context
	// r.Context() : reads the context from the request
	// .Value("claims") : get the value with key "claims", see jwtMiddleware function
	// .(*Claims) parse it into the Claims object
	claims := r.Context().Value("claims").(*Claims)
	userID := claims.ID
	log.Println(claims)

	// the insert query
	// to add new Project into the
	// project table in the DB
	inserQuery := `
        INSERT INTO projects ("user",name,repo_url,site_url,description,dependencies,dev_dependencies,status)
                      VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING xata_id
  `
	var xataID string
	log.Println(userID)
	err = app.DB.QueryRow(inserQuery,
		userID,
		project.Name,
		project.RepoURL,
		project.SiteURL,
		project.Description,
		pq.Array(project.Dependencies),
		pq.Array(project.DevDependencies),
		project.Status).Scan(&xataID)
	// check if there was an error trying to
	// add project into database
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Creating new Project: "+err.Error())
		return
	}

	project.XataID = xataID
	project.UserID = userID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
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
func (app *App) getProjects(w http.ResponseWriter, r *http.Request) {
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
