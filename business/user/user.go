// Package user provides user (related) entities and business logic.
package user

import (
	"fmt"
	"time"

	"github.com/appinesshq/caservice/foundation/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

type User struct {
	ID           string    `validate:"required,uuid"`
	Name         string    `validate:"required"`
	Email        string    `validate:"required,email"`
	PasswordHash []byte    `validate:"required,notEmptyPassword"`
	Roles        []string  `validate:"required,min=1"`
	DateCreated  time.Time `validate:"required"`
	DateUpdated  time.Time `validate:"required"`
}

func New(id, name, email, password string, roles []string, now time.Time) (User, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("encrypting password: %w", err)
	}

	u := User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: b,
		Roles:        roles,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := u.Validate(); err != nil {
		return User{}, fmt.Errorf("validation error: %w", err)
	}
	return u, nil
}

func NewWithID(name, email, password string, roles []string, now time.Time) (User, error) {
	return New(uuid.New().String(), name, email, password, roles, now)
}

func (u User) Validate() error {
	if err := validation.DefaultValidationProvider.Check(u); err != nil {
		return err
	}

	return nil
}

func (u User) HasPassword(s string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(s)); err != nil {
		return false
	}
	return true
}
