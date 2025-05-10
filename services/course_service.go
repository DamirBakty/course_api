package services

import (
	"errors"
	"web/models"
	"web/repos"
	"web/schemas"
)

type CourseServiceInterface interface {
	GetAllCourses() ([]schemas.CourseResponseWithChaptersCount, error)
	GetCourseByID(id uint) (models.Course, error)
	GetCourseByIDWithChapterCount(id uint) (schemas.CourseResponseWithChaptersCount, error)
	CreateCourse(courseDTO schemas.CreateCourseRequest) (schemas.CourseResponse, error)
	UpdateCourse(course models.Course, courseRequest schemas.UpdateCourseRequest) (schemas.CourseResponse, error)
	DeleteCourse(id uint) error
}

// Ensure CourseService implements CourseServiceInterface
var _ CourseServiceInterface = (*CourseService)(nil)

type CourseService struct {
	repo repos.CourseRepositoryInterface
}

func (s *CourseService) GetCourseByIDWithChapterCount(id uint) (schemas.CourseResponseWithChaptersCount, error) {
	return s.repo.GetByIDWithChaptersCount(id)
}

func NewCourseService(repo repos.CourseRepositoryInterface) *CourseService {
	return &CourseService{
		repo: repo,
	}
}

func (s *CourseService) GetCourseByID(id uint) (models.Course, error) {
	if id == 0 {
		return models.Course{}, errors.New("course ID is required")
	}

	return s.repo.GetByID(id)
}

func (s *CourseService) CreateCourse(courseRequest schemas.CreateCourseRequest) (schemas.CourseResponse, error) {
	if courseRequest.Name == "" {
		return schemas.CourseResponse{}, errors.New("course name is required")
	}
	if courseRequest.Description == "" {
		return schemas.CourseResponse{}, errors.New("course description is required")
	}

	// Convert DTO to model
	course := models.Course{
		Name:        courseRequest.Name,
		Description: courseRequest.Description,
	}
	course, err := s.repo.Create(course)
	if err != nil {
		return schemas.CourseResponse{}, err
	}
	courseResponse := schemas.CourseResponse{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		CreatedAt:   course.CreatedAt,
	}
	return courseResponse, nil
}

func (s *CourseService) UpdateCourse(course models.Course, courseRequest schemas.UpdateCourseRequest) (schemas.CourseResponse, error) {
	if course.ID == 0 {
		return schemas.CourseResponse{}, errors.New("course ID is required")
	}

	if course.Name == "" {
		return schemas.CourseResponse{}, errors.New("course name is required")
	}
	course, err := s.repo.Update(course, courseRequest)
	if err != nil {
		return schemas.CourseResponse{}, err
	}
	courseResponse := schemas.CourseResponse{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		CreatedAt:   course.CreatedAt,
	}
	return courseResponse, nil
}

func (s *CourseService) DeleteCourse(id uint) error {
	if id == 0 {
		return errors.New("course ID is required")
	}

	return s.repo.Delete(id)
}

func (s *CourseService) GetAllCourses() ([]schemas.CourseResponseWithChaptersCount, error) {
	return s.repo.GetAll()
}
