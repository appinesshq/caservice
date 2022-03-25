package validation

import (
	"errors"
	"fmt"
)

// Validation error is a validation error for a single field.
type ValidationError struct {
	Field string
	Err   string
	Value interface{}
}

// Error implements the Error interface.
func (ve ValidationError) Error() string {
	return fmt.Sprintf("field %q with value %q is invalid: %s", ve.Field, ve.Value, ve.Err)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// ToMap returns a map with the error string for each field.
func (ve ValidationErrors) ToMap() map[string]string {
	m := make(map[string]string, len(ve))
	for _, err := range ve {
		m[err.Field] = err.Error()
	}
	return m
}

// IsValidationErrors returns true if the provided error is
// a ValidationsErrors
func IsValidationErrors(err error) bool {
	var ve ValidationErrors
	return errors.As(err, &ve)
}

// IsValidationError returns true if the provided error is
// a ValidationsErrors or ValidationError.
func IsValidationError(err error) bool {
	var ve ValidationError
	return errors.As(err, &ve)
}

// IsAnyValidationError returns true if the provided error is
// a ValidationsErrors or ValidationError.
func IsAnyValidationError(err error) bool {
	return IsValidationErrors(err) || IsValidationError(err)
}

// GetValidationErrors returns the error as ValidationErrors
// if having that type.
func GetValidationErrors(err error) ValidationErrors {
	var ve ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}
	return ve
}

// Error implements the Error interface.
func (ve ValidationErrors) Error() string {
	m := ve.ToMap()
	return fmt.Sprintf("%+v", m)
}
