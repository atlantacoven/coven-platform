package main

import (
	"testing"

	"github.com/atlantacoven/door-lock/server/api"
	"github.com/atlantacoven/door-lock/server/database"
)

func TestServer(t *testing.T) *api.TestServer {
	t.Helper()

	ctx := database.CreateTest(t)
	db := database.Get(ctx)
	return api.NewTestServer(t, ctx, NewServer(db))
}
