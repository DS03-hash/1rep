package userService

import (
	"time"

	"task-api/internal/domain"

	"gorm.io/gorm"
)

// User описывает пользователя в БД.
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password  string         `json:"password" gorm:"size:255;not null"`
	Tasks     []domain.Task  `json:"tasks,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
