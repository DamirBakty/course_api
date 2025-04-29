package services

import (
	"errors"
	"web/models"
	"web/repos"
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

func (s *CourseService) CreateCourse(course models.Course) (uint, error) {
	if course.Name == "" {
		return 0, errors.New("course name is required")
	}

	return s.repo.Create(course)
}

func (s *CourseService) UpdateCourse(course models.Course) error {
	if course.ID == 0 {
		return errors.New("course ID is required")
	}

	if course.Name == "" {
		return errors.New("course name is required")
	}

	return s.repo.Update(course)
}

func (s *CourseService) DeleteCourse(id uint) error {
	if id == 0 {
		return errors.New("course ID is required")
	}

	return s.repo.Delete(id)
}
