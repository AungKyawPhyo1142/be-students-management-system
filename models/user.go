package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)" json:"name"`
	Username string `gorm:"type:varchar(100)" json:"username"`
	Role     string `gorm:"type:varchar(100);default:'admin'" json:"role"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	UserName  string `json:"username"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetAllUsersResponse struct {
	Data []UserResponse `json:"data"`
}

func (u User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		UserName:  u.Username,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
	}
}

func (u User) GetAllUsersResponse(users []User) GetAllUsersResponse {
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToUserResponse())
	}
	return GetAllUsersResponse{
		Data: userResponses,
	}
}
