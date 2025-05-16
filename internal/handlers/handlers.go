package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"ip_info_server/internal/middleware"
	"ip_info_server/internal/service"
	"log"
	"net/http"
)

type Handler struct {
	ipService   *service.IPService
	userService *service.UserService
}

func NewHandler(ipService *service.IPService, userService *service.UserService) *Handler {
	return &Handler{
		ipService:   ipService,
		userService: userService,
	}
}

func (h *Handler) SelfIpHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	ipInfo, err := h.ipService.GetSelfIPInfo(userID)
	if err != nil {
		http.Error(w, "Failed to get IP info: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error fetching IP info:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ipInfo); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func (h *Handler) ExtIpHandler(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	ipInfo, err := h.ipService.GetExternalIPInfo(ip, userID)
	if err != nil {
		http.Error(w, "Failed to get IP info: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error fetching IP info:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ipInfo); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func (h *Handler) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	records, err := h.ipService.GetIPHistory(userID)
	if err != nil {
		http.Error(w, "Failed to fetch history: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.RegisterUser(username, password)
	if err != nil {
		http.Error(w, "Cannot create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": user.Token,
	})
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.LoginUser(username, password)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": user.Token,
	})
}
