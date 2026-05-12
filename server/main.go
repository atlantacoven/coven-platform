package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"rabidaudio.com/coven-door/server/users"
)

type RouteBuilder func(r chi.Router)

var routers = []RouteBuilder{
	// Add route handlers here
	users.Router,
}

func main() {
	db, err := sqlx.Open("sqlite3", "development.db")
	if err != nil {
		panic(fmt.Errorf("open db: %w", err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("db ping: %w", err))
	}

	r := chi.NewRouter()
	for _, rb := range routers {
		rb(r)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Starting server on %v", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(fmt.Errorf("start server: %v", err))
	}
}
