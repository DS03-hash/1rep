package service

import (
	"errors"
	"strings"

	"task-api/internal/domain"
	"task-api/internal/repository"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)
// ErrNotFound - это ошибка, которая возвращается, когда запрашиваемая задача не найдена в базе данных.
// ErrInvalidInput - это ошибка, которая возвращается, когда входные данные для создания или обновления задачи недопустимы (например, пустая строка для задачи).

type TaskService struct {
	repo repository.TaskRepository
}
// TaskService - это структура, которая содержит ссылку на репозиторий задач.
// Она отвечает за бизнес-логику приложения, связанную с задачами. 
// Внутри TaskService определены методы для создания задачи, получения списка задач, обновления задачи и удаления задачи. 
// Эти методы взаимодействуют с репозиторием для выполнения операций над данными и возвращают результаты или ошибки в зависимости от ситуации.

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}
// NewTaskService - это функция-конструктор для TaskService. 
// Она принимает ссылку на TaskRepository и возвращает новый экземпляр TaskService, который будет использовать этот репозиторий для управления задачами.

func (s *TaskService) Create(task string, isDone bool) (*domain.Task, error) {
	if strings.TrimSpace(task) == "" {
		return nil, ErrInvalidInput
	}
	t := &domain.Task{Task: task, IsDone: isDone}
	if err := s.repo.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}
// Create - это метод, который создает новую задачу. Он принимает строку task и булевое значение isDone в качестве входных данных.
// Если строка task пустая или состоит только из пробелов, метод возвращает ошибку ErrInvalidInput. 
// В противном случае он создает новый объект domain.Task, сохраняет его в базе данных через репозиторий и возвращает созданную задачу или ошибку, если операция не удалась.

func (s *TaskService) List() ([]domain.Task, error) {
	return s.repo.List()
}
// List - это метод, который возвращает список всех задач. Он просто вызывает метод List репозитория и возвращает результат или ошибку, если операция не удалась.

func (s *TaskService) Patch(id uint, task *string, isDone *bool) (*domain.Task, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {

		return nil, ErrNotFound
	}

	if task != nil {
		if strings.TrimSpace(*task) == "" {
			return nil, ErrInvalidInput
		}
		t.Task = *task
	}
	if isDone != nil {
		t.IsDone = *isDone
	}

	if err := s.repo.Update(t); err != nil {
		return nil, err
	}
	return t, nil
}
// Patch - это метод, который обновляет существующую задачу. Он принимает идентификатор задачи, а также указатели на строку task и булевое значение isDone, которые могут быть nil.
// Сначала метод пытается получить задачу по идентификатору. Если задача не найдена, возвращается ошибка ErrNotFound. 
// Затем, если task не nil, проверяется, что строка не пустая или не состоит только из пробелов. Если это так, возвращается ошибка ErrInvalidInput. 
// В противном случае обновляется поле Task. Аналогично, если isDone не nil, обновляется поле IsDone. 
// Наконец, обновленная задача сохраняется в базе данных через репозиторий и возвращается или возвращается ошибка, если операция не удалась.

func (s *TaskService) Delete(id uint) error {
	rows, err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
// Delete - это метод, который удаляет задачу по ее идентификатору. Он вызывает метод DeleteByID репозитория и получает количество удаленных строк и ошибку, если операция не удалась.
// Если возникает ошибка при удалении, она возвращается. Если количество удаленных строк равно нулю, это означает, что задача не была найдена, и возвращается ошибка ErrNotFound. 
// В противном случае возвращается nil, что означает успешное удаление задачи.