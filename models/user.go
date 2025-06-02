package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
// swagger:model
type User struct {
	tableName struct{}       `gorm:"table:users"`
	ID        uint           `gorm:"primaryKey" json:"id,omitempty" example:"1"`
	Username  string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"username" example:"johndoe"`
	Email     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email" example:"john.doe@example.com"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // Password is not exposed in JSON
	Roles     string         `gorm:"type:varchar(255);not null" json:"roles" example:"user,admin"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Sub       string         `gorm:"type:varchar(36);not null" json:"sub"`
}

func (User) TableName() string {
	return "users"
}
