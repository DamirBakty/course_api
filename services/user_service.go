package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"web/models"
	"web/repos"
	"web/schemas"
)

type UserServiceInterface interface {
	RegisterUser(userDTO schemas.RegisterUserRequest) (schemas.UserResponse, error)
}

var _ UserServiceInterface = (*UserService)(nil)

type UserService struct {
	repo repos.UserRepositoryInterface
}

func NewUserService(repo repos.UserRepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(userDTO schemas.RegisterUserRequest) (schemas.UserResponse, error) {
	if userDTO.Username == "" {
		return schemas.UserResponse{}, errors.New("username is required")
	}
	if userDTO.Email == "" {
		return schemas.UserResponse{}, errors.New("email is required")
	}
	if userDTO.Password == "" {
		return schemas.UserResponse{}, errors.New("password is required")
	}
	if userDTO.Roles == "" {
		return schemas.UserResponse{}, errors.New("roles are required")
	}

	// Check if username already exists
	_, err := s.repo.GetByUsername(userDTO.Username)
	if err == nil {
		return schemas.UserResponse{}, errors.New("username already exists")
	} else if err.Error() != "user not found" {
		return schemas.UserResponse{}, err
	}

	// Check if email already exists
	_, err = s.repo.GetByEmail(userDTO.Email)
	if err == nil {
		return schemas.UserResponse{}, errors.New("email already exists")
	} else if err.Error() != "user not found" {
		return schemas.UserResponse{}, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return schemas.UserResponse{}, errors.New("failed to hash password")
	}

	// Create user
	user := models.User{
		Username: userDTO.Username,
		Email:    userDTO.Email,
		Password: string(hashedPassword),
		Roles:    userDTO.Roles,
	}

	user, err = s.repo.Create(user)
	if err != nil {
		return schemas.UserResponse{}, err
	}

	userResponse := schemas.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
	}

	return userResponse, nil
}
