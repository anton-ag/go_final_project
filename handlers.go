package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type IDResponse struct {
	ID string `json:"id"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func respondError(w http.ResponseWriter, s string) {
	body, _ := json.Marshal(ErrorResponse{Error: s})
	w.WriteHeader(http.StatusBadRequest)
	w.Write(body)
}

func respondOk(w http.ResponseWriter, s string) {
	body, _ := json.Marshal(IDResponse{ID: s})
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(DateFormat, now)
	if err != nil {
		http.Error(w, "Некорректный формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(nextDate))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewTaskHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		now := time.Now()
		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			respondError(w, "Ошибка подключения к базе данных")
			return
		}
		defer db.Close()

		var task Task
		var buf bytes.Buffer

		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			respondError(w, "Неверный формат запроса")
			return
		}

		if task.Title == "" {
			respondError(w, "Не указан заголовок задачи")
			return
		}

		if task.Date == "" {
			task.Date = now.Format(DateFormat)
		}

		if _, err = time.Parse(DateFormat, task.Date); err != nil {
			respondError(w, "Неверный формат времени")
			return
		}

		if task.Date < now.Format(DateFormat) {
			task.Date = now.Format(DateFormat)
		}

		if task.Repeat != "" {
			_, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				respondError(w, "Неверный формат правила повтора")
				return
			}
		}

		id, err := addTask(db, task)
		if err != nil {
			respondError(w, "Ошибка работы с БД")
			return
		}

		idS := strconv.Itoa(int(id))
		respondOk(w, idS)
	}
}

func GetTasksHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			respondError(w, "Ошибка подключения к базе данных")
			return
		}
		defer db.Close()

		// TODO: add search parameter
		tasks, err := getTasks(db)
		if err != nil {
			respondError(w, err.Error())
			return
		}
		body, _ := json.Marshal(TasksResponse{Tasks: tasks})
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func GetTaskHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			respondError(w, "Пустой ID")
			return
		}

		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			respondError(w, "Ошибка подключения к базе данных")
			return
		}
		defer db.Close()

		task, err := getTaskByID(db, id)
		if err != nil {
			respondError(w, "Задача с данным ID не найдена")
			return
		}

		body, _ := json.Marshal(task)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func PutTaskHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		now := time.Now()
		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			respondError(w, "Ошибка подключения к базе данных")
			return
		}
		defer db.Close()

		var task Task
		var buf bytes.Buffer

		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			respondError(w, "Неверный формат запроса")
			return
		}

		_, err = strconv.Atoi(task.ID)
		if err != nil {
			respondError(w, "Неверный ID")
			return
		}

		if task.Title == "" {
			respondError(w, "Не указан заголовок задачи")
			return
		}

		if _, err = time.Parse(DateFormat, task.Date); err != nil {
			respondError(w, "Неверный формат времени")
			return
		}

		// TODO: add function to check the correct format
		if task.Repeat != "" {
			_, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				respondError(w, "Неверный формат правила повтора")
				return
			}
		}

		id, err := updateTask(db, task)
		if err != nil {
			respondError(w, err.Error())
			return
		}

		respondOk(w, strconv.Itoa(int(id)))
	}
}
