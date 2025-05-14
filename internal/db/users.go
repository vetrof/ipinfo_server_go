package db

import (
	"errors"
	"ip_info_server/internal/models"
)

// GetUserByToken получает пользователя по токену
func GetUserByToken(token string) (models.User, error) {
	var user models.User
	row := DB.QueryRow("SELECT id, username, password, token FROM users WHERE token = ?", token)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Token)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// CreateUser создает нового пользователя
func CreateUser(username, password string) (models.User, error) {
	// Проверка существования пользователя
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return models.User{}, err
	}
	if count > 0 {
		return models.User{}, errors.New("username already exists")
	}

	// Генерация токена (в реальном приложении нужно использовать более надежный способ)
	token := GenerateToken(username, password)

	// Создание пользователя
	stmt, err := DB.Prepare("INSERT INTO users (username, password, token) VALUES (?, ?, ?)")
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password, token)
	if err != nil {
		return models.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:       int(id),
		Username: username,
		Password: password,
		Token:    token,
	}, nil
}

// GetUser получает пользователя по логину и паролю
func GetUser(username, password string) (models.User, error) {
	var user models.User
	row := DB.QueryRow("SELECT id, username, password, token FROM users WHERE username = ? AND password = ?", username, password)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Token)
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}
	return user, nil
}

// GenerateToken генерирует простой токен (для демонстрации)
func GenerateToken(username, password string) string {
	// В реальном приложении здесь должна быть более сложная логика
	return username + ":" + password
}
