package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"rabidaudio.com/coven-door/server/database"
)

func UserFixture(t *testing.T, ctx context.Context) *User {
	t.Helper()

	u := User{Name: "Julien", Email: "julien@example.com"}
	err := u.SetPassword("password")
	assert.NoError(t, err)
	err = create(ctx, &u)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, u.Id)
	return &u
}

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
