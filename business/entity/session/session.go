package session

import (
	"time"

	"github.com/appinesshq/caservice/business/entity/user"
)

// Session is an entity that contains session data and
// related business logic.
type Session struct {
	user    *user.User
	expires time.Time
}

// New returns a pointer to an instantiated session.
func New(u *user.User, expires time.Time) *Session {
	return &Session{
		user:    u,
		expires: expires,
	}
}

// Authorized returns true if the session's user has one of the
// provided roles.
func (s *Session) UserHasRole(roles ...user.Role) bool {
	if s.user == nil {
		return false
	}

	for _, has := range s.user.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}
	return false
}

// IsExpired returns true if the session
// is expired.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.expires)
}

// IsValid returns true if the session is valid.
// The session is valid if it's not expires and
// contains a valid User entity for the
// authenticated user
func (s *Session) IsValid() bool {
	return !s.IsExpired() && s.user != nil
}
