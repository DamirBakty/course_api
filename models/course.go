package models

import (
	"time"

	"gorm.io/gorm"
)

// Course represents a course in the system
// swagger:model
type Course struct {
	tableName   struct{}       `gorm:"table:course"`
	ID          uint           `gorm:"primaryKey" json:"id,omitempty" example:"1"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name" example:"Introduction to Go Programming"`
	Description string         `gorm:"type:text" json:"description" example:"Learn the basics of Go programming language"`
	CreatedBy   *uint          `gorm:"column:created_by" json:"created_by,omitempty"`
	Creator     *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Chapters    []Chapter      `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"chapters,omitempty"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Course) TableName() string {
	return "course"
}
