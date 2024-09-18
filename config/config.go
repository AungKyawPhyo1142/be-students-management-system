package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbString := os.Getenv("DATABASE_URL")
	if dbString == "" {
		log.Fatal("Database url is not set yet")
	}

	// connect to datbase
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{})

	if err != nil {
		log.Fatalf("Database connection error: %v", err)
		return
	}

	log.Print("Database connected")

	DB = db.Debug()
}
