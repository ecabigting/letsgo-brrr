package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateInput(input interface{}) error {
	validate := validator.New()
	return validate.Struct(input)
}

func PasswordComplexity(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	/*
	* Regex for PasswordComplexity
	* Min 8
	* Must have 1 Uppercase
	* Must have 1 Lowercase
	* Must have 1 Symbol
	* Must have 1 Number
	 */

	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Check for at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	// Check for at least one special character
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)
	// Check for minimum length
	isValidLength := len(password) >= 8

	return hasLower && hasUpper && hasDigit && hasSpecial && isValidLength
}
