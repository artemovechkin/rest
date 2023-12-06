package db

import (
	"database/sql"
	"log"
)

func InitDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "TaskManager.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
			CREATE TABLE if not exists taskmanager (
	   			ID INTEGER PRIMARY KEY AUTOINCREMENT,
	   			Title VARCHAR(255),
			    Description VARCHAR(255),
			    Status INTEGER,
			    Priority VARCHAR(255)
	);
`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
