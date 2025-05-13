package services_test

import (
	"testing"
	"web/services"
)

func TestLessonServiceMethods(t *testing.T) {
	
	var lessonService *services.LessonService

	_ = lessonService.GetLessonsByChapterID
	_ = lessonService.GetLessonByID
	_ = lessonService.CreateLesson
	_ = lessonService.UpdateLesson
	_ = lessonService.DeleteLesson
}
