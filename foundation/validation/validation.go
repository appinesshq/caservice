// Package validation provides functionality for validation.
package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/crypto/bcrypt"
)

// ValueValidator is an interface for value validation
// functionality. Should return an error when validation fails.
type ValueValidator interface {
	Validate() error
}

// ValidationProvider is a provider for validation services.
// Any ValidationProvider should return ValidationErrors or nil.
type ValidationProvider interface {
	Check(any) error
}

// StandardValidationProvider uses go playground validation.
type StandardValidationProvider struct {
	validate   *validator.Validate
	translator ut.Translator
}

// Check implements the ValidationProvider interface.
func (v *StandardValidationProvider) Check(s any) error {
	err := v.validate.Struct(s)
	if err != nil {
		// Convert implementation specific errors to package errors for consistency.
		errs := err.(validator.ValidationErrors)
		var res ValidationErrors
		for _, e := range errs {
			res = append(res, FieldValidationError{Field: e.StructField(), Err: e.Translate(v.translator)})
		}
		return res.ToValidationError()
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

	// Register custom validation functions.
	// validate.RegisterValidation("uuid", isValidUUID)
	validate.RegisterValidation("notEmptyPassword", isNotEmptyPassword)
	validate.RegisterTranslation("notEmptyPassword", translator, func(ut ut.Translator) error {
		return ut.Add("notEmptyPassword", "{0} can't be an empty string", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notEmptyPassword", fe.Field())
		return t
	})

	return &StandardValidationProvider{validate: validate, translator: translator}
}

// DefaultValidationProvider is a singleton StandardValidationProvider.
var DefaultValidationProvider = NewStandardValidationProvider()

// isValidUUID returns true if the field does NOT contain a hashed empty password.
func isNotEmptyPassword(f validator.FieldLevel) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(f.Field().String()), []byte("")); err == nil {
		return false
	}

	return true
}
