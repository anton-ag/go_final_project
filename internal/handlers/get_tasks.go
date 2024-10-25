package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/anton-ag/todolist/internal/database"
)

func GetTasks(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// TODO: add search parameter
		tasks, err := database.GetTasks(db)
		if err != nil {
			respondError(w, err.Error())
			return
		}
		body, _ := json.Marshal(TasksResponse{Tasks: tasks})
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}
