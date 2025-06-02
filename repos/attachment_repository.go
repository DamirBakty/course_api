package repos

import (
	"errors"
	"gorm.io/gorm"
	"web/models"
)

type AttachmentRepositoryInterface interface {
	GetByID(id uint) (models.Attachment, error)
	GetByLessonID(lessonID uint) ([]models.Attachment, error)
	Create(attachment models.Attachment) (uint, error)
	Delete(id uint) error
}

var _ AttachmentRepositoryInterface = (*AttachmentRepository)(nil)

type AttachmentRepository struct {
	DB *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{
		DB: db,
	}
}

func (r *AttachmentRepository) GetByID(id uint) (models.Attachment, error) {
	var attachment models.Attachment
	result := r.DB.First(&attachment, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return attachment, errors.New("attachment not found")
		}
		return attachment, result.Error
	}
	return attachment, nil
}

func (r *AttachmentRepository) GetByLessonID(lessonID uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	result := r.DB.Where("lesson_id = ?", lessonID).Find(&attachments)
	if result.Error != nil {
		return nil, result.Error
	}
	return attachments, nil
}

func (r *AttachmentRepository) Create(attachment models.Attachment) (uint, error) {
	result := r.DB.Create(&attachment)
	if result.Error != nil {
		return 0, result.Error
	}
	return attachment.ID, nil
}

func (r *AttachmentRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Attachment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("attachment not found")
	}
	return nil
}