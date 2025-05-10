package repos

import (
	"errors"
	"gorm.io/gorm"
	"web/models"
	"web/schemas"
)

type ChapterRepositoryInterface interface {
	GetByID(id, courseId uint) (models.Chapter, error)
	GetByIDWithLessonsCount(id uint, courseID uint) (schemas.ChapterResponseWithLessonsCount, error)
	GetByCourseID(courseID uint) ([]schemas.ChapterResponseWithLessonsCount, error)
	Create(chapter models.Chapter) (uint, error)
	Update(chapter models.Chapter) error
	Delete(id uint) error
}

var _ ChapterRepositoryInterface = (*ChapterRepository)(nil)

type ChapterRepository struct {
	DB *gorm.DB
}

func NewChapterRepository(db *gorm.DB) *ChapterRepository {
	return &ChapterRepository{
		DB: db,
	}
}

func (r *ChapterRepository) GetByCourseID(courseID uint) ([]schemas.ChapterResponseWithLessonsCount, error) {
	var chapterResponses []schemas.ChapterResponseWithLessonsCount

	subQuery := r.DB.Model(&models.Lesson{}).
		Select("chapter_id, count(*) as lessons_count").
		Group("chapter_id")

	err := r.DB.Model(&models.Chapter{}).
		Select("chapter.id, chapter.name, chapter.description, chapter.created_at, chapter.updated_at, COALESCE(lessons_count, 0) as lessons_count").
		Joins("LEFT JOIN (?) AS lesson_counts ON chapter.id = lesson_counts.chapter_id", subQuery).
		Where("chapter.course_id = ?", courseID).
		Order(`"order" ASC`).
		Scan(&chapterResponses).Error

	if err != nil {
		return nil, err
	}

	return chapterResponses, nil
}
func (r *ChapterRepository) GetByID(id, courseId uint) (models.Chapter, error) {
	var chapter models.Chapter
	result := r.DB.Where("course_id = ? and id = ?", courseId, id).First(&chapter)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return chapter, errors.New("chapter not found")
		}
		return chapter, result.Error
	}

	return chapter, nil
}

func (r *ChapterRepository) GetByIDWithLessonsCount(id, courseID uint) (schemas.ChapterResponseWithLessonsCount, error) {
	var chapterResponse schemas.ChapterResponseWithLessonsCount

	subQuery := r.DB.Model(&models.Lesson{}).
		Select("chapter_id, count(*) as lessons_count").
		Where("chapter_id = ?", id).
		Group("chapter_id")

	err := r.DB.Model(&models.Chapter{}).
		Select("chapter.id, chapter.name, chapter.description, chapter.created_at, chapter.updated_at, COALESCE(lessons_count, 0) as lessons_count").
		Joins("LEFT JOIN (?) AS lesson_counts ON chapter.id = lesson_counts.chapter_id", subQuery).
		Where("chapter.id = ? and chapter.course_id = ?", id, courseID).
		Scan(&chapterResponse).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return chapterResponse, errors.New("chapter not found")
		}
		return chapterResponse, err
	}

	if chapterResponse.ID == 0 {
		return chapterResponse, errors.New("chapter not found")
	}

	return chapterResponse, nil
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
