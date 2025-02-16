package services

import (
	"awesomeProject/dao"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *gorm.DB

// Инициализация БД перед запуском приложения.
func init() {

	// Извлекаем значения из .env.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s", dbHost, dbUsername, dbName, dbPort, dbPassword)
	log.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		log.Fatal(err)
	}

	// Миграция БД.
	db = conn
	db.Debug().AutoMigrate(&dao.Merch{}, &dao.Purchase{}, &dao.Transaction{}, &dao.User{})

	// Инициализация таблицы Merch.
	for _, item := range dao.MerchList {
		var existing dao.Merch
		if err := db.Where("type = ?", item.Type).First(&existing).Error; err != nil {
			db.Create(&item)
		}
	}
}

func GetDB() *gorm.DB {
	return db
}
