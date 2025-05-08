package services

import (
	"web/models"
	"web/schemas"
)

type CourseServiceInterface interface {
	GetAllCourses() ([]schemas.CourseResponseWithChaptersCount, error)
	GetCourseByID(id uint) (models.Course, error)
	GetCourseByIDWithChapterCount(id uint) (schemas.CourseResponseWithChaptersCount, error)
	CreateCourse(courseDTO schemas.CreateCourseRequest) (schemas.CourseResponse, error)
	UpdateCourse(course models.Course, courseRequest schemas.UpdateCourseRequest) (schemas.CourseResponse, error)
	DeleteCourse(id uint) error
}

type ChapterServiceInterface interface {
	GetAllChapters() ([]models.Chapter, error)
	GetChapterByID(id uint) (schemas.ChapterResponseWithLessonsCount, error)
	GetChaptersByCourseID(courseID uint) ([]schemas.ChapterResponseWithLessonsCount, error)
	CountChaptersByCourseID(courseID uint) (int, error)
	CreateChapter(chapter models.Chapter) (uint, error)
	UpdateChapter(chapter models.Chapter) error
	DeleteChapter(id uint) error
}

type LessonServiceInterface interface {
	GetLessonByID(id uint) (models.Lesson, error)
	GetLessonsByChapterID(chapterID, courseID uint) ([]schemas.LessonResponse, error)
	CreateLesson(lessonRequest schemas.LessonRequest, chapterId uint) (uint, error)
	UpdateLesson(courseID, chapterID, id uint, lessonRequest schemas.LessonRequest) error
	DeleteLesson(id uint) error
}
