// Package user provides user storage functionality.
package user

import "github.com/appinesshq/caservice/business/user"

// UserStorage is an interface to be implemented by user storage mechanisms.
type UserStorage interface {
	Create(user.User) error
	Query() ([]user.User, error)
	QueryByID(string) (user.User, error)
	QueryByEmail(string) (user.User, error)
	Update(user.User) error
	Delete(string) error
}
