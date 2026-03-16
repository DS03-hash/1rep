package domain

import "gorm.io/gorm"

type Task struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
// Сущность обработки
// Task - это структура, которая представляет собой задачу. Она содержит поля ID, Task, IsDone и DeletedAt. ID - это уникальный идентификатор задачи, Task - это текст задачи, IsDone - это булевое значение, которое указывает, выполнена ли задача, а DeletedAt - это поле для мягкого удаления задачи.