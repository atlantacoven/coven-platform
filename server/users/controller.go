package users

import (
	"context"
	"fmt"
)

var ErrInvalidPassword = fmt.Errorf("Password Invalid")

// Authenticate takes login credentials and, if valid, returns
// the user instance. Returns ErrInvalidPassword if user is not found
// or password is invalid.
func AuthenticatePassword(ctx context.Context, email, password string) (*User, error) {
	u, err := findByEmail(ctx, email)
	if err != nil {
		return nil, err // TODO: email existence could be tested by timing attack
	}
	if !u.VerifyPassword(password) {
		return nil, ErrInvalidPassword
	}
	return u, nil
}
