package services

import (
	"web/models"
	"web/schemas"
)

type CourseServiceInterface interface {
	GetAllCourses() ([]models.Course, error)
	GetCourseByID(id uint) (models.Course, error)
	CreateCourse(courseDTO schemas.CreateCourseRequest) (models.Course, error)
	UpdateCourse(course models.Course, courseRequest schemas.UpdateCourseRequest) (models.Course, error)
	DeleteCourse(id uint) error
}

type ChapterServiceInterface interface {
	GetAllChapters() ([]models.Chapter, error)
	GetChapterByID(id uint) (models.Chapter, error)
	GetChaptersByCourseID(courseID uint) ([]models.Chapter, error)
	CreateChapter(chapter models.Chapter) (uint, error)
	UpdateChapter(chapter models.Chapter) error
	DeleteChapter(id uint) error
}

type LessonServiceInterface interface {
	GetAllLessons() ([]models.Lesson, error)
	GetLessonByID(id uint) (models.Lesson, error)
	GetLessonsByChapterID(chapterID uint) ([]models.Lesson, error)
	CreateLesson(lesson models.Lesson) (uint, error)
	UpdateLesson(lesson models.Lesson) error
	DeleteLesson(id uint) error
}
