package tests

import (
	"errors"
	"testing"
	"web/models"
	"web/repos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChapterRepository is a mock implementation of ChapterRepositoryInterface
type MockChapterRepository struct {
	mock.Mock
}

// Ensure MockChapterRepository implements ChapterRepositoryInterface
var _ repos.ChapterRepositoryInterface = (*MockChapterRepository)(nil)

// GetAll mocks the GetAll method
func (m *MockChapterRepository) GetAll() ([]models.Chapter, error) {
	args := m.Called()
	return args.Get(0).([]models.Chapter), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockChapterRepository) GetByID(id uint) (models.Chapter, error) {
	args := m.Called(id)
	return args.Get(0).(models.Chapter), args.Error(1)
}

// GetByCourseID mocks the GetByCourseID method
func (m *MockChapterRepository) GetByCourseID(courseID uint) ([]models.Chapter, error) {
	args := m.Called(courseID)
	return args.Get(0).([]models.Chapter), args.Error(1)
}

// Create mocks the Create method
func (m *MockChapterRepository) Create(chapter models.Chapter) (uint, error) {
	args := m.Called(chapter)
	return args.Get(0).(uint), args.Error(1)
}

// Update mocks the Update method
func (m *MockChapterRepository) Update(chapter models.Chapter) error {
	args := m.Called(chapter)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockChapterRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestChapterService is a test-specific implementation of ChapterService
// that accepts our mock repository
type TestChapterService struct {
	repo repos.ChapterRepositoryInterface
}

// NewTestChapterService creates a new TestChapterService
func NewTestChapterService(repo repos.ChapterRepositoryInterface) *TestChapterService {
	return &TestChapterService{
		repo: repo,
	}
}

// GetAllChapters delegates to the mock repository
func (s *TestChapterService) GetAllChapters() ([]models.Chapter, error) {
	return s.repo.GetAll()
}

// GetChapterByID delegates to the mock repository
func (s *TestChapterService) GetChapterByID(id uint) (models.Chapter, error) {
	return s.repo.GetByID(id)
}

// GetChaptersByCourseID delegates to the mock repository
func (s *TestChapterService) GetChaptersByCourseID(courseID uint) ([]models.Chapter, error) {
	if courseID == 0 {
		return nil, errors.New("course ID is required")
	}
	return s.repo.GetByCourseID(courseID)
}

// CreateChapter delegates to the mock repository
func (s *TestChapterService) CreateChapter(chapter models.Chapter) (uint, error) {
	if chapter.Name == "" {
		return 0, errors.New("chapter name is required")
	}

	if chapter.CourseID == 0 {
		return 0, errors.New("course ID is required")
	}

	return s.repo.Create(chapter)
}

// UpdateChapter delegates to the mock repository
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

// DeleteChapter delegates to the mock repository
func (s *TestChapterService) DeleteChapter(id uint) error {
	if id == 0 {
		return errors.New("chapter ID is required")
	}

	return s.repo.Delete(id)
}

// TestChapterService_GetAllChapters tests the GetAllChapters method
func TestChapterService_GetAllChapters(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockChapterRepository)
	
	// Create test data
	chapters := []models.Chapter{
		{ID: 1, Name: "Chapter 1", Description: "Description 1", CourseID: 1},
		{ID: 2, Name: "Chapter 2", Description: "Description 2", CourseID: 1},
	}
	
	// Set up expectations
	mockRepo.On("GetAll").Return(chapters, nil)
	
	// Create service with mock repository
	service := NewTestChapterService(mockRepo)
	
	// Call the method
	result, err := service.GetAllChapters()
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Chapter 1", result[0].Name)
	assert.Equal(t, uint(2), result[1].ID)
	assert.Equal(t, "Chapter 2", result[1].Name)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestChapterService_GetChapterByID tests the GetChapterByID method
func TestChapterService_GetChapterByID(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockChapterRepository)
	
	// Create test data
	chapter := models.Chapter{ID: 1, Name: "Chapter 1", Description: "Description 1", CourseID: 1}
	
	// Test case 1: Chapter found
	t.Run("Chapter found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByID", uint(1)).Return(chapter, nil).Once()
		
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method
		result, err := service.GetChapterByID(1)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Chapter 1", result.Name)
		assert.Equal(t, "Description 1", result.Description)
		assert.Equal(t, uint(1), result.CourseID)
	})
	
	// Test case 2: Chapter not found
	t.Run("Chapter not found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByID", uint(999)).Return(models.Chapter{}, errors.New("chapter not found")).Once()
		
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method
		_, err := service.GetChapterByID(999)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter not found", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestChapterService_GetChaptersByCourseID tests the GetChaptersByCourseID method
func TestChapterService_GetChaptersByCourseID(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockChapterRepository)
	
	// Create test data
	chapters := []models.Chapter{
		{ID: 1, Name: "Chapter 1", Description: "Description 1", CourseID: 1},
		{ID: 2, Name: "Chapter 2", Description: "Description 2", CourseID: 1},
	}
	
	// Test case 1: Chapters found
	t.Run("Chapters found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByCourseID", uint(1)).Return(chapters, nil).Once()
		
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method
		result, err := service.GetChaptersByCourseID(1)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, "Chapter 1", result[0].Name)
		assert.Equal(t, uint(2), result[1].ID)
		assert.Equal(t, "Chapter 2", result[1].Name)
	})
	
	// Test case 2: No chapters found
	t.Run("No chapters found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByCourseID", uint(999)).Return([]models.Chapter{}, nil).Once()
		
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method
		result, err := service.GetChaptersByCourseID(999)
		
		// Assert
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
	
	// Test case 3: Missing course ID
	t.Run("Missing course ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.GetChaptersByCourseID(0)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestChapterService_CreateChapter tests the CreateChapter method
func TestChapterService_CreateChapter(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockChapterRepository)
	
	// Create test data
	chapter := models.Chapter{
		Name:        "New Chapter",
		Description: "New Description",
		CourseID:    1,
	}
	
	// Test case 1: Valid chapter creation
	t.Run("Valid chapter creation", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Create", chapter).Return(uint(1), nil).Once()
		
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method
		id, err := service.CreateChapter(chapter)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
	})
	
	// Test case 2: Missing name
	t.Run("Missing name", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.CreateChapter(models.Chapter{
			CourseID: 1,
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter name is required", err.Error())
	})
	
	// Test case 3: Missing course ID
	t.Run("Missing course ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.CreateChapter(models.Chapter{
			Name: "New Chapter",
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestChapterService_UpdateChapter tests the UpdateChapter method
func TestChapterService_UpdateChapter(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockChapterRepository)
	
	// Create test data
	chapter := models.Chapter{
		ID:          1,
		Name:        "Updated Chapter",
		Description: "Updated Description",
		CourseID:    1,
	}
	
	// Test case 1: Valid chapter update
	t.Run("Valid chapter update", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Update", chapter).Return(nil).Once()
		
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method
		err := service.UpdateChapter(chapter)
		
		// Assert
		assert.NoError(t, err)
	})
	
	// Test case 2: Missing ID
	t.Run("Missing ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method with invalid data
		err := service.UpdateChapter(models.Chapter{
			Name:     "Updated Chapter",
			CourseID: 1,
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})
	
	// Test case 3: Missing name
	t.Run("Missing name", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method with invalid data
		err := service.UpdateChapter(models.Chapter{
			ID:       1,
			CourseID: 1,
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter name is required", err.Error())
	})
	
	// Test case 4: Missing course ID
	t.Run("Missing course ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method with invalid data
		err := service.UpdateChapter(models.Chapter{
			ID:   1,
			Name: "Updated Chapter",
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestChapterService_DeleteChapter tests the DeleteChapter method
func TestChapterService_DeleteChapter(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockChapterRepository)
	
	// Test case 1: Valid chapter deletion
	t.Run("Valid chapter deletion", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Delete", uint(1)).Return(nil).Once()
		
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method
		err := service.DeleteChapter(1)
		
		// Assert
		assert.NoError(t, err)
	})
	
	// Test case 2: Chapter not found
	t.Run("Chapter not found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Delete", uint(999)).Return(errors.New("chapter not found")).Once()
		
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method
		err := service.DeleteChapter(999)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter not found", err.Error())
	})
	
	// Test case 3: Missing ID
	t.Run("Missing ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestChapterService(mockRepo)
		
		// Call the method with invalid data
		err := service.DeleteChapter(0)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}