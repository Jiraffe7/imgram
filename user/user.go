package user

import (
	"context"
)

// User contains data of a user
type User struct {
	ID uint64
}

type key int

const userKey key = iota

// NewContext creates a context with a User
func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext retrieves a User from context if exists
func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}
