package middleware

import (
	"context"
	"ip_info_server/internal/db"
	"ip_info_server/internal/interfaces"
	"net/http"
	"strings"
)

// UserIDKey ключ контекста для хранения ID пользователя
type UserIDKey string

const userIDKey UserIDKey = "userID"

// userRepository репозиторий для работы с пользователями
var userRepository interfaces.UserRepository

// InitMiddleware инициализирует middleware
func InitMiddleware() {
	userRepository = db.GetRepository()
}

// AuthMiddleware проверяет авторизацию пользователя
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		user, ok := userRepository.GetUserByToken(token)
		if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем ID пользователя в контекст
		ctx := context.WithValue(r.Context(), userIDKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID получает ID пользователя из контекста
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}
