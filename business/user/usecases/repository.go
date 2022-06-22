package usecases

import "github.com/appinesshq/caservice/business/user"

type UserRepository interface {
	Create(user.User) error
	Query() ([]user.User, error)
	QueryByID(string) (user.User, error)
	QueryByEmail(string) (user.User, error)
	Update(user.User) error
	Delete(string) error
}
