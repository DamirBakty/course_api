package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	tableName   struct{}       `gorm:"table:course"`
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Chapters    []Chapter      `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
