package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"rabidaudio.com/coven-door/server/database"
)

func TestAuthenticatePassword(t *testing.T) {
	ctx := database.PrepareForTest(t)
	u1 := UserFixture(t, ctx)

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
