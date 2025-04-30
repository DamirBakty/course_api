package repos

import (
	"web/models"
	"web/schemas"
)

// CourseRepositoryInterface defines the interface for course repository operations
type CourseRepositoryInterface interface {
	GetAll() ([]models.Course, error)
	GetByID(id uint) (models.Course, error)
	Create(course models.Course) (models.Course, error)
	Update(course models.Course, courseRequest schemas.UpdateCourseRequest) (models.Course, error)
	Delete(id uint) error
}

// ChapterRepositoryInterface defines the interface for chapter repository operations
type ChapterRepositoryInterface interface {
	GetAll() ([]models.Chapter, error)
	GetByID(id uint) (models.Chapter, error)
	GetByCourseID(courseID uint) ([]models.Chapter, error)
	Create(chapter models.Chapter) (uint, error)
	Update(chapter models.Chapter) error
	Delete(id uint) error
}

// LessonRepositoryInterface defines the interface for lesson repository operations
type LessonRepositoryInterface interface {
	GetAll() ([]models.Lesson, error)
	GetByID(id uint) (models.Lesson, error)
	GetByChapterID(chapterID uint) ([]models.Lesson, error)
	Create(lesson models.Lesson) (uint, error)
	Update(lesson models.Lesson) error
	Delete(id uint) error
}
