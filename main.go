package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "modernc.org/sqlite"
)

func main() {
	// initialize database or use existent
	initDb()
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal("Ошибка соединения с БД")
	}
	defer db.Close()

	// router
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	r.Get("/api/nextdate", NextDateHandler)
	r.Post("/api/task", NewTaskHandler(db))
	r.Get("/api/tasks", GetTasksHandler(db))

	// launch server
	err = http.ListenAndServe(":7540", r) // TODO: get port from env variables
	if err != nil {
		panic(err)
	}
}
