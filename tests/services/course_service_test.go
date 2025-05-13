package services_test

import (
	"errors"
	"testing"
	"time"
	"web/mocks/repos"
	"web/models"
	"web/schemas"
	"web/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCourseService_GetAllCourses(t *testing.T) {

	mockRepo := new(mocks.CourseRepositoryInterface)

	service := services.NewCourseService(mockRepo)

	expectedCourses := []schemas.CourseResponseWithChaptersCount{
		{
			ID:            1,
			Name:          "Test Course 1",
			Description:   "Test Description 1",
			ChaptersCount: 2,
		},
		{
			ID:            2,
			Name:          "Test Course 2",
			Description:   "Test Description 2",
			ChaptersCount: 0,
		},
	}

	mockRepo.On("GetAll").Return(expectedCourses, nil)

	courses, err := service.GetAllCourses()

	assert.NoError(t, err)
	assert.Equal(t, expectedCourses, courses)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetCourseByID(t *testing.T) {

	mockRepo := new(mocks.CourseRepositoryInterface)

	service := services.NewCourseService(mockRepo)

	testCases := []struct {
		name           string
		courseID       uint
		mockSetup      func()
		expectedError  error
		expectedCourse models.Course
	}{
		{
			name:     "Success",
			courseID: 1,
			mockSetup: func() {
				expectedCourse := models.Course{
					ID:          1,
					Name:        "Test Course",
					Description: "Test Description",
				}
				mockRepo.On("GetByID", uint(1)).Return(expectedCourse, nil)
			},
			expectedError: nil,
			expectedCourse: models.Course{
				ID:          1,
				Name:        "Test Course",
				Description: "Test Description",
			},
		},
		{
			name:     "Course Not Found",
			courseID: 999,
			mockSetup: func() {
				mockRepo.On("GetByID", uint(999)).Return(models.Course{}, errors.New("course not found"))
			},
			expectedError:  errors.New("course not found"),
			expectedCourse: models.Course{},
		},
		{
			name:     "Invalid ID",
			courseID: 0,
			mockSetup: func() {

			},
			expectedError:  errors.New("course ID is required"),
			expectedCourse: models.Course{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.mockSetup()

			course, err := service.GetCourseByID(tc.courseID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedCourse, course)
		})
	}

	mockRepo.AssertExpectations(t)
}

func TestCourseService_CreateCourse(t *testing.T) {

	mockRepo := new(mocks.CourseRepositoryInterface)

	service := services.NewCourseService(mockRepo)

	testCases := []struct {
		name           string
		courseRequest  schemas.CreateCourseRequest
		mockSetup      func()
		expectedError  error
		expectedCourse schemas.CourseResponse
	}{
		{
			name: "Success",
			courseRequest: schemas.CreateCourseRequest{
				Name:        "Test Course",
				Description: "Test Description",
			},
			mockSetup: func() {
				createdCourse := models.Course{
					ID:          1,
					Name:        "Test Course",
					Description: "Test Description",
					CreatedAt:   time.Now(),
				}
				mockRepo.On("Create", mock.MatchedBy(func(course models.Course) bool {
					return course.Name == "Test Course" && course.Description == "Test Description"
				})).Return(createdCourse, nil)
			},
			expectedError: nil,
			expectedCourse: schemas.CourseResponse{
				ID:          1,
				Name:        "Test Course",
				Description: "Test Description",
			},
		},
		{
			name: "Missing Name",
			courseRequest: schemas.CreateCourseRequest{
				Description: "Test Description",
			},
			mockSetup: func() {

			},
			expectedError:  errors.New("course name is required"),
			expectedCourse: schemas.CourseResponse{},
		},
		{
			name: "Missing Description",
			courseRequest: schemas.CreateCourseRequest{
				Name: "Test Course",
			},
			mockSetup: func() {

			},
			expectedError:  errors.New("course description is required"),
			expectedCourse: schemas.CourseResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.mockSetup()

			course, err := service.CreateCourse(tc.courseRequest)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCourse.ID, course.ID)
				assert.Equal(t, tc.expectedCourse.Name, course.Name)
				assert.Equal(t, tc.expectedCourse.Description, course.Description)
			}
		})
	}

	mockRepo.AssertExpectations(t)
}
