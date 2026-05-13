package main

import (
	"testing"

	"github.com/atlantacoven/coven-platform/member-site/users"
	"github.com/stretchr/testify/assert"
)

func TestPostSession(t *testing.T) {
	s := TestServer(t)
	u1 := users.Fixture(t, s)

	// test find
	res, err := s.Request("POST", "/session", `{"email":"julien@example.com","password":"password"}`)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, u1.Id, int(res.Data["id"].(float64)))
}
