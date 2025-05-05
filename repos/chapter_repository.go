package repos

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"web/models"
	"web/schemas"
)

var _ ChapterRepositoryInterface = (*ChapterRepository)(nil)

type ChapterRepository struct {
	DB *gorm.DB
}

func NewChapterRepository(db *gorm.DB) *ChapterRepository {
	return &ChapterRepository{
		DB: db,
	}
}

func (r *ChapterRepository) GetByCourseID(courseID uint) ([]schemas.ChapterResponse, error) {
	var chapters []models.Chapter
	result := r.DB.Where("course_id = ?", courseID).Order(`"order" ASC`).Find(&chapters)
	if result.Error != nil {
		return nil, result.Error
	}
	var chapterResponses []schemas.ChapterResponse

	for _, chapter := range chapters {
		chapterResponses = append(chapterResponses, schemas.ChapterResponse{
			ID:          chapter.ID,
			Name:        chapter.Name,
			Description: chapter.Description,
			CreatedAt:   chapter.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   chapter.UpdatedAt.Format(time.RFC3339),
		})
	}

	return chapterResponses, nil
}

func (r *ChapterRepository) GetByID(id uint) (models.Chapter, error) {
	var chapter models.Chapter
	result := r.DB.First(&chapter, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return chapter, errors.New("chapter not found")
		}
		return chapter, result.Error
	}

	return chapter, nil
}

func (r *ChapterRepository) Create(chapter models.Chapter) (uint, error) {
	result := r.DB.Create(&chapter)
	if result.Error != nil {
		return 0, result.Error
	}

	return chapter.ID, nil
}

func (r *ChapterRepository) Update(chapter models.Chapter) error {
	result := r.DB.Model(&chapter).Updates(models.Chapter{
		Name:        chapter.Name,
		Description: chapter.Description,
		Order:       chapter.Order,
		CourseID:    chapter.CourseID,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("chapter not found or no changes made")
	}

	return nil
}

func (r *ChapterRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Chapter{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("chapter not found")
	}

	return nil
}
