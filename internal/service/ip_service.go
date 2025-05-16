package service

import (
	"encoding/json"
	"io"
	"ip_info_server/internal/models"
	"ip_info_server/internal/repository"
	"net/http"
)

type IPService struct {
	ipRepo *repository.IPRepository
}

func NewIPService(ipRepo *repository.IPRepository) *IPService {
	return &IPService{ipRepo: ipRepo}
}

func (s *IPService) GetSelfIPInfo(userID int) (models.IPInfo, error) {
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

	ipInfo.UserID = userID

	// Сохраняем в БД
	if err := s.ipRepo.SaveIPInfo(ipInfo); err != nil {
		return models.IPInfo{}, err
	}

	return ipInfo, nil
}

func (s *IPService) GetExternalIPInfo(ip string, userID int) (models.IPInfo, error) {
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

	ipInfo.UserID = userID

	// Сохраняем в БД
	if err := s.ipRepo.SaveIPInfo(ipInfo); err != nil {
		return models.IPInfo{}, err
	}

	return ipInfo, nil
}

func (s *IPService) GetIPHistory(userID int) ([]models.IPInfo, error) {
	return s.ipRepo.GetHistoryByUserID(userID)
}
