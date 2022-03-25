// Package validation provides functionality for validation.
package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// ValueValidator is an interface for value validation
// functionality. Should return an error when validation fails.
type ValueValidator interface {
	Validate() error
}

// ValidationProvider is a provider for validation services.
// Any ValidationProvider should return ValidationErrors or nil.
type ValidationProvider interface {
	Check(interface{}) ValidationErrors
}

// StandardValidationProvider uses go playground validation.
type StandardValidationProvider struct {
	validate   *validator.Validate
	translator ut.Translator
}

// Check implements the ValidationProvider interface.
func (v *StandardValidationProvider) Check(s interface{}) ValidationErrors {
	err := v.validate.Struct(s)
	if err != nil {
		// Convert implementation specific errors to package errors for consistency.
		errs := err.(validator.ValidationErrors)
		var res ValidationErrors
		for _, e := range errs {
			res = append(res, ValidationError{Field: e.StructField(), Value: e.Value(), Err: e.Translate(v.translator)})
		}
		return res
	}
	return nil
}

// NewStandardValidationProvider returns an initialized
// StandardValidationProvider.
func NewStandardValidationProvider() ValidationProvider {
	validate := validator.New()
	// Create a translator for english so the error messages are
	// more human-readable than technical.
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the english error messages for use.
	en_translations.RegisterDefaultTranslations(validate, translator)
	return &StandardValidationProvider{validate: validate, translator: translator}
}

// DefaultValidationProvider is a singleton StandardValidationProvider.
var DefaultValidationProvider = NewStandardValidationProvider()
