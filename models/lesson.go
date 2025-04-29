package models

import (
	"time"

	"gorm.io/gorm"
)

// Lesson represents a lesson in a chapter
// swagger:model
type Lesson struct {
	tableName   struct{}       `gorm:"table:lesson"`
	ID          uint           `gorm:"primaryKey" json:"id,omitempty" example:"1"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name" example:"Lesson 1: Introduction"`
	Description string         `gorm:"type:text" json:"description" example:"Overview of the chapter content"`
	Content     string         `gorm:"type:text" json:"content" example:"This lesson covers the basic concepts of the chapter."`
	Order       int            `gorm:"not null" json:"order" example:"1"`
	ChapterID   uint           `gorm:"not null" json:"chapter_id,omitempty" example:"1"`
	Chapter     Chapter        `gorm:"foreignKey:ChapterID" json:"chapter,omitempty"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Lesson) TableName() string {
	return "lesson"
}
