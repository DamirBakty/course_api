package tests

import (
	"errors"
	"testing"
	"web/models"
	"web/repos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChapterRepository struct {
	mock.Mock
}

var _ repos.ChapterRepositoryInterface = (*MockChapterRepository)(nil)

func (m *MockChapterRepository) GetAll() ([]models.Chapter, error) {
	args := m.Called()
	return args.Get(0).([]models.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByID(id uint) (models.Chapter, error) {
	args := m.Called(id)
	return args.Get(0).(models.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByCourseID(courseID uint) ([]models.Chapter, error) {
	args := m.Called(courseID)
	return args.Get(0).([]models.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Create(chapter models.Chapter) (uint, error) {
	args := m.Called(chapter)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockChapterRepository) Update(chapter models.Chapter) error {
	args := m.Called(chapter)
	return args.Error(0)
}

func (m *MockChapterRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type TestChapterService struct {
	repo repos.ChapterRepositoryInterface
}

func NewTestChapterService(repo repos.ChapterRepositoryInterface) *TestChapterService {
	return &TestChapterService{
		repo: repo,
	}
}

func (s *TestChapterService) GetAllChapters() ([]models.Chapter, error) {
	return s.repo.GetAll()
}

func (s *TestChapterService) GetChapterByID(id uint) (models.Chapter, error) {
	return s.repo.GetByID(id)
}

func (s *TestChapterService) GetChaptersByCourseID(courseID uint) ([]models.Chapter, error) {
	if courseID == 0 {
		return nil, errors.New("course ID is required")
	}
	return s.repo.GetByCourseID(courseID)
}

func (s *TestChapterService) CreateChapter(chapter models.Chapter) (uint, error) {
	if chapter.Name == "" {
		return 0, errors.New("chapter name is required")
	}

	if chapter.CourseID == 0 {
		return 0, errors.New("course ID is required")
	}

	return s.repo.Create(chapter)
}

func (s *TestChapterService) UpdateChapter(chapter models.Chapter) error {
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

func (s *TestChapterService) DeleteChapter(id uint) error {
	if id == 0 {
		return errors.New("chapter ID is required")
	}

	return s.repo.Delete(id)
}

func TestChapterService_GetAllChapters(t *testing.T) {

	mockRepo := new(MockChapterRepository)

	chapters := []models.Chapter{
		{ID: 1, Name: "Chapter 1", Description: "Description 1", CourseID: 1},
		{ID: 2, Name: "Chapter 2", Description: "Description 2", CourseID: 1},
	}

	mockRepo.On("GetAll").Return(chapters, nil)

	service := NewTestChapterService(mockRepo)

	result, err := service.GetAllChapters()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Chapter 1", result[0].Name)
	assert.Equal(t, uint(2), result[1].ID)
	assert.Equal(t, "Chapter 2", result[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestChapterService_GetChapterByID(t *testing.T) {

	mockRepo := new(MockChapterRepository)

	chapter := models.Chapter{ID: 1, Name: "Chapter 1", Description: "Description 1", CourseID: 1}

	t.Run("Chapter found", func(t *testing.T) {

		mockRepo.On("GetByID", uint(1)).Return(chapter, nil).Once()

		service := NewTestChapterService(mockRepo)

		result, err := service.GetChapterByID(1)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Chapter 1", result.Name)
		assert.Equal(t, "Description 1", result.Description)
		assert.Equal(t, uint(1), result.CourseID)
	})

	t.Run("Chapter not found", func(t *testing.T) {

		mockRepo.On("GetByID", uint(999)).Return(models.Chapter{}, errors.New("chapter not found")).Once()

		service := NewTestChapterService(mockRepo)

		_, err := service.GetChapterByID(999)

		assert.Error(t, err)
		assert.Equal(t, "chapter not found", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestChapterService_GetChaptersByCourseID(t *testing.T) {

	mockRepo := new(MockChapterRepository)

	chapters := []models.Chapter{
		{ID: 1, Name: "Chapter 1", Description: "Description 1", CourseID: 1},
		{ID: 2, Name: "Chapter 2", Description: "Description 2", CourseID: 1},
	}

	t.Run("Chapters found", func(t *testing.T) {

		mockRepo.On("GetByCourseID", uint(1)).Return(chapters, nil).Once()

		service := NewTestChapterService(mockRepo)

		result, err := service.GetChaptersByCourseID(1)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, "Chapter 1", result[0].Name)
		assert.Equal(t, uint(2), result[1].ID)
		assert.Equal(t, "Chapter 2", result[1].Name)
	})

	t.Run("No chapters found", func(t *testing.T) {

		mockRepo.On("GetByCourseID", uint(999)).Return([]models.Chapter{}, nil).Once()

		service := NewTestChapterService(mockRepo)

		result, err := service.GetChaptersByCourseID(999)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Missing course ID", func(t *testing.T) {

		service := NewTestChapterService(mockRepo)

		_, err := service.GetChaptersByCourseID(0)

		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestChapterService_CreateChapter(t *testing.T) {

	mockRepo := new(MockChapterRepository)

	chapter := models.Chapter{
		Name:        "New Chapter",
		Description: "New Description",
		CourseID:    1,
	}

	t.Run("Valid chapter creation", func(t *testing.T) {

		mockRepo.On("Create", chapter).Return(uint(1), nil).Once()

		service := NewTestChapterService(mockRepo)

		id, err := service.CreateChapter(chapter)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
	})

	t.Run("Missing name", func(t *testing.T) {

		service := NewTestChapterService(mockRepo)

		_, err := service.CreateChapter(models.Chapter{
			CourseID: 1,
		})

		assert.Error(t, err)
		assert.Equal(t, "chapter name is required", err.Error())
	})

	t.Run("Missing course ID", func(t *testing.T) {

		service := NewTestChapterService(mockRepo)

		_, err := service.CreateChapter(models.Chapter{
			Name: "New Chapter",
		})

		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestChapterService_UpdateChapter(t *testing.T) {

	mockRepo := new(MockChapterRepository)

	chapter := models.Chapter{
		ID:          1,
		Name:        "Updated Chapter",
		Description: "Updated Description",
		CourseID:    1,
	}

	t.Run("Valid chapter update", func(t *testing.T) {

		mockRepo.On("Update", chapter).Return(nil).Once()

		service := NewTestChapterService(mockRepo)

		err := service.UpdateChapter(chapter)

		assert.NoError(t, err)
	})

	t.Run("Missing ID", func(t *testing.T) {

		service := NewTestChapterService(mockRepo)

		err := service.UpdateChapter(models.Chapter{
			Name:     "Updated Chapter",
			CourseID: 1,
		})

		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})

	t.Run("Missing name", func(t *testing.T) {

		service := NewTestChapterService(mockRepo)

		err := service.UpdateChapter(models.Chapter{
			ID:       1,
			CourseID: 1,
		})

		assert.Error(t, err)
		assert.Equal(t, "chapter name is required", err.Error())
	})

	t.Run("Missing course ID", func(t *testing.T) {

		service := NewTestChapterService(mockRepo)

		err := service.UpdateChapter(models.Chapter{
			ID:   1,
			Name: "Updated Chapter",
		})

		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestChapterService_DeleteChapter(t *testing.T) {

	mockRepo := new(MockChapterRepository)

	t.Run("Valid chapter deletion", func(t *testing.T) {

		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		service := NewTestChapterService(mockRepo)

		err := service.DeleteChapter(1)

		assert.NoError(t, err)
	})

	t.Run("Chapter not found", func(t *testing.T) {

		mockRepo.On("Delete", uint(999)).Return(errors.New("chapter not found")).Once()

		service := NewTestChapterService(mockRepo)

		err := service.DeleteChapter(999)

		assert.Error(t, err)
		assert.Equal(t, "chapter not found", err.Error())
	})

	t.Run("Missing ID", func(t *testing.T) {

		service := NewTestChapterService(mockRepo)

		err := service.DeleteChapter(0)

		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}
