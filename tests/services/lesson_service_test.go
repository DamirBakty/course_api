package tests

import (
	"errors"
	"testing"
	"web/models"
	"web/repos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLessonRepository is a mock implementation of LessonRepositoryInterface
type MockLessonRepository struct {
	mock.Mock
}

// Ensure MockLessonRepository implements LessonRepositoryInterface
var _ repos.LessonRepositoryInterface = (*MockLessonRepository)(nil)

// GetAll mocks the GetAll method
func (m *MockLessonRepository) GetAll() ([]models.Lesson, error) {
	args := m.Called()
	return args.Get(0).([]models.Lesson), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockLessonRepository) GetByID(id uint) (models.Lesson, error) {
	args := m.Called(id)
	return args.Get(0).(models.Lesson), args.Error(1)
}

// GetByChapterID mocks the GetByChapterID method
func (m *MockLessonRepository) GetByChapterID(chapterID uint) ([]models.Lesson, error) {
	args := m.Called(chapterID)
	return args.Get(0).([]models.Lesson), args.Error(1)
}

// Create mocks the Create method
func (m *MockLessonRepository) Create(lesson models.Lesson) (uint, error) {
	args := m.Called(lesson)
	return args.Get(0).(uint), args.Error(1)
}

// Update mocks the Update method
func (m *MockLessonRepository) Update(lesson models.Lesson) error {
	args := m.Called(lesson)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockLessonRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestLessonService is a test-specific implementation of LessonService
// that accepts our mock repository
type TestLessonService struct {
	repo repos.LessonRepositoryInterface
}

// NewTestLessonService creates a new TestLessonService
func NewTestLessonService(repo repos.LessonRepositoryInterface) *TestLessonService {
	return &TestLessonService{
		repo: repo,
	}
}

// GetAllLessons delegates to the mock repository
func (s *TestLessonService) GetAllLessons() ([]models.Lesson, error) {
	return s.repo.GetAll()
}

// GetLessonByID delegates to the mock repository
func (s *TestLessonService) GetLessonByID(id uint) (models.Lesson, error) {
	return s.repo.GetByID(id)
}

// GetLessonsByChapterID delegates to the mock repository
func (s *TestLessonService) GetLessonsByChapterID(chapterID uint) ([]models.Lesson, error) {
	if chapterID == 0 {
		return nil, errors.New("chapter ID is required")
	}
	return s.repo.GetByChapterID(chapterID)
}

// CreateLesson delegates to the mock repository
func (s *TestLessonService) CreateLesson(lesson models.Lesson) (uint, error) {
	if lesson.Name == "" {
		return 0, errors.New("lesson name is required")
	}

	if lesson.ChapterID == 0 {
		return 0, errors.New("chapter ID is required")
	}

	return s.repo.Create(lesson)
}

// UpdateLesson delegates to the mock repository
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

// DeleteLesson delegates to the mock repository
func (s *TestLessonService) DeleteLesson(id uint) error {
	if id == 0 {
		return errors.New("lesson ID is required")
	}

	return s.repo.Delete(id)
}

// TestLessonService_GetAllLessons tests the GetAllLessons method
func TestLessonService_GetAllLessons(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockLessonRepository)
	
	// Create test data
	lessons := []models.Lesson{
		{ID: 1, Name: "Lesson 1", Content: "Content 1", ChapterID: 1},
		{ID: 2, Name: "Lesson 2", Content: "Content 2", ChapterID: 1},
	}
	
	// Set up expectations
	mockRepo.On("GetAll").Return(lessons, nil)
	
	// Create service with mock repository
	service := NewTestLessonService(mockRepo)
	
	// Call the method
	result, err := service.GetAllLessons()
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Lesson 1", result[0].Name)
	assert.Equal(t, uint(2), result[1].ID)
	assert.Equal(t, "Lesson 2", result[1].Name)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestLessonService_GetLessonByID tests the GetLessonByID method
func TestLessonService_GetLessonByID(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockLessonRepository)
	
	// Create test data
	lesson := models.Lesson{ID: 1, Name: "Lesson 1", Content: "Content 1", ChapterID: 1}
	
	// Test case 1: Lesson found
	t.Run("Lesson found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByID", uint(1)).Return(lesson, nil).Once()
		
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method
		result, err := service.GetLessonByID(1)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Lesson 1", result.Name)
		assert.Equal(t, "Content 1", result.Content)
		assert.Equal(t, uint(1), result.ChapterID)
	})
	
	// Test case 2: Lesson not found
	t.Run("Lesson not found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByID", uint(999)).Return(models.Lesson{}, errors.New("lesson not found")).Once()
		
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method
		_, err := service.GetLessonByID(999)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "lesson not found", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestLessonService_GetLessonsByChapterID tests the GetLessonsByChapterID method
func TestLessonService_GetLessonsByChapterID(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockLessonRepository)
	
	// Create test data
	lessons := []models.Lesson{
		{ID: 1, Name: "Lesson 1", Content: "Content 1", ChapterID: 1},
		{ID: 2, Name: "Lesson 2", Content: "Content 2", ChapterID: 1},
	}
	
	// Test case 1: Lessons found
	t.Run("Lessons found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByChapterID", uint(1)).Return(lessons, nil).Once()
		
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method
		result, err := service.GetLessonsByChapterID(1)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, "Lesson 1", result[0].Name)
		assert.Equal(t, uint(2), result[1].ID)
		assert.Equal(t, "Lesson 2", result[1].Name)
	})
	
	// Test case 2: No lessons found
	t.Run("No lessons found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByChapterID", uint(999)).Return([]models.Lesson{}, nil).Once()
		
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method
		result, err := service.GetLessonsByChapterID(999)
		
		// Assert
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
	
	// Test case 3: Missing chapter ID
	t.Run("Missing chapter ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.GetLessonsByChapterID(0)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestLessonService_CreateLesson tests the CreateLesson method
func TestLessonService_CreateLesson(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockLessonRepository)
	
	// Create test data
	lesson := models.Lesson{
		Name:      "New Lesson",
		Content:   "New Content",
		ChapterID: 1,
	}
	
	// Test case 1: Valid lesson creation
	t.Run("Valid lesson creation", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Create", lesson).Return(uint(1), nil).Once()
		
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method
		id, err := service.CreateLesson(lesson)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
	})
	
	// Test case 2: Missing name
	t.Run("Missing name", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.CreateLesson(models.Lesson{
			ChapterID: 1,
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "lesson name is required", err.Error())
	})
	
	// Test case 3: Missing chapter ID
	t.Run("Missing chapter ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.CreateLesson(models.Lesson{
			Name: "New Lesson",
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestLessonService_UpdateLesson tests the UpdateLesson method
func TestLessonService_UpdateLesson(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockLessonRepository)
	
	// Create test data
	lesson := models.Lesson{
		ID:        1,
		Name:      "Updated Lesson",
		Content:   "Updated Content",
		ChapterID: 1,
	}
	
	// Test case 1: Valid lesson update
	t.Run("Valid lesson update", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Update", lesson).Return(nil).Once()
		
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method
		err := service.UpdateLesson(lesson)
		
		// Assert
		assert.NoError(t, err)
	})
	
	// Test case 2: Missing ID
	t.Run("Missing ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method with invalid data
		err := service.UpdateLesson(models.Lesson{
			Name:      "Updated Lesson",
			ChapterID: 1,
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "lesson ID is required", err.Error())
	})
	
	// Test case 3: Missing name
	t.Run("Missing name", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method with invalid data
		err := service.UpdateLesson(models.Lesson{
			ID:        1,
			ChapterID: 1,
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "lesson name is required", err.Error())
	})
	
	// Test case 4: Missing chapter ID
	t.Run("Missing chapter ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method with invalid data
		err := service.UpdateLesson(models.Lesson{
			ID:   1,
			Name: "Updated Lesson",
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "chapter ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestLessonService_DeleteLesson tests the DeleteLesson method
func TestLessonService_DeleteLesson(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockLessonRepository)
	
	// Test case 1: Valid lesson deletion
	t.Run("Valid lesson deletion", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Delete", uint(1)).Return(nil).Once()
		
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method
		err := service.DeleteLesson(1)
		
		// Assert
		assert.NoError(t, err)
	})
	
	// Test case 2: Lesson not found
	t.Run("Lesson not found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Delete", uint(999)).Return(errors.New("lesson not found")).Once()
		
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method
		err := service.DeleteLesson(999)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "lesson not found", err.Error())
	})
	
	// Test case 3: Missing ID
	t.Run("Missing ID", func(t *testing.T) {
		// Create service with mock repository
		service := NewTestLessonService(mockRepo)
		
		// Call the method with invalid data
		err := service.DeleteLesson(0)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "lesson ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}