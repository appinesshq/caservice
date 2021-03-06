// Package mem provides memory storage functionality for users.
package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/appinesshq/caservice/business/user"
	"github.com/appinesshq/caservice/business/user/usecases"
)

type Store struct {
	mu      sync.RWMutex
	users   map[string]user.User
	indexes map[string]string
}

func New() *Store {
	return &Store{users: make(map[string]user.User), indexes: make(map[string]string)}
}

func (m *Store) hasID(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.users[id]
	return exists
}

func (m *Store) hasEmail(email string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.indexes["email:"+email]
	return exists
}

func (m *Store) Create(ctx context.Context, u user.User) error {
	// Return an error if the ID or Email already exist.
	if m.hasID(u.ID) {
		return usecases.ErrUniqueID
	}
	if m.hasEmail(u.Email) {
		return usecases.ErrUniqueEmail
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Store the user and index for the email.
	m.users[u.ID] = u
	m.indexes["email:"+u.Email] = u.ID

	return nil
}

func (m *Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]user.User, error) {
	// TODO: Implement pagination.
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := []user.User{}
	for _, u := range m.users {
		users = append(users, u)
	}

	return users, nil
}

func (m *Store) QueryByID(ctx context.Context, id string) (user.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u, ok := m.users[id]
	if !ok {
		return user.User{}, usecases.ErrNotFound
	}

	return u, nil
}

func (m *Store) QueryByEmail(ctx context.Context, email string) (user.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u := user.User{}
	// Get the ID from the email index.
	id, ok := m.indexes["email:"+email]

	// Loop through the users if no id was indexed.
	if !ok {
		for _, usr := range m.users {
			// User found: Store in u and break the loop.
			if usr.Email == email {
				u = usr
				break
			}
		}
		// No user found. Return an error.
		return u, usecases.ErrNotFound
	} else {
		// Get the user data via the id.
		u = m.users[id]
	}

	return u, nil

}
func (m *Store) Update(ctx context.Context, u user.User) error {
	current, err := m.QueryByID(ctx, u.ID)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[u.ID] = u
	if u.Email != current.Email {
		delete(m.indexes, "email:"+current.Email)
		m.indexes[u.Email] = u.ID
	}

	return nil
}

func (m *Store) Delete(ctx context.Context, id string) error {
	current, err := m.QueryByID(ctx, id)
	if err != nil && errors.Is(err, usecases.ErrNotFound) {
		// Delete should return a nil error in case of not found.
		return nil
	} else if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.users, id)
	delete(m.indexes, "email:"+current.Email)
	return nil
}
