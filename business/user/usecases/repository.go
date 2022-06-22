package usecases

import (
	"errors"
	"fmt"

	"github.com/appinesshq/caservice/business/user"
)

type FieldNotUniqueError string

func (err *FieldNotUniqueError) Error() string {
	return fmt.Sprintf("Non-unique value for field %q", string(*err))
}

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
	Create(user.User) error
	Query() ([]user.User, error)
	QueryByID(string) (user.User, error)
	QueryByEmail(string) (user.User, error)
	Update(user.User) error
	Delete(string) error
}
