package main

import (
	"testing"

	"github.com/atlantacoven/coven-platform/member-site/api"
	"github.com/atlantacoven/coven-platform/member-site/database"
)

func TestServer(t *testing.T) *api.TestServer {
	t.Helper()

	ctx := database.CreateTest(t)
	db := database.Get(ctx)
	return api.NewTestServer(t, ctx, NewServer(db))
}
