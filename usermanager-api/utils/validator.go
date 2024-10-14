package utils

import "github.com/go-playground/validator/v10"

func ValidateInput(input interface{}) error {
	validate := validator.New()
	return validate.Struct(input)
}
