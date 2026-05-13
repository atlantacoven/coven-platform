package users

import (
	"context"
	"fmt"
)

var ErrNotFound = fmt.Errorf("user not found")
var ErrInvalidPassword = fmt.Errorf("password invalid")

// Authenticate takes login credentials and, if valid, returns
// the user instance. Returns ErrInvalidPassword if user is not found
// or password is invalid.
func AuthenticatePassword(ctx context.Context, email, password string) (*User, error) {
	u, err := findByEmail(ctx, email)
	if err != nil {
		// TODO: email existence could be tested by timing attack
		return nil, fmt.Errorf("db: %w", err)
	}
	if !u.VerifyPassword(password) {
		return nil, ErrInvalidPassword
	}
	return u, nil
}
