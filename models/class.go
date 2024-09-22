package models

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	ClassName  string    `gorm:"type:varchar(100)" json:"className" validate:"required"`
	ClassCode  string    `gorm:"type:varchar(100)" json:"classCode" validate:"required"`
	Instructor string    `gorm:"type:varchar(100)" json:"instructor" validate:"required"`
	Semester   string    `gorm:"type:varchar(100)" json:"semester" validate:"required"`
	Credits    string    `gorm:"type:varchar(100)" json:"credits" validate:"required"`
	Students   []Student `gorm:"many2many:student_classes;"`
}
