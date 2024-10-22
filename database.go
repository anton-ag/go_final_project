package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const Limit = 50

func initDb() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db") // TODO: customize dbfile with env variable
	_, err = os.Stat(dbFile)

	var needSetup bool
	if err != nil {
		needSetup = true
	}
	if needSetup {
		os.Create(dbFile)
		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS scheduler (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `date` CHAR(8) NOT NULL DEFAULT '', `title` VARCHAR(128) NOT NULL DEFAULT '', `comment` VARCHAR(256) NOT NULL DEFAULT '', `repeat` VARCHAR(128) NOT NULL DEFAULT '');")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func addTask(db *sql.DB, task Task) (int64, error) {
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getTasks(db *sql.DB) ([]Task, error) {
	var task Task
	var tasks []Task

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", Limit))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Close(); err != nil {
		return []Task{}, err
	}
	if len(tasks) == 0 {
		return []Task{}, err
	}
	return tasks, nil
}
