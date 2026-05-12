package users

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int
	Name              string
	Email             string
	encryptedPassword []byte
}

func (u *User) SetPassword(p string) error {
	enc, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("set password: %w", err)
	}
	u.encryptedPassword = enc
	return nil
}

// VerifyPassword checks if p matches EncryptedPassword
// in a
func (u *User) VerifyPassword(p string) bool {
	err := bcrypt.CompareHashAndPassword(u.encryptedPassword, []byte(p))
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
