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
	AdminCreateUser(userDTO schemas.AdminCreateUserRequest, authService *AuthService) (schemas.UserInfoResponse, error)
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

	_, err = s.repo.GetByEmail(userDTO.Email)
	if err == nil {
		return schemas.UserResponse{}, errors.New("email already exists")
	} else if err.Error() != "user not found" {
		return schemas.UserResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return schemas.UserResponse{}, errors.New("failed to hash password")
	}

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

	roles := "ROLE_USER"
	if len(claims.RealmAccess.Roles) > 0 {
		roles = strings.Join(claims.RealmAccess.Roles, ",")
	}

	randomPassword := make([]byte, 32)
	_, err = rand.Read(randomPassword)
	if err != nil {
		return schemas.UserResponse{}, errors.New("failed to generate random password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(randomPassword, bcrypt.DefaultCost)
	if err != nil {
		return schemas.UserResponse{}, errors.New("failed to hash password")
	}

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

func (s *UserService) AdminCreateUser(userDTO schemas.AdminCreateUserRequest, authService *AuthService) (schemas.UserInfoResponse, error) {
	if userDTO.Username == "" {
		return schemas.UserInfoResponse{}, errors.New("username is required")
	}
	if userDTO.Email == "" {
		return schemas.UserInfoResponse{}, errors.New("email is required")
	}
	if userDTO.Password == "" {
		return schemas.UserInfoResponse{}, errors.New("password is required")
	}
	if len(userDTO.Roles) == 0 {
		return schemas.UserInfoResponse{}, errors.New("roles are required")
	}

	_, err := s.repo.GetByUsername(userDTO.Username)
	if err == nil {
		return schemas.UserInfoResponse{}, errors.New("username already exists")
	} else if err.Error() != "user not found" {
		return schemas.UserInfoResponse{}, err
	}

	_, err = s.repo.GetByEmail(userDTO.Email)
	if err == nil {
		return schemas.UserInfoResponse{}, errors.New("email already exists")
	} else if err.Error() != "user not found" {
		return schemas.UserInfoResponse{}, err
	}

	err = authService.RegisterUserInKeycloak(userDTO.Username, userDTO.Email, userDTO.Password, userDTO.Roles)
	if err != nil {
		return schemas.UserInfoResponse{}, fmt.Errorf("failed to register user in Keycloak: %w", err)
	}

	userResponse := schemas.UserInfoResponse{
		Username: userDTO.Username,
		Email:    userDTO.Email,
		Roles:    userDTO.Roles,
	}

	return userResponse, nil
}
