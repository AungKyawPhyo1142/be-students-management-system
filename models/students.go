package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	FirstName string  `gorm:"type:varchar(100)" json:"firstName" validate:"required"`
	LastName  string  `gorm:"type:varchar(100)" json:"lastName" validate:"required"`
	DOB       string  `gorm:"type:varchar(100)" json:"dob" validate:"required"`
	Email     string  `gorm:"type:varchar(100)" json:"email" validate:"required"`
	Phone     string  `gorm:"type:varchar(100)" json:"phone" validate:"required"`
	Classes   []Class `gorm:"many2many:student_classes;"` // sepecify that this has RS with classes table
}

type StudentResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	DOB       string `json:"date_of_birth"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetAllStudentsResponse struct {
	Data []StudentResponse `json:"data"`
}

func (s Student) ToStudentResponse() StudentResponse {
	return StudentResponse{
		ID:        s.ID,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		DOB:       s.DOB,
		Phone:     s.Phone,
		Email:     s.Email,
		CreatedAt: s.CreatedAt.String(),
		UpdatedAt: s.UpdatedAt.String(),
	}
}

func (s Student) GetAllStudentsResponse(students []Student) GetAllStudentsResponse {
	var studentResponse []StudentResponse
	for _, student := range students {
		studentResponse = append(studentResponse, student.ToStudentResponse())
	}
	return GetAllStudentsResponse{
		Data: studentResponse,
	}
}
