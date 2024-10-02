package main

import (
	"database/sql"
)

type App struct {
	DB *sql.DB
}

type RouteResponse struct {
	// annotate it to json to allow marshalling and unmarshalling
	Message string `json:"message"`
	//for endpoints that doesnt have the id as part of the url request
	ID string `json:"id,omitempty"`
}

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	XataID   string `json:"xata_id"`
	Username string `json:"username"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
