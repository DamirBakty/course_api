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
	UpdateUser(userID uint, userDTO schemas.UpdateUserRequest) (schemas.UserResponse, error)
	UpdatePassword(userID uint, passwordDTO schemas.UpdatePasswordRequest) error
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

	err := authService.RegisterUserInKeycloak(userDTO.Username, userDTO.Email, userDTO.Password, userDTO.Roles)
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

func (s *UserService) UpdateUser(userID uint, userDTO schemas.UpdateUserRequest) (schemas.UserResponse, error) {
	if userDTO.Username == "" {
		return schemas.UserResponse{}, errors.New("username is required")
	}
	if userDTO.Email == "" {
		return schemas.UserResponse{}, errors.New("email is required")
	}

	existingUser, err := s.repo.GetByUsername(userDTO.Username)
	if err == nil && existingUser.ID != userID {
		return schemas.UserResponse{}, errors.New("username already exists")
	} else if err != nil && err.Error() != "user not found" {
		return schemas.UserResponse{}, err
	}

	existingUser, err = s.repo.GetByEmail(userDTO.Email)
	if err == nil && existingUser.ID != userID {
		return schemas.UserResponse{}, errors.New("email already exists")
	} else if err != nil && err.Error() != "user not found" {
		return schemas.UserResponse{}, err
	}

	var user models.User
	user.ID = userID
	user.Username = userDTO.Username
	user.Email = userDTO.Email

	updatedUser, err := s.repo.Update(user)
	if err != nil {
		return schemas.UserResponse{}, err
	}

	userResponse := schemas.UserResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Roles:     updatedUser.Roles,
		Sub:       updatedUser.Sub,
		CreatedAt: updatedUser.CreatedAt,
	}

	return userResponse, nil
}

func (s *UserService) UpdatePassword(userID uint, passwordDTO schemas.UpdatePasswordRequest) error {
	if passwordDTO.CurrentPassword == "" {
		return errors.New("current password is required")
	}
	if passwordDTO.NewPassword == "" {
		return errors.New("new password is required")
	}
	if len(passwordDTO.NewPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	var user models.User
	user.ID = userID

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordDTO.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	err = s.repo.UpdatePassword(userID, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}
