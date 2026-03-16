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
// TaskRepository - это интерфейс, который определяет методы для взаимодействия с хранилищем данных задач.
// Он включает методы для создания, получения, обновления и удаления задач. 
// Реализации этого интерфейса будут использоваться сервисом для управления задачами без необходимости знать детали реализации хранилища данных.

type GormTaskRepository struct {
	db *gorm.DB
}
// GormTaskRepository - это структура, которая реализует интерфейс TaskRepository с использованием GORM для взаимодействия с базой данных.

func NewGormTaskRepository(db *gorm.DB) *GormTaskRepository {
	return &GormTaskRepository{db: db}
}
//	NewGormTaskRepository - это функция-конструктор, которая принимает указатель на gorm.DB и возвращает новый экземпляр GormTaskRepository.

func (r *GormTaskRepository) Create(t *domain.Task) error {
	return r.db.Create(t).Error
}
// Create - это метод, который сохраняет новую задачу в базе данных. Он использует метод Create GORM и возвращает ошибку, если операция не удалась.

func (r *GormTaskRepository) List() ([]domain.Task, error) {
	var tasks []domain.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
// List - это метод, который извлекает все задачи из базы данных. Он использует метод Find GORM и возвращает срез задач и ошибку, если операция не удалась.

func (r *GormTaskRepository) GetByID(id uint) (*domain.Task, error) {
	var t domain.Task
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
// GetByID - это метод, который извлекает задачу по ее идентификатору из базы данных. Он использует метод First GORM и 
// возвращает указатель на задачу и ошибку, если операция не удалась.

func (r *GormTaskRepository) Update(t *domain.Task) error {
	return r.db.Save(t).Error
}
// Update - это метод, который обновляет существующую задачу в базе данных. Он использует метод Save GORM и возвращает ошибку, 
// если операция не удалась.

func (r *GormTaskRepository) DeleteByID(id uint) (int64, error) {
	res := r.db.Delete(&domain.Task{}, id)
	return res.RowsAffected, res.Error
}
// DeleteByID - это метод, который удаляет задачу по ее идентификатору из базы данных. 
// Он использует метод Delete GORM и возвращает количество удаленных строк и ошибку, если операция не удалась.
