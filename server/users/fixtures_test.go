package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
