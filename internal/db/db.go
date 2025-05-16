package db

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	log.Println("Using DB path:", filepath.Join(".", "sqlite3.db"))
	DB, err = sql.Open("sqlite3", filepath.Join(".", "sqlite3.db"))
	if err != nil {
		log.Fatal("Cannot open database:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS ip_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT,
		hostname TEXT,
		city TEXT,
		region TEXT,
		country TEXT,
		loc TEXT,
		org TEXT,
		postal TEXT,
		timezone TEXT,
		readme TEXT,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Cannot create table:", err)
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		token TEXT UNIQUE
	);`

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Cannot create users table:", err)
	}
}
