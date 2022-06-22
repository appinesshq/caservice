package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/appinesshq/caservice/business/user"
)

var (
	AuthorizationErr  = errors.New("unauthorized")
	AuthenticationErr = errors.New("authentication error")
)

// Config is used to configure UserUseCases.
type Config struct {
	SessionExpires time.Duration
}

// UserUseCases contain application logic for user entities.
type UserUseCases struct {
	Cfg  Config
	Repo UserRepository
}

// New returns an initialized UserUseCases.
func New(cfg Config, r UserRepository) UserUseCases {
	return UserUseCases{Cfg: cfg, Repo: r}
}

// Authenticate returns a Session after succesfully authenticating a user by email and password.
func (uc UserUseCases) Authenticate(ctx context.Context, email, password string) (user.Session, error) {
	u, err := uc.Repo.QueryByEmail(email)
	if err != nil {
		return user.Session{}, AuthenticationErr
	}

	if !u.HasPassword(password) {
		return user.Session{}, AuthenticationErr
	}

	return user.NewSession(u, time.Now().Add(uc.Cfg.SessionExpires)), nil
}

// Create inserts the provided user at the repository.
func (uc UserUseCases) Create(ctx context.Context, u user.User) error {
	s, err := user.GetSession(ctx)
	if err != nil {
		return err
	}

	// Only ADMIN can do this action.
	if !s.UserHasRole(user.RoleAdmin) {
		return AuthorizationErr
	}

	return uc.Repo.Create(u)
}

// Query retrieves all users from the repository,
func (uc UserUseCases) Query(ctx context.Context) ([]user.User, error) {
	s, err := user.GetSession(ctx)
	if err != nil {
		return []user.User{}, err
	}

	// Only ADMIN can do this action.
	if !s.UserHasRole(user.RoleAdmin) {
		return []user.User{}, AuthorizationErr
	}

	return uc.Repo.Query()
}

// QueryByID retrieves a single user from the repository by its id.
func (uc UserUseCases) QueryByID(ctx context.Context, id string) (user.User, error) {
	s, err := user.GetSession(ctx)
	if err != nil {
		return user.User{}, err
	}

	// Only ADMIN or owner can do this action.
	if !s.UserHasRole(user.RoleAdmin) || s.User.ID != id {
		return user.User{}, AuthorizationErr
	}

	return uc.Repo.QueryByID(id)
}

// QueryByEmail retrieves a single user from the repository by its email address.
func (uc UserUseCases) QueryByEmail(ctx context.Context, email string) (user.User, error) {
	s, err := user.GetSession(ctx)
	if err != nil {
		return user.User{}, err
	}

	// Only ADMIN or owner can do this action.
	if !s.UserHasRole(user.RoleAdmin) || s.User.Email != email {
		return user.User{}, AuthorizationErr
	}

	return uc.Repo.QueryByEmail(email)
}

// Update updates the provided user at the repository.
func (uc UserUseCases) Update(ctx context.Context, u user.User) error {
	s, err := user.GetSession(ctx)
	if err != nil {
		return err
	}

	// Only ADMIN or owner can do this action.
	if !s.UserHasRole(user.RoleAdmin) || s.User.ID != u.ID {
		return AuthorizationErr
	}

	return uc.Repo.Update(u)
}

// Delete removes a user from the repository by its id.
func (uc UserUseCases) Delete(ctx context.Context, id string) error {
	s, err := user.GetSession(ctx)
	if err != nil {
		return err
	}

	// Only ADMIN or owner can do this action.
	if !s.UserHasRole(user.RoleAdmin) || s.User.ID != id {
		return AuthorizationErr
	}

	return uc.Repo.Delete(id)
}
