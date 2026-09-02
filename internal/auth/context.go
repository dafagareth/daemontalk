package auth

import "context"

type contextKey struct{}

var userContextKey = contextKey{}

// WithUser returns a new Context that carries the provided User pointer.
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUser returns the User from the Context if present, otherwise nil.
func GetUser(ctx context.Context) *User {
	if ctx == nil {
		return nil
	}
	if user, ok := ctx.Value(userContextKey).(*User); ok {
		return user
	}
	return nil
}
