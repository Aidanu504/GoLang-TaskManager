package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3" // _ allows for indirect usage
)

// basic connection function from database/sql package docs
func DatabaseConnect() *sql.DB {
	db, err := sql.Open("sqlite", "./database/TaskDatabase.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}


// Creates tables if DB file is empty or missing
func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS Tasks (
		TaskID INTEGER PRIMARY KEY AUTOINCREMENT,
		TaskName TEXT NOT NULL,
		TaskDescription TEXT NOT NULL,
		IsCompleted BOOLEAN NOT NULL DEFAULT 0,
		CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	return err
}