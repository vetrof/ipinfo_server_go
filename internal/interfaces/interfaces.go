package interfaces

import (
	"context"
	"ip_info_server/internal/models"
)

// IPInfoRepository определяет методы для работы с IP информацией в БД
type IPInfoRepository interface {
	SaveIPInfo(info models.IPInfo) error
	GetIPInfoByUser(userID int) ([]models.IPInfo, error)
}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	CreateUser(username, password string) (models.User, error)
	GetUser(username, password string) (models.User, error)
	GetUserByToken(token string) (models.User, bool)
}

// IPInfoService определяет сервисные методы для работы с IP информацией
type IPInfoService interface {
	GetUserHistory(userID int) ([]models.IPInfo, error)
	FetchSelfIPInfo(ctx context.Context) (models.IPInfo, error)
	FetchExternalIPInfo(ctx context.Context, ip string) (models.IPInfo, error)
}
