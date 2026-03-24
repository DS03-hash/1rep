package userService

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"task-api/internal/domain"
)

// TestService_Create проверяет сценарии создания пользователя в сервисе.
func TestService_Create(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("create error")

	tests := []struct {
		name       string
		email      string
		password   string
		setupMock  func(repo *RepositoryMock)
		wantErr    error
		assertUser func(t *testing.T, got *User)
	}{
		{
			name:     "success",
			email:    "user@example.com",
			password: "secret",
			setupMock: func(repo *RepositoryMock) {
				repo.On("Create", mock.MatchedBy(func(u *User) bool {
					return u.Email == "user@example.com" && u.Password == "secret"
				})).
					Return(nil).
					Once()
			},
			assertUser: func(t *testing.T, got *User) {
				require.NotNil(t, got)
				assert.Equal(t, "user@example.com", got.Email)
				assert.Equal(t, "secret", got.Password)
			},
		},
		{
			name:     "invalid email",
			email:    "bad_email",
			password: "secret",
			wantErr:  ErrInvalidInput,
		},
		{
			name:     "empty password",
			email:    "user@example.com",
			password: "   ",
			wantErr:  ErrInvalidInput,
		},
		{
			name:     "repository returns error",
			email:    "user@example.com",
			password: "secret",
			setupMock: func(repo *RepositoryMock) {
				repo.On("Create", mock.MatchedBy(func(u *User) bool {
					return u.Email == "user@example.com" && u.Password == "secret"
				})).
					Return(repoErr).
					Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &RepositoryMock{}
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			svc := NewService(repo)

			got, err := svc.Create(tt.email, tt.password)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				if tt.assertUser != nil {
					tt.assertUser(t, got)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestService_List проверяет получение списка пользователей из сервиса.
func TestService_List(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo list error")
	expectedUsers := []User{
		{ID: 1, Email: "user1@example.com", Password: "secret1"},
		{ID: 2, Email: "user2@example.com", Password: "secret2"},
	}

	tests := []struct {
		name      string
		setupMock func(repo *RepositoryMock)
		wantUsers []User
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(repo *RepositoryMock) {
				repo.On("List").Return(expectedUsers, nil).Once()
			},
			wantUsers: expectedUsers,
		},
		{
			name: "repository returns error",
			setupMock: func(repo *RepositoryMock) {
				repo.On("List").Return([]User(nil), repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &RepositoryMock{}
			tt.setupMock(repo)
			svc := NewService(repo)

			got, err := svc.List()

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUsers, got)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestService_GetTasksForUser проверяет получение задач конкретного пользователя.
func TestService_GetTasksForUser(t *testing.T) {
	t.Parallel()

	getUserErr := errors.New("user not found")
	getTasksErr := errors.New("get tasks error")
	expectedTasks := []domain.Task{
		{ID: 1, Task: "first", IsDone: false, UserID: 5},
		{ID: 2, Task: "second", IsDone: true, UserID: 5},
	}

	tests := []struct {
		name      string
		userID    uint
		setupMock func(repo *RepositoryMock)
		wantTasks []domain.Task
		wantErr   error
	}{
		{
			name:   "user not found",
			userID: 17,
			setupMock: func(repo *RepositoryMock) {
				repo.On("GetByID", uint(17)).Return((*User)(nil), getUserErr).Once()
			},
			wantErr: ErrNotFound,
		},
		{
			name:   "repository get tasks error",
			userID: 6,
			setupMock: func(repo *RepositoryMock) {
				repo.On("GetByID", uint(6)).Return(&User{ID: 6}, nil).Once()
				repo.On("GetTasksForUser", uint(6)).Return([]domain.Task(nil), getTasksErr).Once()
			},
			wantErr: getTasksErr,
		},
		{
			name:   "success",
			userID: 5,
			setupMock: func(repo *RepositoryMock) {
				repo.On("GetByID", uint(5)).Return(&User{ID: 5}, nil).Once()
				repo.On("GetTasksForUser", uint(5)).Return(expectedTasks, nil).Once()
			},
			wantTasks: expectedTasks,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &RepositoryMock{}
			tt.setupMock(repo)
			svc := NewService(repo)

			got, err := svc.GetTasksForUser(tt.userID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantTasks, got)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestService_Patch проверяет частичное обновление пользователя.
func TestService_Patch(t *testing.T) {
	t.Parallel()

	getErr := errors.New("get error")
	updateErr := errors.New("update error")

	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name       string
		id         uint
		email      *string
		password   *string
		setupMock  func(repo *RepositoryMock)
		wantErr    error
		assertUser func(t *testing.T, got *User)
	}{
		{
			name:     "user not found",
			id:       7,
			email:    strPtr("new@example.com"),
			password: strPtr("new-secret"),
			setupMock: func(repo *RepositoryMock) {
				repo.On("GetByID", uint(7)).Return((*User)(nil), getErr).Once()
			},
			wantErr: ErrNotFound,
		},
		{
			name:     "invalid email input",
			id:       1,
			email:    strPtr("bad_email"),
			password: nil,
			setupMock: func(repo *RepositoryMock) {
				repo.On("GetByID", uint(1)).
					Return(&User{ID: 1, Email: "old@example.com", Password: "old-secret"}, nil).
					Once()
			},
			wantErr: ErrInvalidInput,
		},
		{
			name:     "invalid password input",
			id:       2,
			email:    nil,
			password: strPtr("   "),
			setupMock: func(repo *RepositoryMock) {
				repo.On("GetByID", uint(2)).
					Return(&User{ID: 2, Email: "old@example.com", Password: "old-secret"}, nil).
					Once()
			},
			wantErr: ErrInvalidInput,
		},
		{
			name:     "repository update error",
			id:       3,
			email:    strPtr("updated@example.com"),
			password: strPtr("new-secret"),
			setupMock: func(repo *RepositoryMock) {
				repo.On("GetByID", uint(3)).
					Return(&User{ID: 3, Email: "old@example.com", Password: "old-secret"}, nil).
					Once()
				repo.On("Update", mock.MatchedBy(func(u *User) bool {
					return u.ID == 3 && u.Email == "updated@example.com" && u.Password == "new-secret"
				})).
					Return(updateErr).
					Once()
			},
			wantErr: updateErr,
		},
		{
			name:     "successfully updates user",
			id:       9,
			email:    strPtr("new@example.com"),
			password: strPtr("new-password"),
			setupMock: func(repo *RepositoryMock) {
				repo.On("GetByID", uint(9)).
					Return(&User{ID: 9, Email: "old@example.com", Password: "old-password"}, nil).
					Once()
				repo.On("Update", mock.MatchedBy(func(u *User) bool {
					return u.ID == 9 && u.Email == "new@example.com" && u.Password == "new-password"
				})).
					Return(nil).
					Once()
			},
			assertUser: func(t *testing.T, got *User) {
				require.NotNil(t, got)
				assert.Equal(t, uint(9), got.ID)
				assert.Equal(t, "new@example.com", got.Email)
				assert.Equal(t, "new-password", got.Password)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &RepositoryMock{}
			tt.setupMock(repo)
			svc := NewService(repo)

			got, err := svc.Patch(tt.id, tt.email, tt.password)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				if tt.assertUser != nil {
					tt.assertUser(t, got)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestService_Delete проверяет удаление пользователя в сервисе.
func TestService_Delete(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("delete error")

	tests := []struct {
		name      string
		id        uint
		setupMock func(repo *RepositoryMock)
		wantErr   error
	}{
		{
			name: "repository returns error",
			id:   5,
			setupMock: func(repo *RepositoryMock) {
				repo.On("DeleteByID", uint(5)).Return(int64(0), repoErr).Once()
			},
			wantErr: repoErr,
		},
		{
			name: "user not found by id",
			id:   6,
			setupMock: func(repo *RepositoryMock) {
				repo.On("DeleteByID", uint(6)).Return(int64(0), nil).Once()
			},
			wantErr: ErrNotFound,
		},
		{
			name: "success",
			id:   8,
			setupMock: func(repo *RepositoryMock) {
				repo.On("DeleteByID", uint(8)).Return(int64(1), nil).Once()
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &RepositoryMock{}
			tt.setupMock(repo)
			svc := NewService(repo)

			err := svc.Delete(tt.id)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
