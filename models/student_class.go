package models

import (
	"time"

	"gorm.io/gorm"
)

// this is the join-table of
// student-class table
// M to M relationship
type StudentClass struct {
	gorm.Model
	StudentID      uint      `gorm:"not null" json:"student_id"`
	ClassID        uint      `gorm:"not null" json:"class_id"`
	EnrollmentDate time.Time `gorm:"not null" json:"enrollment_date"`
}
