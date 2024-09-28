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
	ClassID        string    `gorm:"not null" json:"class_code"`
	EnrollmentDate time.Time `gorm:"not null" json:"enrollment_date"`
}

type ToStudentClassResponse struct {
	ID             uint      `json:"id"`
	StudentID      uint      `json:"student_id"`
	ClassID        string    `json:"class_code"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ToAllStudentClassResponse struct {
	Data []ToStudentClassResponse `json:"data"`
}

func (sc StudentClass) ToStudentClassResponse() ToStudentClassResponse {
	return ToStudentClassResponse{
		ID:             sc.ID,
		StudentID:      sc.StudentID,
		ClassID:        sc.ClassID,
		EnrollmentDate: sc.EnrollmentDate,
		CreatedAt:      sc.CreatedAt,
		UpdatedAt:      sc.UpdatedAt,
	}
}

func (sc StudentClass) ToAllStudentClassResponse(values []StudentClass) ToAllStudentClassResponse {
	var response []ToStudentClassResponse

	for _, value := range values {
		response = append(response, value.ToStudentClassResponse())
	}

	return ToAllStudentClassResponse{
		Data: response,
	}

}
