package tests

import (
	"errors"
	"testing"
	"web/models"
	"web/repos"
	"web/schemas"
	"web/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCourseRepository is a mock implementation of CourseRepositoryInterface
type MockCourseRepository struct {
	mock.Mock
}

// Ensure MockCourseRepository implements CourseRepositoryInterface
var _ repos.CourseRepositoryInterface = (*MockCourseRepository)(nil)

// GetAll mocks the GetAll method
func (m *MockCourseRepository) GetAll() ([]models.Course, error) {
	args := m.Called()
	return args.Get(0).([]models.Course), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockCourseRepository) GetByID(id uint) (models.Course, error) {
	args := m.Called(id)
	return args.Get(0).(models.Course), args.Error(1)
}

// Create mocks the Create method
func (m *MockCourseRepository) Create(course models.Course) (models.Course, error) {
	args := m.Called(course)
	return args.Get(0).(models.Course), args.Error(1)
}

// Update mocks the Update method
func (m *MockCourseRepository) Update(course models.Course, courseRequest schemas.UpdateCourseRequest) (models.Course, error) {
	args := m.Called(course, courseRequest)
	return args.Get(0).(models.Course), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockCourseRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestCourseService_GetAllCourses tests the GetAllCourses method
func TestCourseService_GetAllCourses(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockCourseRepository)
	
	// Create test data
	courses := []models.Course{
		{ID: 1, Name: "Course 1", Description: "Description 1"},
		{ID: 2, Name: "Course 2", Description: "Description 2"},
	}
	
	// Set up expectations
	mockRepo.On("GetAll").Return(courses, nil)
	
	// Create service with mock repository
	service := services.NewCourseService(mockRepo)
	
	// Call the method
	result, err := service.GetAllCourses()
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Course 1", result[0].Name)
	assert.Equal(t, uint(2), result[1].ID)
	assert.Equal(t, "Course 2", result[1].Name)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestCourseService_GetCourseByID tests the GetCourseByID method
func TestCourseService_GetCourseByID(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockCourseRepository)
	
	// Create test data
	course := models.Course{ID: 1, Name: "Course 1", Description: "Description 1"}
	
	// Test case 1: Course found
	t.Run("Course found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByID", uint(1)).Return(course, nil).Once()
		
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method
		result, err := service.GetCourseByID(1)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Course 1", result.Name)
		assert.Equal(t, "Description 1", result.Description)
	})
	
	// Test case 2: Course not found
	t.Run("Course not found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("GetByID", uint(999)).Return(models.Course{}, errors.New("course not found")).Once()
		
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method
		_, err := service.GetCourseByID(999)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course not found", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestCourseService_CreateCourse tests the CreateCourse method
func TestCourseService_CreateCourse(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockCourseRepository)
	
	// Create test data
	courseRequest := schemas.CreateCourseRequest{
		Name:        "New Course",
		Description: "New Description",
	}
	
	createdCourse := models.Course{
		ID:          1,
		Name:        "New Course",
		Description: "New Description",
	}
	
	// Test case 1: Valid course creation
	t.Run("Valid course creation", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Create", mock.AnythingOfType("models.Course")).Return(createdCourse, nil).Once()
		
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method
		result, err := service.CreateCourse(courseRequest)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "New Course", result.Name)
		assert.Equal(t, "New Description", result.Description)
	})
	
	// Test case 2: Missing name
	t.Run("Missing name", func(t *testing.T) {
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.CreateCourse(schemas.CreateCourseRequest{
			Description: "New Description",
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course name is required", err.Error())
	})
	
	// Test case 3: Missing description
	t.Run("Missing description", func(t *testing.T) {
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.CreateCourse(schemas.CreateCourseRequest{
			Name: "New Course",
		})
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course description is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestCourseService_UpdateCourse tests the UpdateCourse method
func TestCourseService_UpdateCourse(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockCourseRepository)
	
	// Create test data
	course := models.Course{
		ID:          1,
		Name:        "Course 1",
		Description: "Description 1",
	}
	
	updateRequest := schemas.UpdateCourseRequest{
		Name:        "Updated Course",
		Description: "Updated Description",
	}
	
	updatedCourse := models.Course{
		ID:          1,
		Name:        "Updated Course",
		Description: "Updated Description",
	}
	
	// Test case 1: Valid course update
	t.Run("Valid course update", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Update", course, updateRequest).Return(updatedCourse, nil).Once()
		
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method
		result, err := service.UpdateCourse(course, updateRequest)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Updated Course", result.Name)
		assert.Equal(t, "Updated Description", result.Description)
	})
	
	// Test case 2: Missing ID
	t.Run("Missing ID", func(t *testing.T) {
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.UpdateCourse(models.Course{Name: "Course"}, updateRequest)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})
	
	// Test case 3: Missing name
	t.Run("Missing name", func(t *testing.T) {
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method with invalid data
		_, err := service.UpdateCourse(models.Course{ID: 1}, updateRequest)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course name is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestCourseService_DeleteCourse tests the DeleteCourse method
func TestCourseService_DeleteCourse(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockCourseRepository)
	
	// Test case 1: Valid course deletion
	t.Run("Valid course deletion", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Delete", uint(1)).Return(nil).Once()
		
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method
		err := service.DeleteCourse(1)
		
		// Assert
		assert.NoError(t, err)
	})
	
	// Test case 2: Course not found
	t.Run("Course not found", func(t *testing.T) {
		// Set up expectations
		mockRepo.On("Delete", uint(999)).Return(errors.New("course not found")).Once()
		
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method
		err := service.DeleteCourse(999)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course not found", err.Error())
	})
	
	// Test case 3: Missing ID
	t.Run("Missing ID", func(t *testing.T) {
		// Create service with mock repository
		service := services.NewCourseService(mockRepo)
		
		// Call the method with invalid data
		err := service.DeleteCourse(0)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}