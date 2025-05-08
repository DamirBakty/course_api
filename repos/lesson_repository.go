package repos

import (
	"errors"
	"gorm.io/gorm"
	"web/models"
	"web/schemas"
)

var _ LessonRepositoryInterface = (*LessonRepository)(nil)

type LessonRepository struct {
	DB *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{
		DB: db,
	}
}

func (r *LessonRepository) GetByChapterID(courseID, chapterID uint) ([]schemas.LessonResponse, error) {
	var lessons []schemas.LessonResponse
	result := r.DB.Model(&models.Lesson{}).
		Select("lesson.id, lesson.name, lesson.description, lesson.content, lesson.order, lesson.created_at").
		Joins("INNER JOIN chapter ON chapter.id = lesson.chapter_id").
		Where("chapter_id = ? and chapter.course_id = ?", chapterID, courseID).
		Order(
			"lesson.order ASC",
		).
		Find(&lessons)
	if result.Error != nil {
		return nil, result.Error
	}

	return lessons, nil
}

func (r *LessonRepository) GetByID(courseID, chapterID, id uint) (models.Lesson, error) {
	var lesson models.Lesson
	result := r.DB.Model(&models.Lesson{}).
		Select("lesson.id, lesson.name, lesson.description, lesson.content, lesson.order, lesson.created_at").
		Joins("INNER JOIN chapter ON chapter.id = lesson.chapter_id").
		Where("chapter_id = ? and chapter.course_id = ? and id = ?", chapterID, courseID, id).
		Find(&lesson)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return lesson, errors.New("lesson not found")
		}
		return lesson, result.Error
	}

	return lesson, nil
}

func (r *LessonRepository) Create(lesson models.Lesson) (uint, error) {
	result := r.DB.Create(&lesson)
	if result.Error != nil {
		return 0, result.Error
	}

	return lesson.ID, nil
}

func (r *LessonRepository) Update(lesson models.Lesson) error {
	result := r.DB.Model(&lesson).Updates(models.Lesson{
		Name:        lesson.Name,
		Description: lesson.Description,
		Content:     lesson.Content,
		Order:       lesson.Order,
		UpdatedAt:   lesson.UpdatedAt,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("lesson not found or no changes made")
	}

	return nil
}

func (r *LessonRepository) Delete(courseID, chapterID, id uint) error {
	result := r.DB.
		Delete(&models.Lesson{}).
		Joins("INNER JOIN chapter ON chapter.id = lesson.chapter_id").
		Where("lesson.id = ? AND lesson.chapter_id = ? AND chapter.course_id = ?", id, chapterID, courseID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("lesson not found")
	}

	return nil
}
