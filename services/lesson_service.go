package services

import (
	"errors"
	"time"
	"web/models"
	"web/repos"
	"web/schemas"
)

type LessonServiceInterface interface {
	GetLessonByID(courseID, chapterID, id uint) (schemas.LessonResponse, error)
	GetLessonsByChapterID(courseID, chapterID uint) ([]schemas.LessonResponse, error)
	CreateLesson(lessonRequest schemas.LessonRequest, courseId, chapterId uint) (uint, error)
	UpdateLesson(courseID, chapterID, id uint, lessonRequest schemas.LessonRequest) error
	DeleteLesson(courseID, chapterID, id uint) error
}

type LessonService struct {
	repo        *repos.LessonRepository
	chapterRepo *repos.ChapterRepository
	courseRepo  *repos.CourseRepository
}

func NewLessonService(repo *repos.LessonRepository, chapterRepo *repos.ChapterRepository, courseRepo *repos.CourseRepository) *LessonService {
	return &LessonService{
		repo:        repo,
		chapterRepo: chapterRepo,
		courseRepo:  courseRepo,
	}
}

func (s *LessonService) GetLessonsByChapterID(courseID, chapterID uint) ([]schemas.LessonResponse, error) {
	if chapterID == 0 {
		return nil, errors.New("chapter ID is required")
	}
	if courseID == 0 {
		return nil, errors.New("course ID is required")
	}
	return s.repo.GetByChapterID(courseID, chapterID)
}

func (s *LessonService) GetLessonByID(courseID, chapterID, id uint) (schemas.LessonResponse, error) {
	lesson, err := s.repo.GetByID(courseID, chapterID, id)
	if err != nil {
		return schemas.LessonResponse{}, err
	}
	lessonResponse := schemas.LessonResponse{
		ID:          lesson.ID,
		Name:        lesson.Name,
		Description: lesson.Description,
		Content:     lesson.Content,
		Order:       lesson.Order,
		CreatedAt:   lesson.CreatedAt.Format(time.RFC3339),
	}
	return lessonResponse, nil
}

func (s *LessonService) CreateLesson(lessonRequest schemas.LessonRequest, courseId, chapterId uint) (uint, error) {
	course, err := s.courseRepo.GetByID(courseId)
	if err != nil {
		return 0, err
	}
	chapter, err := s.chapterRepo.GetByID(course.ID, chapterId)
	if err != nil {
		return 0, err
	}

	lesson := models.Lesson{
		Name:        lessonRequest.Name,
		Description: lessonRequest.Description,
		Content:     lessonRequest.Content,
		Order:       lessonRequest.Order,
		ChapterID:   chapter.ID,
		CreatedBy:   lessonRequest.CreatedBy,
	}
	if lesson.Name == "" {
		return 0, errors.New("lesson name is required")
	}
	if lesson.Description == "" {
		return 0, errors.New("lesson description is required")
	}
	if lesson.Content == "" {
		return 0, errors.New("lesson content is required")
	}
	if lesson.Order == 0 {
		return 0, errors.New("lesson order is required")
	}

	if lesson.ChapterID == 0 {
		return 0, errors.New("chapter ID is required")
	}

	return s.repo.Create(lesson)
}

func (s *LessonService) UpdateLesson(courseID, chapterID, id uint, lessonRequest schemas.LessonRequest) error {
	lesson, err := s.repo.GetByID(courseID, chapterID, id)
	if err != nil {
		return err
	}
	lesson.Name = lessonRequest.Name
	lesson.Description = lessonRequest.Description
	lesson.Content = lessonRequest.Content
	lesson.Order = lessonRequest.Order
	lesson.UpdatedAt = time.Now()

	return s.repo.Update(lesson)
}

func (s *LessonService) DeleteLesson(courseID, chapterID, id uint) error {
	lesson, err := s.repo.GetByID(courseID, chapterID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(lesson.ID)
}
