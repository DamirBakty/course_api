package repos

import (
	"web/models"
	"web/schemas"
)

type CourseRepositoryInterface interface {
	GetAll() ([]schemas.CourseResponseWithChaptersCount, error)
	GetByID(id uint) (models.Course, error)
	Create(course models.Course) (models.Course, error)
	Update(course models.Course, courseRequest schemas.UpdateCourseRequest) (models.Course, error)
	Delete(id uint) error
	GetByIDWithChaptersCount(id uint) (schemas.CourseResponseWithChaptersCount, error)
}

type ChapterRepositoryInterface interface {
	GetByID(id, courseId uint) (models.Chapter, error)
	GetByIDWithLessonsCount(id uint, courseID uint) (schemas.ChapterResponseWithLessonsCount, error)
	GetByCourseID(courseID uint) ([]schemas.ChapterResponseWithLessonsCount, error)
	Create(chapter models.Chapter) (uint, error)
	Update(chapter models.Chapter) error
	Delete(id uint) error
}

type LessonRepositoryInterface interface {
	GetByID(courseID, chapterID, id uint) (models.Lesson, error)
	GetByChapterID(chapterID, courseID uint) ([]schemas.LessonResponse, error)
	Create(lesson models.Lesson) (uint, error)
	Update(lesson models.Lesson) error
	Delete(courseID, chapterID, id uint) error
}
