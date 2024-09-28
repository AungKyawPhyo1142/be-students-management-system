package models

import (
	"time"

	"gorm.io/gorm"
)

// this is the join-table of
// instructor-class table
// M to M relationship
type InstructorClass struct {
	gorm.Model
	InstructorID uint   `gorm:"not null" json:"instructor_id"`
	ClassID      string `gorm:"not null" json:"class_code"`
}

type ToInstructorClassResponse struct {
	ID           uint      `json:"id"`
	InstructorID uint      `gorm:"not null" json:"instructor_id"`
	ClassID      string    `json:"class_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ToAllInstructorClassResponse struct {
	Data []ToInstructorClassResponse `json:"data"`
}

func (ic InstructorClass) ToStudentClassResponse() ToInstructorClassResponse {
	return ToInstructorClassResponse{
		ID:             ic.ID,
		InstructorID:   ic.InstructorID,
		ClassID:        ic.ClassID,
		CreatedAt:      ic.CreatedAt,
		UpdatedAt:      ic.UpdatedAt,
	}
}

func (sc InstructorClass) ToAllStudentClassResponse(values []InstructorClass) ToAllInstructorClassResponse {
	var response []ToInstructorClassResponse

	for _, value := range values {
		response = append(response, value.ToStudentClassResponse())
	}

	return ToAllInstructorClassResponse{
		Data: response,
	}

}
