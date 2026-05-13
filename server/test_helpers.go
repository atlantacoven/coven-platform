package main

import (
	"testing"

	"rabidaudio.com/coven-door/server/api"
	"rabidaudio.com/coven-door/server/database"
)

func TestServer(t *testing.T) *api.TestServer {
	t.Helper()

	ctx := database.PrepareForTest(t)
	db := database.Get(ctx)
	return api.NewTestServer(t, ctx, NewServer(db))
}
