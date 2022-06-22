package user

import "github.com/appinesshq/caservice/business/user"

// Storage is an interface to be implemented by user storage mechanisms.
type Storage interface {
	Create(user.User) error
	Query() ([]user.User, error)
	QueryByID(string) (user.User, error)
	QueryByEmail(string) (user.User, error)
	Update(user.User) error
	Delete(string) error
}
