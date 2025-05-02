package tests

import (
	"errors"
	"testing"
	"web/models"
	"web/repos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLessonRepository struct {
	mock.Mock
}

var _ repos.LessonRepositoryInterface = (*MockLessonRepository)(nil)

func (m *MockLessonRepository) GetAll() ([]models.Lesson, error) {
	args := m.Called()
	return args.Get(0).([]models.Lesson), args.Error(1)
}

func (m *MockLessonRepository) GetByID(id uint) (models.Lesson, error) {
	args := m.Called(id)
	return args.Get(0).(models.Lesson), args.Error(1)
}

func (m *MockLessonRepository) GetByChapterID(chapterID uint) ([]models.Lesson, error) {
	args := m.Called(chapterID)
	return args.Get(0).([]models.Lesson), args.Error(1)
}

func (m *MockLessonRepository) Create(lesson models.Lesson) (uint, error) {
	args := m.Called(lesson)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockLessonRepository) Update(lesson models.Lesson) error {
	args := m.Called(lesson)
	return args.Error(0)
}

func (m *MockLessonRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type TestLessonService struct {
	repo repos.LessonRepositoryInterface
}

func NewTestLessonService(repo repos.LessonRepositoryInterface) *TestLessonService {
	return &TestLessonService{
		repo: repo,
	}
}

func (s *TestLessonService) GetAllLessons() ([]models.Lesson, error) {
	return s.repo.GetAll()
}

func (s *TestLessonService) GetLessonByID(id uint) (models.Lesson, error) {
	return s.repo.GetByID(id)
}

func (s *TestLessonService) GetLessonsByChapterID(chapterID uint) ([]models.Lesson, error) {
	if chapterID == 0 {
		return nil, errors.New("chapter ID is required")
	}
	return s.repo.GetByChapterID(chapterID)
}

func (s *TestLessonService) CreateLesson(lesson models.Lesson) (uint, error) {
	if lesson.Name == "" {
		return 0, errors.New("lesson name is required")
	}

	if lesson.ChapterID == 0 {
		return 0, errors.New("chapter ID is required")
	}

	return s.repo.Create(lesson)
}

func (s *TestLessonService) UpdateLesson(lesson models.Lesson) error {
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

func (s *TestLessonService) DeleteLesson(id uint) error {
	if id == 0 {
		return errors.New("lesson ID is required")
	}

	return s.repo.Delete(id)
}

func TestLessonService_GetAllLessons(t *testing.T) {

	mockRepo := new(MockLessonRepository)

	lessons := []models.Lesson{
		{ID: 1, Name: "Lesson 1", Content: "Content 1", ChapterID: 1},
		{ID: 2, Name: "Lesson 2", Content: "Content 2", ChapterID: 1},
	}

	mockRepo.On("GetAll").Return(lessons, nil)

	service := NewTestLessonService(mockRepo)

	result, err := service.GetAllLessons()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Lesson 1", result[0].Name)
	assert.Equal(t, uint(2), result[1].ID)
	assert.Equal(t, "Lesson 2", result[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestLessonService_GetLessonByID(t *testing.T) {

	mockRepo := new(MockLessonRepository)

	lesson := models.Lesson{ID: 1, Name: "Lesson 1", Content: "Content 1", ChapterID: 1}

	t.Run("Lesson found", func(t *testing.T) {

		mockRepo.On("GetByID", uint(1)).Return(lesson, nil).Once()

		service := NewTestLessonService(mockRepo)

		result, err := service.GetLessonByID(1)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Lesson 1", result.Name)
		assert.Equal(t, "Content 1", result.Content)
		assert.Equal(t, uint(1), result.ChapterID)
	})

	t.Run("Lesson not found", func(t *testing.T) {

		mockRepo.On("GetByID", uint(999)).Return(models.Lesson{}, errors.New("lesson not found")).Once()

		service := NewTestLessonService(mockRepo)

		_, err := service.GetLessonByID(999)

		assert.Error(t, err)
		assert.Equal(t, "lesson not found", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestLessonService_GetLessonsByChapterID(t *testing.T) {

	mockRepo := new(MockLessonRepository)

	lessons := []models.Lesson{
		{ID: 1, Name: "Lesson 1", Content: "Content 1", ChapterID: 1},
		{ID: 2, Name: "Lesson 2", Content: "Content 2", ChapterID: 1},
	}

	t.Run("Lessons found", func(t *testing.T) {

		mockRepo.On("GetByChapterID", uint(1)).Return(lessons, nil).Once()

		service := NewTestLessonService(mockRepo)

		result, err := service.GetLessonsByChapterID(1)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, "Lesson 1", result[0].Name)
		assert.Equal(t, uint(2), result[1].ID)
		assert.Equal(t, "Lesson 2", result[1].Name)
	})

	t.Run("No lessons found", func(t *testing.T) {

		mockRepo.On("GetByChapterID", uint(999)).Return([]models.Lesson{}, nil).Once()

		service := NewTestLessonService(mockRepo)

		result, err := service.GetLessonsByChapterID(999)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Missing chapter ID", func(t *testing.T) {

		service := NewTestLessonService(mockRepo)

		_, err := service.GetLessonsByChapterID(0)

		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestLessonService_CreateLesson(t *testing.T) {

	mockRepo := new(MockLessonRepository)

	lesson := models.Lesson{
		Name:      "New Lesson",
		Content:   "New Content",
		ChapterID: 1,
	}

	t.Run("Valid lesson creation", func(t *testing.T) {

		mockRepo.On("Create", lesson).Return(uint(1), nil).Once()

		service := NewTestLessonService(mockRepo)

		id, err := service.CreateLesson(lesson)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
	})

	t.Run("Missing name", func(t *testing.T) {

		service := NewTestLessonService(mockRepo)

		_, err := service.CreateLesson(models.Lesson{
			ChapterID: 1,
		})

		assert.Error(t, err)
		assert.Equal(t, "lesson name is required", err.Error())
	})

	t.Run("Missing chapter ID", func(t *testing.T) {

		service := NewTestLessonService(mockRepo)

		_, err := service.CreateLesson(models.Lesson{
			Name: "New Lesson",
		})

		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestLessonService_UpdateLesson(t *testing.T) {

	mockRepo := new(MockLessonRepository)

	lesson := models.Lesson{
		ID:        1,
		Name:      "Updated Lesson",
		Content:   "Updated Content",
		ChapterID: 1,
	}

	t.Run("Valid lesson update", func(t *testing.T) {

		mockRepo.On("Update", lesson).Return(nil).Once()

		service := NewTestLessonService(mockRepo)

		err := service.UpdateLesson(lesson)

		assert.NoError(t, err)
	})

	t.Run("Missing ID", func(t *testing.T) {

		service := NewTestLessonService(mockRepo)

		err := service.UpdateLesson(models.Lesson{
			Name:      "Updated Lesson",
			ChapterID: 1,
		})

		assert.Error(t, err)
		assert.Equal(t, "lesson ID is required", err.Error())
	})

	t.Run("Missing name", func(t *testing.T) {

		service := NewTestLessonService(mockRepo)

		err := service.UpdateLesson(models.Lesson{
			ID:        1,
			ChapterID: 1,
		})

		assert.Error(t, err)
		assert.Equal(t, "lesson name is required", err.Error())
	})

	t.Run("Missing chapter ID", func(t *testing.T) {

		service := NewTestLessonService(mockRepo)

		err := service.UpdateLesson(models.Lesson{
			ID:   1,
			Name: "Updated Lesson",
		})

		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestLessonService_DeleteLesson(t *testing.T) {

	mockRepo := new(MockLessonRepository)

	t.Run("Valid lesson deletion", func(t *testing.T) {

		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		service := NewTestLessonService(mockRepo)

		err := service.DeleteLesson(1)

		assert.NoError(t, err)
	})

	t.Run("Lesson not found", func(t *testing.T) {

		mockRepo.On("Delete", uint(999)).Return(errors.New("lesson not found")).Once()

		service := NewTestLessonService(mockRepo)

		err := service.DeleteLesson(999)

		assert.Error(t, err)
		assert.Equal(t, "lesson not found", err.Error())
	})

	t.Run("Missing ID", func(t *testing.T) {

		service := NewTestLessonService(mockRepo)

		err := service.DeleteLesson(0)

		assert.Error(t, err)
		assert.Equal(t, "lesson ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}
