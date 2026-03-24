package domain

import "gorm.io/gorm"

type Task struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	UserID    uint           `json:"user_id"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
