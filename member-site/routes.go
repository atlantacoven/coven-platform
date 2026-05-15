package main

import (
	"net/http"

	"github.com/atlantacoven/coven-platform/member-site/database"
	"github.com/atlantacoven/coven-platform/member-site/users"
	"github.com/go-chi/chi/v5"
)

type RouteBuilder func(r chi.Router)

var routers = []RouteBuilder{
	// ADD ROUTES HERE
	users.Router,
}

func NewServer(db database.DB) http.Handler {
	r := chi.NewRouter()

	// attach db to context
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := database.WithDB(db, r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// health check
	r.Get("/", HealthCheck)
	// routes from other packages
	for _, rb := range routers {
		rb(r)
	}
	return r
}
