package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"rabidaudio.com/coven-door/server/database"
)

type RouteBuilder func(r chi.Router)

func New(db database.DB, routers ...RouteBuilder) http.Handler {
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
