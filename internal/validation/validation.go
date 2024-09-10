package validation

import (
	"errors"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"net/http"
	"reflect"
)

// Validator struct holds the validator and translation objects
type Validator struct {
	validate             *validator.Validate
	universalTranslation ut.Translator
}

// NewValidator creates a new instance of Validator with custom validation rules
func NewValidator(validate *validator.Validate, universalTranslation ut.Translator) *Validator {
	// Register custom validation for ksuid type
	//nolint:errcheck
	validate.RegisterValidation("ksuid", func(fl validator.FieldLevel) bool {
		switch fl.Field().Kind() {
		case reflect.String:
			// Parse the string and check if it is a valid ksuid
			_, err := ksuid.Parse(fl.Field().String())
			return err == nil
		default:
			// Return false if the field is not a string
			return false
		}
	})
	// Return the Validator instance
	return &Validator{validate: validate, universalTranslation: universalTranslation}
}

// HandleError processes validation errors and returns a slice of response.Error
func (v Validator) HandleError(errs error) []*response.Error {
	// Initialize a slice to store validation errors
	validationErrors := []*response.Error{}
	if errs != nil {
		// Iterate over each validation error
		// Check if the error can be cast to validator.ValidationErrors using errors.As
		var validationErrs validator.ValidationErrors
		if errors.As(errs, &validationErrs) {
			for _, err := range validationErrs {
				// Create a new response.Error for each validation error
				elem := &response.Error{
					Code:   "Unprocessable Entity",
					Status: http.StatusUnprocessableEntity,
					Title:  "Validation Error",
					Detail: fmt.Sprintf("Field '%s' with tag '%s' and value '%v': %s",
						err.Field(), err.Tag(), err.Value(), err.Translate(v.universalTranslation)),
				}
				// Append the error to the slice
				validationErrors = append(validationErrors, elem)
			}
		}
		// Return the slice of validation errors
		return validationErrors
	}
	// Return nil if there are no errors
	return nil
}

// ValidateVars validates a single variable against a tag and returns validation errors if any
func (v Validator) ValidateVars(data interface{}, tag string) []*response.Error {
	// Validate the variable against the provided tag
	errs := v.validate.Var(data, tag)

	// Handle and return any validation errors
	if validationErrors := v.HandleError(errs); validationErrors != nil {
		return validationErrors
	}
	// Return nil if there are no validation errors
	return nil
}

// Validate validates a struct and returns validation errors if any
func (v Validator) Validate(data interface{}) []*response.Error {
	// Validate the struct fields
	errs := v.validate.Struct(data)
	// Handle and return any validation errors
	if validationErrors := v.HandleError(errs); validationErrors != nil {
		return validationErrors
	}
	// Return nil if there are no validation errors
	return nil
}
