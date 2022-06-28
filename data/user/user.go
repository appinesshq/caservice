// Package user provides user storage functionality.
package user

import (
	"context"

	"github.com/appinesshq/caservice/business/user"
)

// UserStorage is an interface to be implemented by user storages.
type UserStorage interface {
	Create(context.Context, user.User) error
	Query(context.Context) ([]user.User, error)
	QueryByID(context.Context, string) (user.User, error)
	QueryByEmail(context.Context, string) (user.User, error)
	Update(context.Context, user.User) error
	Delete(context.Context, string) error
}

// UserRepository implements the usecases' repository.
//
// All storages in this repo also implement the usecases' repository,
// so the UserRepository struct might seem obsolete. However it has been implemented
// to show the hardcore Clean Architecture filosophy.
//
// In more complex applications the UserRepository layer has a clear added value.
// For example it could use multiple storages, like a cache and persistent database.
type UserRepository struct {
	Storage UserStorage
}

func (r UserRepository) Create(ctx context.Context, u user.User) error {
	return r.Storage.Create(ctx, u)
}

func (r UserRepository) Query(ctx context.Context) ([]user.User, error) {
	return r.Storage.Query(ctx)
}

func (r UserRepository) QueryByID(ctx context.Context, id string) (user.User, error) {
	return r.Storage.QueryByID(ctx, id)
}

func (r UserRepository) QueryByEmail(ctx context.Context, email string) (user.User, error) {
	return r.Storage.QueryByEmail(ctx, email)
}

func (r UserRepository) Update(ctx context.Context, u user.User) error {
	return r.Storage.Update(ctx, u)
}

func (r UserRepository) Delete(ctx context.Context, id string) error {
	return r.Storage.Delete(ctx, id)
}
