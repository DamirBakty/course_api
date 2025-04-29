package services

import (
	"errors"
	"web/models"
	"web/repos"
	"web/schemas"
)

type CourseService struct {
	repo *repos.CourseRepository
}

func NewCourseService(repo *repos.CourseRepository) *CourseService {
	return &CourseService{
		repo: repo,
	}
}

func (s *CourseService) GetAllCourses() ([]models.Course, error) {
	return s.repo.GetAll()
}

func (s *CourseService) GetCourseByID(id uint) (models.Course, error) {
	return s.repo.GetByID(id)
}

func (s *CourseService) CreateCourse(courseDTO schemas.CreateCourseRequest) (models.Course, error) {
	if courseDTO.Name == "" {
		return models.Course{}, errors.New("course name is required")
	}
	if courseDTO.Description == "" {
		return models.Course{}, errors.New("course description is required")
	}

	// Convert DTO to model
	course := models.Course{
		Name:        courseDTO.Name,
		Description: courseDTO.Description,
	}

	return s.repo.Create(course)
}

func (s *CourseService) UpdateCourse(course models.Course, courseRequest schemas.UpdateCourseRequest) (models.Course, error) {
	if course.ID == 0 {
		return models.Course{}, errors.New("course ID is required")
	}

	if course.Name == "" {
		return models.Course{}, errors.New("course name is required")
	}

	return s.repo.Update(course, courseRequest)
}

func (s *CourseService) DeleteCourse(id uint) error {
	if id == 0 {
		return errors.New("course ID is required")
	}

	return s.repo.Delete(id)
}
