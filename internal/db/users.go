package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
)

type User struct {
	ID       int
	Username string
	Password string
	Token    string
}

func CreateUser(username, password string) (*User, error) {
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	stmt, err := DB.Prepare("INSERT INTO users (username, password, token) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, password, token)
	if err != nil {
		return nil, err
	}

	return &User{Username: username, Password: password, Token: token}, nil
}

func GetUser(username, password string) (*User, error) {
	var token string
	err := DB.QueryRow(
		"SELECT token FROM users WHERE username = ? AND password = ?",
		username, password,
	).Scan(&token)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}

	return &User{Username: username, Password: password, Token: token}, nil
}

func GetUserByToken(token string) (*User, error) {
	row := DB.QueryRow("SELECT id, username FROM users WHERE token = ?", token)
	var user User
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	return &user, nil
}
