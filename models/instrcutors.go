package models

import "gorm.io/gorm"

type Instructor struct {
	gorm.Model
	Image     *string `gorm:"size:255" json:"image"`
	FirstName string  `gorm:"type:varchar(100)" json:"firstName" validate:"required"`
	LastName  string  `gorm:"type:varchar(100)" json:"lastName" validate:"required"`
	Email     string  `gorm:"type:varchar(100)" json:"email" validate:"required"`
	Phone     string  `gorm:"type:varchar(100)" json:"phone" validate:"required"`
	Classes   []Class `gorm:"many2many:instrcutor_classes;"`
}

type InstructorResponse struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Image     string  `json:"image"`
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Classes   []Class `json:"classes"`
}

type GetAllInstructorsResponse struct {
	Data []InstructorResponse `json:"data"`
}

func (i Instructor) ToInstructorResponse() InstructorResponse {
	image := ""
	if i.Image != nil {
		image = *i.Image
	}

	return InstructorResponse{
		ID:        i.ID,
		FirstName: i.FirstName,
		LastName:  i.LastName,
		Image:     image,
		Phone:     i.Phone,
		Email:     i.Email,
		CreatedAt: i.CreatedAt.String(),
		UpdatedAt: i.UpdatedAt.String(),
	}
}

func (i Instructor) GetAllInstructorsResponse(instructors []Instructor) GetAllInstructorsResponse {
	var instructorResponse []InstructorResponse

	for _, instructor := range instructors {
		instructorResponse = append(instructorResponse, instructor.ToInstructorResponse())
	}
	return GetAllInstructorsResponse{
		Data: instructorResponse,
	}
}
