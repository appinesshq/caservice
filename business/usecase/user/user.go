// package user contains the user use cases on application logic level.
package user

import (
	"fmt"
	"time"

	entity "github.com/appinesshq/caservice/business/entity/user"
)

// UseCase contains application logic for user.
type UseCase struct {
	Repo Repository
}

// New returns a pointer to an intantiated User UseCase.
func New(r Repository) *UseCase {
	return &UseCase{
		Repo: r,
	}
}

// Register instantiates a new User entity and saves it at the repository.
func (uc *UseCase) Register(name string, email string, password string, roles []entity.Role) (*entity.User, error) {
	user, err := entity.New(name, email, password, roles, time.Now())
	if err != nil {
		return nil, fmt.Errorf("creating user entity: %w", err)
	}

	if err := uc.Repo.Create(user); err != nil {
		return nil, fmt.Errorf("saving user entity: %w", err)
	}

	return user, nil
}

func (uc *UseCase) Authenticate(email string, password string) (string, error) {
	user, err := uc.Repo.GetByEmail(email)
	if err != nil {
		return "", ErrAuthenticationFailed
	}

	if err := user.CheckPassword(password); err != nil {
		return "", ErrAuthenticationFailed
	}

	//TODO: Generate token
	return "12345", nil
}
