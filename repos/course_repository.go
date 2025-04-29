package repos

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"web/models"
	"web/schemas"
)

type CourseRepository struct {
	DB *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{
		DB: db,
	}
}

func (r *CourseRepository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	result := r.DB.Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}

	return courses, nil
}

func (r *CourseRepository) GetByID(id uint) (models.Course, error) {
	var course models.Course
	result := r.DB.First(&course, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return course, errors.New("course not found")
		}
		return course, result.Error
	}

	return course, nil
}

func (r *CourseRepository) Create(course models.Course) (models.Course, error) {
	result := r.DB.Create(&course)
	if result.Error != nil {
		return models.Course{}, result.Error
	}

	return course, nil
}

func (r *CourseRepository) Update(course models.Course, courseRequest schemas.UpdateCourseRequest) (models.Course, error) {
	result := r.DB.Model(&course).Updates(models.Course{
		Name:        courseRequest.Name,
		Description: courseRequest.Description,
		UpdatedAt:   time.Now(),
	})

	if result.Error != nil {
		return models.Course{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Course{}, errors.New("course not found or no changes made")
	}

	return course, nil
}

func (r *CourseRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Course{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("course not found")
	}

	return nil
}
