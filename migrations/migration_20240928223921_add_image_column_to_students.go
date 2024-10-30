package migrations

import (
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
	"gorm.io/gorm"
)

func Migrate_20240928223921_add_image_column_to_students(tx *gorm.DB) error {
	return tx.Migrator().AddColumn(&models.Student{}, "Image")

}

func Rollback_20240928223921_add_image_column_to_students(tx *gorm.DB) error {
	return tx.Migrator().DropColumn(&models.Student{}, "Image")
}
