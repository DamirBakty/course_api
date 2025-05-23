package repos

import (
	"errors"
	"gorm.io/gorm"
	"web/models"
)

type UserRepositoryInterface interface {
	Create(user models.User) (models.User, error)
	GetByUsername(username string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetBySub(sub string) (models.User, error)
}

var _ UserRepositoryInterface = (*UserRepository)(nil)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(user models.User) (models.User, error) {
	result := r.DB.Create(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}
func (r *UserRepository) GetBySub(sub string) (models.User, error) {
	var user models.User
	err := r.DB.Where("sub = ?", sub).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil

}
func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}
