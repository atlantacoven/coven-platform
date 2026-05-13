package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"rabidaudio.com/coven-door/server/api"
	"rabidaudio.com/coven-door/server/database"
	"rabidaudio.com/coven-door/server/users"
)

type RouteBuilder func(r chi.Router)

var routers = []RouteBuilder{
	// Add route handlers here
	users.Router,
}

func main() {
	db, err := database.Create()
	if err != nil {
		panic(fmt.Errorf("open db: %w", err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("db ping: %w", err))
	}

	r := chi.NewRouter()

	// attach db to context
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := database.WithDB(db, r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	// health check
	r.Get("/", api.HealthCheck)
	// routes from other packages
	for _, rb := range routers {
		rb(r)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Starting server on %v\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(fmt.Errorf("start server: %v", err))
	}
}
