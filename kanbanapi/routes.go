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
func (app *App) updateProject(w http.ResponseWriter, r *http.Request) {
	// check if the request body
	// is valid to use to update the project
	var project Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload:  "+err.Error())
		return
	}

	// Read the parameters from the request 'r'
	vars := mux.Vars(r)
	// Get the value of id from the parameter
	id := vars["xata_id"]

	// get the userID from context
	claims := r.Context().Value("claims").(*Claims)
	userID := claims.ID

	// get the project
	// by id provided in the
	// request params
	var storedUserID string
	strQry := `SELECT "user" FROM projects WHERE xata_id=$1`
	err := app.DB.QueryRow(strQry, id).Scan(&storedUserID)
	// check for any errors
	// after querying for
	// the project
	if err != nil {
		// check if we found
		// the project with the
		// requested id
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching project DB.QueryRow():"+id)
		return
	}

	// check if the updating user
	// has rights to update the
	// found project
	if storedUserID != userID {
		respondWithError(w, http.StatusForbidden, "You do not have permission to update this project")
		return
	}

	// query to update the project record
	strQry = `UPDATE projects SET name=$1,
                                repo_url=$2,
                                site_url=$3,
                                description=$4,
                                dependencies=$5,
                                dev_dependencies=$6,
                                status=$7
            WHERE xata_id=$8 and "user"=$9`
	// run the query
	// with DB.Exec
	_, err = app.DB.Exec(strQry, project.Name,
		project.RepoURL,
		project.SiteURL,
		project.Description,
		pq.Array(project.Dependencies),
		pq.Array(project.DevDependencies),
		project.Status,
		id, userID)
	// check if we encountered
	// any errors while updating
	// the updating the project
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Error trying to update the project with id:"+id+" :"+err.Error())
		return
	}
	project.XataID = id
	project.UserID = userID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// Get Project by ID
func (app *App) getProject(w http.ResponseWriter, r *http.Request) {
	// Read the parameters from the request 'r'
	vars := mux.Vars(r)
	// Get the value of id from the parameter
	id := vars["xata_id"]
	claim := r.Context().Value("claims").(*Claims)
	userID := claim.ID

	// handling the returned project
	// by id below
	var project Project
	var dependencies, dev_dependencies []string
	// query the database using
	// the project id passed to the url
	// and the userid found in the context
	strQry := `SELECT xata_id,"user",name,repo_url,site_url,description,dependencies,dev_dependencies,status
            FROM projects WHERE xata_id=$1 AND "user"=$2`
	err := app.DB.QueryRow(strQry, id, userID).Scan(&project.XataID,
		&project.UserID,
		&project.Name,
		&project.RepoURL,
		&project.SiteURL,
		&project.Description,
		pq.Array(&dependencies),
		pq.Array(&dev_dependencies),
		&project.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Project not found with id "+id)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching project: DB.QueryRow() "+id)
		return
	}

	project.Dependencies = dependencies
	project.DevDependencies = dev_dependencies

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// Get all Projects
func (app *App) getProjects(w http.ResponseWriter, r *http.Request) {
	// get the user id from
	// the jwt token via claims
	claims := r.Context().Value("claims").(*Claims)
	userID := claims.ID

	// Query the project
	// records owned by userID
	// note of the function Query
	// compared to what we use prev
	// which is QueryRow(to executy)
	strQry := `SELECT xata_id,"user",name,repo_url,site_url,description,dependencies,dev_dependencies,status
            FROM projects WHERE "user"=$1`
	rows, err := app.DB.Query(strQry, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching project: DB.Query()")
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		var dependencies, dev_dependencies []string
		// scanning each field returned
		// by the query, note the positions
		// are important where you are setting them
		err = rows.Scan(&project.XataID,
			&project.UserID,
			&project.Name,
			&project.RepoURL,
			&project.SiteURL,
			&project.Description,
			pq.Array(dependencies),
			pq.Array(dev_dependencies),
			&project.Status)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning project: rows.scan()")
			return
		}

		// set the array strings
		// to the project variable
		project.Dependencies = dependencies
		project.DevDependencies = dev_dependencies
		// append the project
		// to the projects slice
		projects = append(projects, project)

	}
	if err := rows.Err(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error scanning projects, rows.err()")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// Delete Project by ID
func (app *App) deleteProject(w http.ResponseWriter, r *http.Request) {
	// Read the parameters from the request 'r'
	vars := mux.Vars(r)
	// Get the value of id from the parameter
	id := vars["id"]

	// get the user id from the
	// request context
	claims := r.Context().Value("claims").(*Claims)
	userID := claims.ID

	var storedUserID string
	strQry := `SELECT "user" FROM projects WHERE xata_id=$1`
	err := app.DB.QueryRow(strQry, id).Scan(&storedUserID)
	// check for any errors
	// after querying for
	// the project
	if err != nil {
		// check if we found
		// the project with the
		// requested id
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching project DB.QueryRow():"+id)
		return
	}

	// check if the deleting user
	// has rights to delete the
	// found project
	if storedUserID != userID {
		respondWithError(w, http.StatusForbidden, "You do not have permission to delete this project")
		return
	}

	_, err = app.DB.Exec(`DELETE FROM projects WHERE xata_id=$1 AND "user"=$2`, id, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting project DB.Exec():"+id)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(RouteResponse{Message: "Project successfully deleted. " + id})
}
