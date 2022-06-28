// Package usecases provides user usecases with application logic.
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
	u, err := uc.Repo.QueryByEmail(ctx, email)
	if err != nil {
		return user.Session{}, ErrAuthenticationFailed
	}

	if !u.HasPassword(password) {
		return user.Session{}, ErrAuthenticationFailed
	}

	return user.NewSession(u, now.Add(uc.SessionDuration)), nil
}

// Create inserts the provided user at the repository.
func (uc UserUseCases) Create(ctx context.Context, n NewUser, now time.Time) (user.User, error) {
	s, err := user.GetSession(ctx)
	if err != nil {
		return user.User{}, err
	}

	// Only ADMIN can do this action.
	if !s.UserHasRole(user.RoleAdmin) {
		return user.User{}, ErrAuthenticationFailed
	}

	u, err := user.NewWithID(n.Name, n.Email, n.Password, n.Roles, now)
	if err != nil {
		return user.User{}, err
	}

	if err := uc.Repo.Create(ctx, u); err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Register inserts the provided user at the repository.
// If there are no users yet, the user will get the ADMIN and USER role,
// in any other case the user will get only the USER role.
//
// Unlike Create, Register requires no admin priviliges. It is meant
// to register the first admin of the system and user signups.
func (uc UserUseCases) Register(ctx context.Context, n NewUser, now time.Time) (user.User, error) {
	users, err := uc.Repo.Query(ctx, 1, 1)
	if err != nil {
		return user.User{}, err
	}
	if len(users) == 0 {
		n.Roles = []string{user.RoleAdmin, user.RoleUser}
	} else {
		n.Roles = []string{user.RoleUser}
	}

	u, err := user.NewWithID(n.Name, n.Email, n.Password, n.Roles, now)
	if err != nil {
		return user.User{}, err
	}

	if err := uc.Repo.Create(ctx, u); err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Query retrieves all users from the repository,
func (uc UserUseCases) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]user.User, error) {
	s, err := user.GetSession(ctx)
	if err != nil {
		return []user.User{}, err
	}

	// Only ADMIN can do this action.
	if !s.UserHasRole(user.RoleAdmin) {
		return []user.User{}, ErrUnauthorized
	}

	return uc.Repo.Query(ctx, pageNumber, rowsPerPage)
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

	return uc.Repo.QueryByID(ctx, id)
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

	return uc.Repo.QueryByEmail(ctx, email)
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

	return uc.Repo.Update(ctx, u)
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

	return uc.Repo.Delete(ctx, id)
}
