package db

import (
	"database/sql"
	"ip_info_server/internal/models"
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

func SaveIPInfo(info models.IPInfo) error {
	stmt, err := DB.Prepare(`
		INSERT INTO ip_info (
			ip, hostname, city, region, country,
			loc, org, postal, timezone, readme, user_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		info.IP, info.Hostname, info.City, info.Region, info.Country,
		info.Loc, info.Org, info.Postal, info.Timezone, info.Readme,
		info.UserID,
	)

	return err
}

func HistoryIPInfoByUser(userID int) ([]models.IPInfo, error) {
	rows, err := DB.Query(`
		SELECT ip, hostname, city, region, country,
		       loc, org, postal, timezone, readme, user_id
		FROM ip_info
		WHERE user_id = ?
		ORDER BY id DESC
	`, userID)
	if err != nil {
		log.Println("DB query error:", err)
		return nil, err
	}
	defer rows.Close()

	var result []models.IPInfo
	for rows.Next() {
		var info models.IPInfo
		err := rows.Scan(
			&info.IP, &info.Hostname, &info.City, &info.Region, &info.Country,
			&info.Loc, &info.Org, &info.Postal, &info.Timezone, &info.Readme, &info.UserID,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		result = append(result, info)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		return nil, err
	}

	return result, nil
}
