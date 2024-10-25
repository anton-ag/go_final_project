package database

import (
	"database/sql"
	"strings"

	"github.com/anton-ag/todolist/internal/models"
)

func GetTasks(db *sql.DB, search string) ([]models.Task, error) {
	var task models.Task
	var tasks []models.Task

	search = strings.Join([]string{"%", search, "%"}, "")

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
	rows, err := db.Query(
		query,
		sql.Named("limit", models.Limit),
		sql.Named("search", search),
	)

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
		return []models.Task{}, err
	}
	if len(tasks) == 0 {
		return []models.Task{}, err
	}
	return tasks, nil
}
