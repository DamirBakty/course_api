package schemas

import "time"

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
	Roles    string `json:"roles" binding:"required" example:"ROLE_USER,ROLE_ADMIN"`
}

type UserResponse struct {
	ID        uint      `json:"id,omitempty" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Roles     string    `json:"roles" example:"ROLE_USER,ROLE_ADMIN"`
	Sub       string    `json:"sub" example:"1234567890"`
	CreatedAt time.Time `json:"created_at,omitempty" example:"2020-01-01T12:00:00Z"`
}
