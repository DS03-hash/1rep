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

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

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

func (s *TaskService) List() ([]domain.Task, error) {
	return s.repo.List()
}

func (s *TaskService) Patch(id uint, task *string, isDone *bool) (*domain.Task, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		// здесь можно точнее различать ErrRecordNotFound, но минимально оставим так
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
