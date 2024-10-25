package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/anton-ag/todolist/internal/config"
	"github.com/go-chi/chi/v5"

	_ "modernc.org/sqlite"
)

func main() {
	var config config.Config
	config.Init()

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
	r.Get("/api/task", GetTaskHandler(db))
	r.Put("/api/task", PutTaskHandler(db))
	r.Post("/api/task/done", PostDoneHandler(db))
	r.Delete("/api/task", DeleteTaskHandler(db))

	// launch server
	err = http.ListenAndServe(config.Port, r)
	if err != nil {
		panic(err)
	}
}
