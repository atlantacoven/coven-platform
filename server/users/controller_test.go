package users

import (
	"testing"

	"github.com/atlantacoven/door-lock/server/database"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticatePassword(t *testing.T) {
	ctx := database.CreateTest(t)
	u1 := Fixture(t, ctx)

	// test find
	u2, err := AuthenticatePassword(ctx, "julien@example.com", "password")
	assert.NoError(t, err)
	assert.EqualValues(t, u1, u2)

	// test not found
	u3, err := AuthenticatePassword(ctx, "no-one@example.com", "password")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, u3)

	// test wrong password
	u4, err := AuthenticatePassword(ctx, "julien@example.com", "h@ck3rz")
	assert.ErrorIs(t, err, ErrInvalidPassword)
	assert.Nil(t, u4)
}
