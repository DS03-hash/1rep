package userService

import "github.com/stretchr/testify/mock"

type RepositoryMock struct {
	mock.Mock
}

// Create mocks repository create.
func (m *RepositoryMock) Create(u *User) error {
	args := m.Called(u)
	return args.Error(0)
}

// List mocks repository list.
func (m *RepositoryMock) List() ([]User, error) {
	args := m.Called()
	users, _ := args.Get(0).([]User)
	return users, args.Error(1)
}

// GetByID mocks repository get by id.
func (m *RepositoryMock) GetByID(id uint) (*User, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(*User)
	return user, args.Error(1)
}

// Update mocks repository update.
func (m *RepositoryMock) Update(u *User) error {
	args := m.Called(u)
	return args.Error(0)
}

// DeleteByID mocks repository delete by id.
func (m *RepositoryMock) DeleteByID(id uint) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}
