package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator holds the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate validates a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// GetValidationErrors extracts validation error messages
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch tag {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = field + " must be a valid email"
			case "min":
				errors[field] = field + " is too short"
			case "max":
				errors[field] = field + " is too long"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}

	return errors
}
