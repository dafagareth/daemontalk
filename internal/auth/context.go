package auth

import "context"

type contextKey struct{}

var userContextKey = contextKey{}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUser(ctx context.Context) *User {
	if ctx == nil {
		return nil
	}
	if user, ok := ctx.Value(userContextKey).(*User); ok {
		return user
	}
	return nil
}
