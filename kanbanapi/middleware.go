package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xeipuuv/gojsonschema"
)

func loggingMiddleware(next http.Handler) http.Handler {
	// Log the request in the console
	// and then return the writer and request to the next step
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (app *App) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Invalid Token")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		claims := &Claims{}

		// parse the jwt Token
		// to get information from claims
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return app.TokenKey, nil })
		// check if there is an
		// error trying to parse the token
		if err != nil {
			// do another check if the error was due to
			// invalid token signature
			if err == jwt.ErrSignatureInvalid {
				respondWithError(w, http.StatusUnauthorized, "Invalid Token signature")
				return
			}
			respondWithError(w, http.StatusBadRequest, "Invalid Token \n"+err.Error())
			return
		}

		if !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		// get the context to pass to the response body
		ctx := context.WithValue(r.Context(), "claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validationMiddleware(schema string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// body as key value object equivalent in js
			var body map[string]interface{}
			// read the body from
			// the request into bytes
			bodyBytes, err := io.ReadAll(r.Body)
			// check if there there is no error
			// reading the body bytes
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request payload, check body format")
				return
			}

			// unmarshall the bodyBytes
			// into a key value object
			err = json.Unmarshal(bodyBytes, &body)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request payload, check body format")
				return
			}

			// load the schema
			// we passed to this function
			schemaLoader := gojsonschema.NewStringLoader(schema)
			// load the data from by req1uest body
			// that we just unmarshalled
			documentLoader := gojsonschema.NewGoLoader(body)
			// validate the body using the schema
			result, err := gojsonschema.Validate(schemaLoader, documentLoader)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Error validating payload as JSON")
				return
			}

			// check we the validation
			// results is valid, if not return
			// the validation errors
			if !result.Valid() {
				var validationErrors []string // create a errors slice
				for _, err := range result.Errors() {
					validationErrors = append(validationErrors, err.String())
				}
				respondWithError(w, http.StatusBadRequest, strings.Join(validationErrors, ","))
				return
			}

			// return back the body
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			next.ServeHTTP(w, r)
		})
	}
}
