package usecases

import (
	"errors"

	"github.com/appinesshq/caservice/business/user"
)

var (
	ErrNotFound    = errors.New("user not found")
	ErrUniqueEmail = errors.New("email already exists")
	ErrUniqueID    = errors.New("id already exists")
)

type UserRepository interface {
	Create(user.User) error
	Query() ([]user.User, error)
	QueryByID(string) (user.User, error)
	QueryByEmail(string) (user.User, error)
	Update(user.User) error
	Delete(string) error
}
