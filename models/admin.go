package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)" json:"name" validate:"required"`
	Username string `gorm:"type:varchar(100)" json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `gorm:"type:varchar(100);default:'admin'" json:"role"`
}

type AdminResponse struct {
	ID        uint   `json:"id"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetAllUsersResponse struct {
	Data []AdminResponse `json:"data"`
}

func (u Admin) ToAdminResponse() AdminResponse {
	return AdminResponse{
		ID:        u.ID,
		Name:      u.Name,
		UserName:  u.Username,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
	}
}

func (a Admin) GetAllUsersResponse(admins []Admin) GetAllUsersResponse {
	var AdminResponses []AdminResponse
	for _, admin := range admins {
		AdminResponses = append(AdminResponses, admin.ToAdminResponse())
	}
	return GetAllUsersResponse{
		Data: AdminResponses,
	}
}
