package db

import (
	"ip_info_server/internal/interfaces"
	"ip_info_server/internal/models"
)

// Проверка соответствия интерфейсу на этапе компиляции
var _ interfaces.IPInfoRepository = (*SQLiteRepository)(nil)
var _ interfaces.UserRepository = (*SQLiteRepository)(nil)

// SQLiteRepository реализует интерфейсы репозиториев
type SQLiteRepository struct{}

// GetRepository возвращает экземпляр репозитория
func GetRepository() *SQLiteRepository {
	return &SQLiteRepository{}
}

// SaveIPInfo сохраняет информацию об IP в БД
func (r *SQLiteRepository) SaveIPInfo(info models.IPInfo) error {
	return SaveIPInfo(info)
}

// GetIPInfoByUser получает историю IP запросов пользователя
func (r *SQLiteRepository) GetIPInfoByUser(userID int) ([]models.IPInfo, error) {
	return HistoryIPInfoByUser(userID)
}

// CreateUser создает нового пользователя
func (r *SQLiteRepository) CreateUser(username, password string) (models.User, error) {
	return CreateUser(username, password)
}

// GetUser получает пользователя по логину и паролю
func (r *SQLiteRepository) GetUser(username, password string) (models.User, error) {
	return GetUser(username, password)
}

// GetUserByToken получает пользователя по токену
func (r *SQLiteRepository) GetUserByToken(token string) (models.User, bool) {
	user, err := GetUserByToken(token)
	if err != nil {
		return models.User{}, false
	}
	return user, true
}
