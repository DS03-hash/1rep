package service

import (
	"github.com/stretchr/testify/mock"

	"task-api/internal/domain"
)

type TaskRepositoryMock struct {
	mock.Mock
}

func (m *TaskRepositoryMock) Create(t *domain.Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *TaskRepositoryMock) List() ([]domain.Task, error) {
	args := m.Called()
	tasks, _ := args.Get(0).([]domain.Task)
	return tasks, args.Error(1)
}

func (m *TaskRepositoryMock) GetByID(id uint) (*domain.Task, error) {
	args := m.Called(id)
	task, _ := args.Get(0).(*domain.Task)
	return task, args.Error(1)
}

func (m *TaskRepositoryMock) Update(t *domain.Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *TaskRepositoryMock) DeleteByID(id uint) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

