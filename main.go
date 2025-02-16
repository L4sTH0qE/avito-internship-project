package main

import (
	"awesomeProject/controllers"
	"awesomeProject/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.Use(services.JwtAuthentication) // Добавляем middleware проверки JWT-токена.

	// Инициализация эндпоинтов.
	initializeRoutes(router)

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeRoutes(router *mux.Router) {
	router.HandleFunc("/api/info", controllers.GetInfo).Methods("GET")
	router.HandleFunc("/api/sendCoin", controllers.SendCoins).Methods("POST")
	router.HandleFunc("/api/buy/{item}", controllers.BuyItem).Methods("GET")
	router.HandleFunc("/api/auth", controllers.Authenticate).Methods("POST")
}
