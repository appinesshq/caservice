// package user contains the user entity with all business logic related to this entity.
package user

import (
	"fmt"
	"time"

	"github.com/appinesshq/caservice/foundation/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User is an entity that contains user data and business logic.
type User struct {
	ID       string    `validate:"required,uuid"`
	Name     string    `validate:"required"`
	Email    string    `validate:"required,email"`
	Password string    `validate:"required"`
	Roles    []string  `validate:"required"`
	Created  time.Time `validate:"required"`
	Modified time.Time `validate:"required"`
}

// New returns an initialized and validated User entity with a generated ID and encrypted password
// or returns an error if password generation or validation fails.
func New(name, email, password string, roles []string, now time.Time) (*User, error) {
	user := User{
		ID:       uuid.NewString(),
		Name:     name,
		Password: password,
		Email:    email,
		Roles:    roles,
		Created:  now,
		Modified: now,
	}

	// Validation takes place before encryption of the password to
	// ensure the password is valid. An encrypted invalid password
	// would always pass. We don't want that.
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Encrypt the password
	if err := user.EncryptPassword(); err != nil {
		return nil, fmt.Errorf("encrypting password: %w", err)
	}

	return &user, nil
}

// Validate implements the Validator interface. Returns an error if validation fails.
func (u *User) Validate() error {
	err := validation.DefaultValidationProvider.Check(u)
	if err != nil {
		return err
	}
	return nil
}

// CheckPassword checks if the provided password matches with the encrypted password
// in the User entity. Returns an error if password is wrong.
func (u *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}

// EncryptPassword will encrypt User.Password.
func (u *User) EncryptPassword() error {
	if u.Password == "" {
		return fmt.Errorf("empty password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = fmt.Sprintf("%s", hash)

	return nil
}
