package services_test

import (
	"testing"
	"web/services"
)

func TestChapterServiceImplementsInterface(t *testing.T) {

	var _ services.ChapterServiceInterface = (*services.ChapterService)(nil)
}
