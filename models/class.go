package models

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	ClassName  string    `gorm:"type:varchar(100)" json:"className" validate:"required"`
	ClassCode  string    `gorm:"type:varchar(100)" json:"classCode" validate:"required"`
	Instructor string    `gorm:"type:varchar(100)" json:"instructor" validate:"required"`
	Semester   string    `gorm:"type:varchar(20)" json:"semester" validate:"required,oneof=Spring Summer Fall Winter"`
	Year       int       `json:"year" validate:"required,numeric"`
	Credits    int       `json:"credits" validate:"required,numeric"`
	Students   []Student `gorm:"many2many:student_classes;"`
}

type ClassResponse struct {
	ClassName  string `json:"className"`
	ClassCode  string `json:"classCode"`
	Instructor string `json:"instructor"`
	Semester   string `json:"semester"`
	Year       int    `json:"year"`
	Credits    int    `json:"credits"`
}

type ClassesResponse struct {
	Data []ClassResponse `json:"data"`
}

func (c Class) ToClassResponse() ClassResponse {
	return ClassResponse{
		ClassName:  c.ClassName,
		ClassCode:  c.ClassCode,
		Instructor: c.Instructor,
		Semester:   c.Semester,
		Year:       c.Year,
		Credits:    c.Credits,
	}
}

func (c Class) GetAllClassesResponse(classes []Class) ClassesResponse {
	var classesResponse []ClassResponse
	for _, class := range classes {
		classesResponse = append(classesResponse, class.ToClassResponse())
	}
	return ClassesResponse{
		Data: classesResponse,
	}
}
