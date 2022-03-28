package user

import (
	"errors"

	entity "github.com/appinesshq/caservice/business/entity/user"
)

var (
	// ErrAuthenticationFailed is returned upon an authentication failure.
	ErrAuthenticationFailed = errors.New("authenticated failed")
	// ErrNotFound is returned when a user is not found.
	ErrNotFound = errors.New("user not found")
	// ErrExists is returned when a user already exists.
	ErrExists = errors.New("user already exists")
)

type Reader interface {
	GetByID(string) (*entity.User, error)
	GetByEmail(string) (*entity.User, error)
	Query() ([]entity.User, error)
}
type Writer interface {
	Create(*entity.User) error
	Update(*entity.User) error
	Delete(string) error
}
type Repository interface {
	Reader
	Writer
}

type Service interface {
	Authenticate(string, string) (string, error)
	Register(string, string, string, []string) (*entity.User, error)
}
