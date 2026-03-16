package repository

import (
	"gorm.io/gorm"

	"task-api/internal/domain"
)

type TaskRepository interface {
	Create(t *domain.Task) error
	List() ([]domain.Task, error)
	GetByID(id uint) (*domain.Task, error)
	Update(t *domain.Task) error
	DeleteByID(id uint) (rowsAffected int64, err error)
}

type GormTaskRepository struct {
	db *gorm.DB
}

// NewGormTaskRepository создает репозиторий задач на базе GORM.
func NewGormTaskRepository(db *gorm.DB) *GormTaskRepository {
	return &GormTaskRepository{db: db}
}

// Create сохраняет новую задачу в базе данных.
func (r *GormTaskRepository) Create(t *domain.Task) error {
	return r.db.Create(t).Error
}

// List возвращает все не удаленные задачи.
func (r *GormTaskRepository) List() ([]domain.Task, error) {
	var tasks []domain.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetByID ищет задачу по идентификатору.
func (r *GormTaskRepository) GetByID(id uint) (*domain.Task, error) {
	var t domain.Task
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// Update обновляет существующую задачу.
func (r *GormTaskRepository) Update(t *domain.Task) error {
	return r.db.Save(t).Error
}

// DeleteByID выполняет мягкое удаление задачи по идентификатору.
func (r *GormTaskRepository) DeleteByID(id uint) (int64, error) {
	res := r.db.Delete(&domain.Task{}, id)
	return res.RowsAffected, res.Error
}
