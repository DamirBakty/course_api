package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"web/models"
	"web/repos"
	"web/schemas"
)

type UserServiceInterface interface {
	RegisterUser(userDTO schemas.RegisterUserRequest) (schemas.UserResponse, error)
	ClaimUserUserFromToken(claims *KeycloakClaims) (schemas.UserResponse, error)
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
	user, err := s.repo.GetByUsername(userDTO.Username)
	if err == nil {
		userResponse := schemas.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
		}

		return userResponse, nil
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
	user = models.User{
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

func (s *UserService) ClaimUserUserFromToken(claims *KeycloakClaims) (schemas.UserResponse, error) {
	if claims.PreferredUsername == "" {
		return schemas.UserResponse{}, errors.New("username is required")
	}
	if claims.Email == "" {
		return schemas.UserResponse{}, errors.New("email is required")
	}
	user, err := s.repo.GetBySub(claims.Sub)
	if err == nil {
		userResponse := schemas.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Roles:     user.Roles,
			Sub:       user.Sub,
			CreatedAt: user.CreatedAt,
		}
		return userResponse, nil
	}

	// Extract roles from token
	roles := "ROLE_USER" // Default role
	if len(claims.RealmAccess.Roles) > 0 {
		roles = strings.Join(claims.RealmAccess.Roles, ",")
	}

	// Generate a random password since we're using token-based auth
	randomPassword := make([]byte, 32)
	_, err = rand.Read(randomPassword)
	if err != nil {
		return schemas.UserResponse{}, errors.New("failed to generate random password")
	}

	// Hash the random password
	hashedPassword, err := bcrypt.GenerateFromPassword(randomPassword, bcrypt.DefaultCost)
	if err != nil {
		return schemas.UserResponse{}, errors.New("failed to hash password")
	}

	// Create user
	user = models.User{
		Username: claims.PreferredUsername,
		Email:    claims.Email,
		Password: string(hashedPassword),
		Roles:    roles,
		Sub:      claims.Sub,
	}

	user, err = s.repo.Create(user)
	if err != nil {
		return schemas.UserResponse{}, err
	}
	fmt.Print(user.ID)

	userResponse := schemas.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
	}

	return userResponse, nil
}
