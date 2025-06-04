package services

import (
	"errors"
	"time"
	"web/models"
	"web/repos"
	"web/schemas"
)

type ChapterServiceInterface interface {
	GetChapterByID(id, courseId uint) (models.Chapter, error)
	GetChaptersByCourseID(courseID uint) ([]schemas.ChapterResponseWithLessonsCount, error)
	CreateChapter(chapterRequest schemas.ChapterRequest, courseID uint) (uint, error)
	UpdateChapter(chapter models.Chapter) error
	DeleteChapter(id uint) error
}

var _ ChapterServiceInterface = (*ChapterService)(nil)

type ChapterService struct {
	repo       *repos.ChapterRepository
	courseRepo *repos.CourseRepository
}

func NewChapterService(repo *repos.ChapterRepository, courseRepo *repos.CourseRepository) *ChapterService {
	return &ChapterService{
		repo:       repo,
		courseRepo: courseRepo,
	}
}

func (s *ChapterService) GetChaptersByCourseID(courseID uint) ([]schemas.ChapterResponseWithLessonsCount, error) {
	if courseID == 0 {
		return nil, errors.New("course ID is required")
	}
	return s.repo.GetByCourseID(courseID)
}

func (s *ChapterService) GetChapterByID(id, courseId uint) (models.Chapter, error) {
	return s.repo.GetByID(id, courseId)
}

func (s *ChapterService) GetChapterByIDWithLessonsCount(id, courseID uint) (schemas.ChapterResponseWithLessonsCount, error) {
	return s.repo.GetByIDWithLessonsCount(id, courseID)
}

func (s *ChapterService) CreateChapter(chapterRequest schemas.ChapterRequest, courseID uint) (uint, error) {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return 0, err
	}
	chapter := models.Chapter{
		Name:        chapterRequest.Name,
		Description: chapterRequest.Description,
		Order:       chapterRequest.Order,
		CourseID:    course.ID,
		CreatedBy:   chapterRequest.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if chapterRequest.Name == "" {
		return 0, errors.New("chapter name is required")
	}
	if chapterRequest.Description == "" {
		return 0, errors.New("chapter description is required")
	}
	if chapterRequest.Order == 0 {
		return 0, errors.New("chapter order is required")
	}

	return s.repo.Create(chapter)
}

func (s *ChapterService) UpdateChapter(chapter models.Chapter) error {
	if chapter.ID == 0 {
		return errors.New("chapter ID is required")
	}

	if chapter.Name == "" {
		return errors.New("chapter name is required")
	}

	if chapter.CourseID == 0 {
		return errors.New("course ID is required")
	}

	return s.repo.Update(chapter)
}

func (s *ChapterService) DeleteChapter(id uint) error {
	if id == 0 {
		return errors.New("chapter ID is required")
	}

	return s.repo.Delete(id)
}
