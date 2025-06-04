package models

import (
	"time"

	"gorm.io/gorm"
)

// Chapter represents a chapter in a course
// swagger:model
type Chapter struct {
	tableName   struct{}       `gorm:"table:chapter"`
	ID          uint           `gorm:"primaryKey" json:"id,omitempty" example:"1"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name" example:"Chapter 1: Getting Started"`
	Description string         `gorm:"type:text" json:"description" example:"Introduction to the course material"`
	Order       int            `gorm:"not null" json:"order" example:"1"`
	CourseID    uint           `gorm:"not null" json:"course_id,omitempty" example:"1"`
	Course      Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	CreatedBy   *uint          `gorm:"column:created_by" json:"created_by,omitempty"`
	Creator     *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Lessons     []Lesson       `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE" json:"lessons,omitempty"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Chapter) TableName() string {
	return "chapter"
}
