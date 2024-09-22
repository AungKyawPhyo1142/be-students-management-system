package migrations

import (
	"log"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
)

func Migrate() {

	db := config.DB

	// check current databse
	currentDB := db.Migrator().CurrentDatabase()

	log.Printf("current database %v", currentDB)

	// make sure it is the correct database i am trying to migrate
	if currentDB != "student-management-system-db" {
		log.Fatalf("Database is not correct: %v", currentDB)
	}

	err := config.DB.AutoMigrate(&models.User{}, &models.Student{}, &models.Class{}, &models.StudentClass{})

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Printf("current database %v", config.DB.Migrator().CurrentDatabase())
	log.Print("Migration success!")

}
