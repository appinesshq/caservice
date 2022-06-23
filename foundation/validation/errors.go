package validation

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Err    string
	Fields map[string]string
}

func (ve ValidationError) Error() string {
	return ve.Err
}

func IsValidationError(e error) bool {
	var ve ValidationError
	return errors.As(e, &ve)
}

// FieldValidationError is a validation error for a single field.
type FieldValidationError struct {
	Field string
	Err   string
}

// Error implements the Error interface.
func (ve FieldValidationError) Error() string {
	return ve.Err
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []FieldValidationError

// ToMap returns a map with the error string for each field.
func (ve ValidationErrors) ToMap() map[string]string {
	m := make(map[string]string, len(ve))
	for _, err := range ve {
		m[err.Field] = err.Error()
	}
	return m
}

//
func (ve ValidationErrors) ToValidationError() ValidationError {
	return ValidationError{Err: "data validation error", Fields: ve.ToMap()}
}

// Error implements the Error interface.
func (ve ValidationErrors) Error() string {
	m := ve.ToMap()
	return fmt.Sprintf("%+v", m)
}

// IsValidationErrors returns true if the provided error is
// a ValidationsErrors
func IsValidationErrors(err error) bool {
	var ve ValidationErrors
	return errors.As(err, &ve)
}

// IsFieldValidationError returns true if the provided error is
// a ValidationsErrors or ValidationError.
func IsFieldValidationError(err error) bool {
	var ve FieldValidationError
	return errors.As(err, &ve)
}

// IsAnyValidationError returns true if the provided error is
// a ValidationsErrors or ValidationError.
func IsAnyValidationError(err error) bool {
	return IsValidationErrors(err) || IsFieldValidationError(err) || IsValidationError(err)
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
