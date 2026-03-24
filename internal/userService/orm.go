package userService

import (
	"time"

	"gorm.io/gorm"
)

// User describes user entity in database.
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password  string         `json:"password" gorm:"size:255;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
