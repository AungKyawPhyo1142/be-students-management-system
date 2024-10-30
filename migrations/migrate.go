package migrations

import (
	"log"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
	"github.com/go-gormigrate/gormigrate/v2"
)

func Migrate() {

	db := config.DB

	// check current databse
	currentDB := db.Migrator().CurrentDatabase()

	log.Printf("current database %v", currentDB)

	// make sure it is the correct database i am trying to migrate
	// if currentDB != "student-management-system-db" && currentDB != "railway" {
	// 	log.Fatalf("Database is not correct: %v", currentDB)
	// }

	err := config.DB.AutoMigrate(&models.Admin{}, &models.Student{}, &models.Class{}, &models.StudentClass{}, &models.Instructor{}, &models.InstructorClass{})

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// migration tracking
	migrations := []*gormigrate.Migration{
		{
			ID:       "20240928223921_add_image_column_to_students",
			Migrate:  Migrate_20240928223921_add_image_column_to_students,
			Rollback: Rollback_20240928223921_add_image_column_to_students,
		},
	}

	m := gormigrate.New(config.DB, gormigrate.DefaultOptions, migrations)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Printf("current database %v", config.DB.Migrator().CurrentDatabase())
	log.Print("Migration success!")

}
