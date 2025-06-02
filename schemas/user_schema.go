package schemas

import "time"

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
	Roles    string `json:"roles" binding:"required" example:"user,admin"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJfT3B2QmJxS0VfdU5NbV..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJfT3B2QmJxS0VfdU5NbV..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" example:"300"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJfT3B2QmJxS0VfdU5NbV..."`
}

type UserResponse struct {
	ID        uint      `json:"id,omitempty" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Roles     string    `json:"roles" example:"user,admin"`
	Sub       string    `json:"sub" example:"1234567890"`
	CreatedAt time.Time `json:"created_at,omitempty" example:"2020-01-01T12:00:00Z"`
}

type AdminCreateUserRequest struct {
	Username string   `json:"username" binding:"required" example:"johndoe"`
	Email    string   `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string   `json:"password" binding:"required,min=8" example:"password123"`
	Roles    []string `json:"roles" binding:"required" example:"user,admin"`
}

type UserInfoResponse struct {
	Username string   `json:"username" example:"johndoe"`
	Email    string   `json:"email" example:"john.doe@example.com"`
	Roles    []string `json:"roles" binding:"required" example:"user,admin"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"oldpassword123"`
	NewPassword     string `json:"new_password" binding:"required,min=8" example:"newpassword123"`
}
