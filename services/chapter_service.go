package services

import (
	"errors"
	"web/models"
	"web/repos"
)

type ChapterService struct {
	repo *repos.ChapterRepository
}

func NewChapterService(repo *repos.ChapterRepository) *ChapterService {
	return &ChapterService{
		repo: repo,
	}
}

func (s *ChapterService) GetAllChapters() ([]models.Chapter, error) {
	return s.repo.GetAll()
}

func (s *ChapterService) GetChaptersByCourseID(courseID uint) ([]models.Chapter, error) {
	if courseID == 0 {
		return nil, errors.New("course ID is required")
	}
	return s.repo.GetByCourseID(courseID)
}

func (s *ChapterService) GetChapterByID(id uint) (models.Chapter, error) {
	return s.repo.GetByID(id)
}

func (s *ChapterService) CreateChapter(chapter models.Chapter) (uint, error) {
	if chapter.Name == "" {
		return 0, errors.New("chapter name is required")
	}

	if chapter.CourseID == 0 {
		return 0, errors.New("course ID is required")
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
