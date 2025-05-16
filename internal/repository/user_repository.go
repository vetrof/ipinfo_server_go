package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"ip_info_server/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username, password string) (*models.User, error) {
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	stmt, err := r.db.Prepare("INSERT INTO users (username, password, token) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(username, password, token)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       int(id),
		Username: username,
		Password: password,
		Token:    token,
	}, nil
}

func (r *UserRepository) GetUserByCredentials(username, password string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, token FROM users WHERE username = ? AND password = ?",
		username, password,
	).Scan(&user.ID, &user.Token)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}

	user.Username = username
	user.Password = password

	return &user, nil
}

func (r *UserRepository) GetUserByToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, username FROM users WHERE token = ?", token).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}

	user.Token = token
	return &user, nil
}

func (r *UserRepository) GetUserIDByToken(token string) (int, error) {
	var userID int
	err := r.db.QueryRow("SELECT id FROM users WHERE token = ?", token).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid token")
		}
		return 0, err
	}
	return userID, nil
}
