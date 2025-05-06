package repos

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"web/models"
	"web/schemas"
)

var _ CourseRepositoryInterface = (*CourseRepository)(nil)

type CourseRepository struct {
	DB *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{
		DB: db,
	}
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

func (r *CourseRepository) GetAll() ([]schemas.CourseResponseWithChaptersCount, error) {
	var courseResponses []schemas.CourseResponseWithChaptersCount

	subQuery := r.DB.Model(&models.Chapter{}).
		Select("course_id, count(*) as chapters_count").
		Group("course_id")

	err := r.DB.Model(&models.Course{}).
		Select("course.id, course.name, course.description, course.created_at, COALESCE(chapters_count, 0) as chapters_count").
		Joins("LEFT JOIN (?) AS chapter_counts ON course.id = chapter_counts.course_id", subQuery).
		Scan(&courseResponses).Error

	if err != nil {
		return nil, err
	}

	return courseResponses, nil
}

func (r *CourseRepository) GetByID(id uint) (models.Course, error) {
	var course models.Course
	err := r.DB.First(&course, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return course, errors.New("course not found")
		}
		return course, err
	}

	return course, nil
}

func (r *CourseRepository) GetByIDWithChaptersCount(id uint) (schemas.CourseResponseWithChaptersCount, error) {
	var courseResponse schemas.CourseResponseWithChaptersCount

	subQuery := r.DB.Model(&models.Chapter{}).
		Select("course_id, count(*) as chapters_count").
		Where("course_id = ?", id).
		Group("course_id")

	err := r.DB.Model(&models.Course{}).
		Select("course.id, course.name, course.description, course.created_at, COALESCE(chapters_count, 0) as chapters_count").
		Joins("LEFT JOIN (?) AS chapter_counts ON course.id = chapter_counts.course_id", subQuery).
		Where("course.id = ?", id).
		Scan(&courseResponse).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return courseResponse, errors.New("course not found")
		}
		return courseResponse, err
	}

	if courseResponse.ID == 0 {
		return courseResponse, errors.New("course not found")
	}

	return courseResponse, nil
}
