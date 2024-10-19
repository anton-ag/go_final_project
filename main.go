package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	// initialize database or use existent
	initDb()

	// router
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	r.Get("/api/nextdate", NextDateHandler)

	// launch server
	err := http.ListenAndServe(":7540", r) // TODO: get port from env variables
	if err != nil {
		panic(err)
	}
}
