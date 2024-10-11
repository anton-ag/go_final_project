package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func initDb() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
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
