package userService

import (
	"task-api/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Create(u *User) error
	List() ([]User, error)
	GetByID(id uint) (*User, error)
	GetTasksForUser(userID uint) ([]domain.Task, error)
	Update(u *User) error
	DeleteByID(id uint) (rowsAffected int64, err error)
}

type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository создает репозиторий пользователей на базе GORM.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create сохраняет нового пользователя в базе данных.
func (r *GormRepository) Create(u *User) error {
	return r.db.Create(u).Error
}

// List возвращает всех не удаленных пользователей.
func (r *GormRepository) List() ([]User, error) {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByID ищет пользователя по идентификатору.
func (r *GormRepository) GetByID(id uint) (*User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// GetTasksForUser возвращает задачи конкретного пользователя.
func (r *GormRepository) GetTasksForUser(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	if err := r.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// Update обновляет существующего пользователя.
func (r *GormRepository) Update(u *User) error {
	return r.db.Save(u).Error
}

// DeleteByID выполняет мягкое удаление пользователя по идентификатору.
func (r *GormRepository) DeleteByID(id uint) (int64, error) {
	res := r.db.Delete(&User{}, id)
	return res.RowsAffected, res.Error
}
