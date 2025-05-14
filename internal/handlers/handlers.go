package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"ip_info_server/internal/db"
	"ip_info_server/internal/interfaces"
	"ip_info_server/internal/middleware"
	"ip_info_server/internal/services"
	"net/http"
)

// ipInfoService сервис для работы с IP информацией
var ipInfoService interfaces.IPInfoService

// userRepository репозиторий для работы с пользователями
var userRepository interfaces.UserRepository

// Init инициализирует обработчики
func Init() {
	ipInfoService = services.NewIPInfoService()
	userRepository = db.GetRepository()
}

// SelfIpHandler обрабатывает запрос получения информации о своем IP
func SelfIpHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo, err := ipInfoService.FetchSelfIPInfo(r.Context())
	if err != nil {
		http.Error(w, "Failed to get IP info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ipInfo); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// ExtIpHandler обрабатывает запрос получения информации о внешнем IP
func ExtIpHandler(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")

	ipInfo, err := ipInfoService.FetchExternalIPInfo(r.Context(), ip)
	if err != nil {
		http.Error(w, "Failed to get IP info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ipInfo); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// HistoryHandler обрабатывает запрос получения истории запросов
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	records, err := ipInfoService.GetUserHistory(userID)
	if err != nil {
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// RegisterHandler обрабатывает запрос регистрации пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	user, err := userRepository.CreateUser(username, password)
	if err != nil {
		http.Error(w, "Cannot create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": user.Token,
	})
}

// LoginHandler обрабатывает запрос входа пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	user, err := userRepository.GetUser(username, password)
	if err != nil {
		http.Error(w, "Cannot login: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": user.Token,
	})
}
