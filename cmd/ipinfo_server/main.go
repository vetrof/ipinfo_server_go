package main

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"ip_info_server/internal/db"
	"ip_info_server/internal/handlers"
	"ip_info_server/internal/middleware"
	"ip_info_server/internal/repository"
	"ip_info_server/internal/service"
	"log"
	"net/http"
)

func main() {
	// DB init
	db.InitDB()
	defer db.DB.Close()

	// Repositories
	ipRepo := repository.NewIPRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	// Services
	ipService := service.NewIPService(ipRepo)
	userService := service.NewUserService(userRepo)

	// Handlers
	handler := handlers.NewHandler(ipService, userService)

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Router init
	router := chi.NewRouter()
	router.Use(chiMiddleware.Logger)

	// Public paths
	router.Post("/register", handler.RegisterHandler)
	router.Get("/login", handler.LoginHandler)

	// Protected paths
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.TokenAuthMiddleware)
		r.Get("/self_ip", handler.SelfIpHandler)
		r.Get("/ext_ip/{ip}", handler.ExtIpHandler)
		r.Get("/history", handler.HistoryHandler)
	})

	// Start server
	log.Println("Starting server on localhost:8080")
	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
