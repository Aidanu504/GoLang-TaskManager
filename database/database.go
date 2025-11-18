package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3" // _ allows for indirect usage
)

// basic connection function from database/sql package docs
func DatabaseConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/TaskDatabase.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
