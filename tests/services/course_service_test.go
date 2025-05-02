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

type MockCourseRepository struct {
	mock.Mock
}

var _ repos.CourseRepositoryInterface = (*MockCourseRepository)(nil)

func (m *MockCourseRepository) GetAll() ([]models.Course, error) {
	args := m.Called()
	return args.Get(0).([]models.Course), args.Error(1)
}

func (m *MockCourseRepository) GetByID(id uint) (models.Course, error) {
	args := m.Called(id)
	return args.Get(0).(models.Course), args.Error(1)
}

func (m *MockCourseRepository) Create(course models.Course) (models.Course, error) {
	args := m.Called(course)
	return args.Get(0).(models.Course), args.Error(1)
}

func (m *MockCourseRepository) Update(course models.Course, courseRequest schemas.UpdateCourseRequest) (models.Course, error) {
	args := m.Called(course, courseRequest)
	return args.Get(0).(models.Course), args.Error(1)
}

func (m *MockCourseRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCourseService_GetAllCourses(t *testing.T) {

	mockRepo := new(MockCourseRepository)

	courses := []models.Course{
		{ID: 1, Name: "Course 1", Description: "Description 1"},
		{ID: 2, Name: "Course 2", Description: "Description 2"},
	}

	mockRepo.On("GetAll").Return(courses, nil)

	service := services.NewCourseService(mockRepo)

	result, err := service.GetAllCourses()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Course 1", result[0].Name)
	assert.Equal(t, uint(2), result[1].ID)
	assert.Equal(t, "Course 2", result[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetCourseByID(t *testing.T) {

	mockRepo := new(MockCourseRepository)

	course := models.Course{ID: 1, Name: "Course 1", Description: "Description 1"}

	t.Run("Course found", func(t *testing.T) {

		mockRepo.On("GetByID", uint(1)).Return(course, nil).Once()

		service := services.NewCourseService(mockRepo)

		result, err := service.GetCourseByID(1)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Course 1", result.Name)
		assert.Equal(t, "Description 1", result.Description)
	})

	t.Run("Course not found", func(t *testing.T) {

		mockRepo.On("GetByID", uint(999)).Return(models.Course{}, errors.New("course not found")).Once()

		service := services.NewCourseService(mockRepo)

		_, err := service.GetCourseByID(999)

		assert.Error(t, err)
		assert.Equal(t, "course not found", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestCourseService_CreateCourse(t *testing.T) {

	mockRepo := new(MockCourseRepository)

	courseRequest := schemas.CreateCourseRequest{
		Name:        "New Course",
		Description: "New Description",
	}

	createdCourse := models.Course{
		ID:          1,
		Name:        "New Course",
		Description: "New Description",
	}

	t.Run("Valid course creation", func(t *testing.T) {

		mockRepo.On("Create", mock.AnythingOfType("models.Course")).Return(createdCourse, nil).Once()

		service := services.NewCourseService(mockRepo)

		result, err := service.CreateCourse(courseRequest)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "New Course", result.Name)
		assert.Equal(t, "New Description", result.Description)
	})

	t.Run("Missing name", func(t *testing.T) {

		service := services.NewCourseService(mockRepo)

		_, err := service.CreateCourse(schemas.CreateCourseRequest{
			Description: "New Description",
		})

		assert.Error(t, err)
		assert.Equal(t, "course name is required", err.Error())
	})

	t.Run("Missing description", func(t *testing.T) {

		service := services.NewCourseService(mockRepo)

		_, err := service.CreateCourse(schemas.CreateCourseRequest{
			Name: "New Course",
		})

		assert.Error(t, err)
		assert.Equal(t, "course description is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestCourseService_UpdateCourse(t *testing.T) {

	mockRepo := new(MockCourseRepository)

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

	t.Run("Valid course update", func(t *testing.T) {

		mockRepo.On("Update", course, updateRequest).Return(updatedCourse, nil).Once()

		service := services.NewCourseService(mockRepo)

		result, err := service.UpdateCourse(course, updateRequest)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Updated Course", result.Name)
		assert.Equal(t, "Updated Description", result.Description)
	})

	t.Run("Missing ID", func(t *testing.T) {

		service := services.NewCourseService(mockRepo)

		_, err := service.UpdateCourse(models.Course{Name: "Course"}, updateRequest)

		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})

	t.Run("Missing name", func(t *testing.T) {

		service := services.NewCourseService(mockRepo)

		_, err := service.UpdateCourse(models.Course{ID: 1}, updateRequest)

		assert.Error(t, err)
		assert.Equal(t, "course name is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestCourseService_DeleteCourse(t *testing.T) {

	mockRepo := new(MockCourseRepository)

	t.Run("Valid course deletion", func(t *testing.T) {

		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		service := services.NewCourseService(mockRepo)

		err := service.DeleteCourse(1)

		assert.NoError(t, err)
	})

	t.Run("Course not found", func(t *testing.T) {

		mockRepo.On("Delete", uint(999)).Return(errors.New("course not found")).Once()

		service := services.NewCourseService(mockRepo)

		err := service.DeleteCourse(999)

		assert.Error(t, err)
		assert.Equal(t, "course not found", err.Error())
	})

	t.Run("Missing ID", func(t *testing.T) {

		service := services.NewCourseService(mockRepo)

		err := service.DeleteCourse(0)

		assert.Error(t, err)
		assert.Equal(t, "course ID is required", err.Error())
	})

	mockRepo.AssertExpectations(t)
}
