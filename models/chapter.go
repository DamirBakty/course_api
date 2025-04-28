package models

import (
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	tableName   struct{}       `gorm:"table:chapter"`
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Order       int            `gorm:"not null"`
	CourseID    uint           `gorm:"not null"`
	Course      Course         `gorm:"foreignKey:CourseID"`
	Lessons     []Lesson       `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
