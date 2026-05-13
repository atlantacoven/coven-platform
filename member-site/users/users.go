package users

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int    `db:"id"`
	Name              string `db:"name"`
	Email             string `db:"email"`
	EncryptedPassword string `db:"encrypted_password"`
}

func (u *User) SetPassword(p string) error {
	enc, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("set password: %w", err)
	}
	u.EncryptedPassword = string(enc)
	return nil
}

// VerifyPassword checks if p matches EncryptedPassword
// in a
func (u *User) VerifyPassword(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(p))
	if err == nil {
		return true
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	panic(fmt.Errorf("password compare: %w", err))
}

// GenerateSecret generates a secret signed by the server that can be used by the
// app to unlock the door.
// func GenerateSecret() (string, error) {

// }
