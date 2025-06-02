package models

import (
	"time"

	"gorm.io/gorm"
)

// Attachment represents a file attached to a lesson
// swagger:model
type Attachment struct {
	tableName struct{}       `gorm:"table:attachment"`
	ID        uint           `gorm:"primaryKey" json:"id,omitempty" example:"1"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name" example:"lecture_slides.pdf"`
	URL       string         `gorm:"type:varchar(255);not null" json:"url" example:"https://storage.example.com/files/lecture_slides.pdf"`
	LessonID  uint           `gorm:"not null" json:"lesson_id,omitempty" example:"1"`
	Lesson    Lesson         `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Attachment) TableName() string {
	return "attachment"
}