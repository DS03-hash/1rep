package userService

import (
	"github.com/stretchr/testify/mock"

	"task-api/internal/domain"
)

type RepositoryMock struct {
	mock.Mock
}

// Create имитирует вызов создания пользователя в репозитории.
func (m *RepositoryMock) Create(u *User) error {
	args := m.Called(u)
	return args.Error(0)
}

// List имитирует получение списка пользователей из репозитория.
func (m *RepositoryMock) List() ([]User, error) {
	args := m.Called()
	users, _ := args.Get(0).([]User)
	return users, args.Error(1)
}

// GetByID имитирует поиск пользователя по идентификатору.
func (m *RepositoryMock) GetByID(id uint) (*User, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(*User)
	return user, args.Error(1)
}

// GetTasksForUser имитирует получение задач пользователя из репозитория.
func (m *RepositoryMock) GetTasksForUser(userID uint) ([]domain.Task, error) {
	args := m.Called(userID)
	tasks, _ := args.Get(0).([]domain.Task)
	return tasks, args.Error(1)
}

// Update имитирует обновление пользователя в репозитории.
func (m *RepositoryMock) Update(u *User) error {
	args := m.Called(u)
	return args.Error(0)
}

// DeleteByID имитирует удаление пользователя по идентификатору.
func (m *RepositoryMock) DeleteByID(id uint) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}
