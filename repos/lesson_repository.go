package repos

import (
	"errors"
	"gorm.io/gorm"
	"web/models"
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

func (r *LessonRepository) GetAll() ([]models.Lesson, error) {
	var lessons []models.Lesson
	result := r.DB.Find(&lessons)
	if result.Error != nil {
		return nil, result.Error
	}

	return lessons, nil
}

func (r *LessonRepository) GetByChapterID(chapterID uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	result := r.DB.Where("chapter_id = ?", chapterID).Find(&lessons)
	if result.Error != nil {
		return nil, result.Error
	}

	return lessons, nil
}

func (r *LessonRepository) GetByID(id uint) (models.Lesson, error) {
	var lesson models.Lesson
	result := r.DB.First(&lesson, id)

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
		ChapterID:   lesson.ChapterID,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("lesson not found or no changes made")
	}

	return nil
}

func (r *LessonRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Lesson{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("lesson not found")
	}

	return nil
}
