package user

import (
	"time"
)

// Session is an entity for user sessions.
type Session struct {
	User    User
	Expires time.Time
}

// NewSession returns an initialized user session.
func NewSession(u User, expires time.Time) Session {
	return Session{User: u, Expires: expires}
}

// Authorized returns true if the session's user has one of the
// provided roles.
func (s Session) UserHasRole(roles ...string) bool {
	for _, has := range s.User.Roles {
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
func (s Session) IsExpired() bool {
	return time.Now().After(s.Expires)
}

// IsValid returns true if the session is valid.
// The session is valid if it's not expires and
// contains a valid User entity for the
// authenticated user
func (s Session) IsValid() bool {
	return !s.IsExpired()
}
