package main

import (
	"fmt"
	"net/http"
	"os"

	"rabidaudio.com/coven-door/server/app"
	"rabidaudio.com/coven-door/server/database"
	"rabidaudio.com/coven-door/server/users"
)

var routers = []app.RouteBuilder{
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

	r := app.New(db, routers...)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := "localhost:" + port
	fmt.Printf("Starting server on %v\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(fmt.Errorf("start server: %v", err))
	}
}
