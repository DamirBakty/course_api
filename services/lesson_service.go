package services

import (
	"errors"
	"web/models"
	"web/repos"
)

type LessonService struct {
	repo *repos.LessonRepository
}

func NewLessonService(repo *repos.LessonRepository) *LessonService {
	return &LessonService{
		repo: repo,
	}
}

func (s *LessonService) GetLessonsByChapterID(chapterID, courseID uint) ([]models.Lesson, error) {
	if chapterID == 0 {
		return nil, errors.New("chapter ID is required")
	}
	return s.repo.GetByChapterID(chapterID)
}

func (s *LessonService) GetLessonByID(id uint) (models.Lesson, error) {
	return s.repo.GetByID(id)
}

func (s *LessonService) CreateLesson(lesson models.Lesson) (uint, error) {
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

func (s *LessonService) UpdateLesson(lesson models.Lesson) error {
	if lesson.ID == 0 {
		return errors.New("lesson ID is required")
	}

	if lesson.Name == "" {
		return errors.New("lesson name is required")
	}

	if lesson.ChapterID == 0 {
		return errors.New("chapter ID is required")
	}

	return s.repo.Update(lesson)
}

func (s *LessonService) DeleteLesson(id uint) error {
	if id == 0 {
		return errors.New("lesson ID is required")
	}

	return s.repo.Delete(id)
}
