package user

import (
	"context"
	"fmt"
)

//  ctxKey is the type of the value for a context key.
type ctxKey int

// key is the key under which values are stored or retrieved.
const key ctxKey = 534783

// ContextWithSession returns a new context with the provided session encapsulated.
func ContextWithSession(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, key, s)
}

// GetSession returns the Session from the context.
// Returns an error is there's no or an invalid session in the context.
func GetSession(ctx context.Context) (Session, error) {
	s, ok := ctx.Value(key).(Session)
	if !ok {
		return Session{}, fmt.Errorf("session missing from context")
	}

	return s, nil
}
