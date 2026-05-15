package main

import (
	"net/http"
	"time"

	"github.com/atlantacoven/coven-platform/member-site/api"
)

type healthResponseBody struct {
	Environment api.Environment `json:"environment"`
	CurrentTime time.Time       `json:"current_time"`
}

// HealthCheck implements a simple endpoint for checking if the server is running.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	api.Respond(w, &healthResponseBody{
		Environment: api.Env(),
		CurrentTime: time.Now(),
	}, "OK")
}
