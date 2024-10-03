package main

import (
	"database/sql"

	"github.com/golang-jwt/jwt/v5"
)

type App struct {
	DB       *sql.DB
	TokenKey []byte
}

type RouteResponse struct {
	// annotate it to json to allow marshalling and unmarshalling
	Message string `json:"message"`
	// for endpoints that doesnt have the id as part of the url request
	ID string `json:"id,omitempty"`
}

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	XataID   string `json:"xata_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Claims struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	// adding the jwt Claims
	jwt.RegisteredClaims
}
