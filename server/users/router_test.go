package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"rabidaudio.com/coven-door/server/app"
)

func TestPostSession(t *testing.T) {
	a := app.NewTest(t, Router)

	u1 := UserFixture(t, a.Context)

	// test find
	res, err := a.Request("POST", "/session", `{"email":"julien@example.com","password":"password"}`)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, u1.Id, int(res.Data["id"].(float64)))
}
