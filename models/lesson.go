package models

import (
	"time"

	"gorm.io/gorm"
)

type Lesson struct {
	tableName   struct{}       `gorm:"table:lesson"`
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Content     string         `gorm:"type:text"`
	Order       int            `gorm:"not null"`
	ChapterID   uint           `gorm:"not null"`
	Chapter     Chapter        `gorm:"foreignKey:ChapterID"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Lesson) TableName() string {
	return "lesson"
}
