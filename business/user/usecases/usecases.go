package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/appinesshq/caservice/business/user"
	"go.uber.org/zap"
)

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrAuthenticationFailed = errors.New("authentication failed")
)

// Config is used to configure UserUseCases.
type Config struct {
}

// UserUseCases contain application logic for user entities.
type UserUseCases struct {
	Log             *zap.SugaredLogger
	Repo            UserRepository
	SessionDuration time.Duration
}

// New returns an initialized UserUseCases.
// Requires a logger, user repository and user session duration as input.
func New(log *zap.SugaredLogger, r UserRepository, s time.Duration) UserUseCases {
	return UserUseCases{Log: log, Repo: r, SessionDuration: s}
}

// Authenticate returns a Session after succesfully authenticating a user by email and password.
func (uc UserUseCases) Authenticate(ctx context.Context, email, password string, now time.Time) (user.Session, error) {
	u, err := uc.Repo.QueryByEmail(email)
	if err != nil {
		return user.Session{}, ErrAuthenticationFailed
	}

	if !u.HasPassword(password) {
		return user.Session{}, ErrAuthenticationFailed
	}

	return user.NewSession(u, now.Add(uc.SessionDuration)), nil
}

// Create inserts the provided user at the repository.
func (uc UserUseCases) Create(ctx context.Context, u user.User) error {
	s, err := user.GetSession(ctx)
	if err != nil {
		return err
	}

	// Only ADMIN can do this action.
	if !s.UserHasRole(user.RoleAdmin) {
		return ErrAuthenticationFailed
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
		return []user.User{}, ErrUnauthorized
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
		return user.User{}, ErrUnauthorized
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
		return user.User{}, ErrUnauthorized
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
		return ErrUnauthorized
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
		return ErrUnauthorized
	}

	return uc.Repo.Delete(id)
}
