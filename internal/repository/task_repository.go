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

func NewGormTaskRepository(db *gorm.DB) *GormTaskRepository {
	return &GormTaskRepository{db: db}
}

func (r *GormTaskRepository) Create(t *domain.Task) error {
	return r.db.Create(t).Error
}

func (r *GormTaskRepository) List() ([]domain.Task, error) {
	var tasks []domain.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *GormTaskRepository) GetByID(id uint) (*domain.Task, error) {
	var t domain.Task
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormTaskRepository) Update(t *domain.Task) error {
	return r.db.Save(t).Error
}

func (r *GormTaskRepository) DeleteByID(id uint) (int64, error) {
	res := r.db.Delete(&domain.Task{}, id)
	return res.RowsAffected, res.Error
}
