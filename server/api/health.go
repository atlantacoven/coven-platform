package api

import (
	"net/http"
	"time"
)

type healthResponseBody struct {
	Environment Environment `json:"environment"`
	CurrentTime time.Time   `json:"current_time"`
}

// HealthCheck implements a simple endpoint for checking if the server is running.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	Respond(w, &healthResponseBody{
		Environment: Env(),
		CurrentTime: time.Now(),
	}, "OK")
}
