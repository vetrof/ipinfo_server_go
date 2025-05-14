package services

import (
	"context"
	"encoding/json"
	"io"
	"ip_info_server/internal/db"
	"ip_info_server/internal/interfaces"
	"ip_info_server/internal/middleware"
	"ip_info_server/internal/models"
	"net/http"
)

// IPInfoServiceImpl реализует интерфейс сервиса для работы с IP информацией
type IPInfoServiceImpl struct {
	repository interfaces.IPInfoRepository
}

// Проверка соответствия интерфейсу на этапе компиляции
var _ interfaces.IPInfoService = (*IPInfoServiceImpl)(nil)

// NewIPInfoService создает новый экземпляр сервиса для работы с IP
func NewIPInfoService() *IPInfoServiceImpl {
	return &IPInfoServiceImpl{
		repository: db.GetRepository(),
	}
}

// GetUserHistory получает историю IP запросов пользователя
func (s *IPInfoServiceImpl) GetUserHistory(userID int) ([]models.IPInfo, error) {
	return s.repository.GetIPInfoByUser(userID)
}

// FetchSelfIPInfo получает информацию о своем IP
func (s *IPInfoServiceImpl) FetchSelfIPInfo(ctx context.Context) (models.IPInfo, error) {
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return models.IPInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.IPInfo{}, err
	}

	var ipInfo models.IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return models.IPInfo{}, err
	}

	// Получаем userID из контекста
	if userID, ok := middleware.GetUserID(ctx); ok {
		ipInfo.UserID = userID
		// Сохраняем в БД
		_ = s.repository.SaveIPInfo(ipInfo)
	}

	return ipInfo, nil
}

// FetchExternalIPInfo получает информацию о внешнем IP
func (s *IPInfoServiceImpl) FetchExternalIPInfo(ctx context.Context, ip string) (models.IPInfo, error) {
	resp, err := http.Get("https://ipinfo.io/" + ip + "/json")
	if err != nil {
		return models.IPInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.IPInfo{}, err
	}

	var ipInfo models.IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return models.IPInfo{}, err
	}

	// Получаем userID из контекста
	if userID, ok := middleware.GetUserID(ctx); ok {
		ipInfo.UserID = userID
		// Сохраняем в БД
		_ = s.repository.SaveIPInfo(ipInfo)
	}

	return ipInfo, nil
}
