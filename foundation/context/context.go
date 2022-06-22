package context

import (
	"context"
	"errors"
	"fmt"
	"time"
)

//  ctxKey is the type of the value for a context key.
type ctxKey int

// key is the key under which values are stored or retrieved.
const key ctxKey = 1

// Values represent a request's state.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// GetValues returns the Values from the context.
// Returns an error is there's no or an invalid Values in the context.
func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return nil, fmt.Errorf("values missing from context")
	}

	return v, nil
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return errors.New("values missing from context")
	}
	v.StatusCode = statusCode
	return nil
}

// ContextWithSession returns a new context with the provided values encapsulated.
func ContextWithValues(ctx context.Context, v *Values) context.Context {
	return context.WithValue(ctx, key, v)
}
