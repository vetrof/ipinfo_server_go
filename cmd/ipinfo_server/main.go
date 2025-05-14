package main

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"ip_info_server/internal/db"
	"ip_info_server/internal/handlers"
	authMiddleware "ip_info_server/internal/middleware"
	"log"
	"net/http"
)

func main() {
	//db init
	db.InitDB()
	defer db.DB.Close()

	//инициализация компонентов
	authMiddleware.InitMiddleware()
	handlers.Init()

	//router init
	router := chi.NewRouter()
	router.Use(chiMiddleware.Logger)

	//public path
	router.Post("/register", handlers.RegisterHandler)
	router.Get("/login", handlers.LoginHandler)

	//with token path
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware)
		r.Get("/self_ip", handlers.SelfIpHandler)
		r.Get("/ext_ip/{ip}", handlers.ExtIpHandler)
		r.Get("/history", handlers.HistoryHandler)
	})

	//start server
	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
