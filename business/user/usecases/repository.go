package usecases

import (
	"context"
	"errors"

	"github.com/appinesshq/caservice/business/user"
)

var (
	ErrNotFound    = errors.New("user not found")
	ErrUniqueEmail = errors.New("email already exists")
	ErrUniqueID    = errors.New("id already exists")
)

// UserRepository is an interface which is to be implemented by the layer
// between user usecases and storages.
type UserRepository interface {
	Create(context.Context, user.User) error
	Query(context.Context) ([]user.User, error)
	QueryByID(context.Context, string) (user.User, error)
	QueryByEmail(context.Context, string) (user.User, error)
	Update(context.Context, user.User) error
	Delete(context.Context, string) error
}
