package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"task-api/internal/domain"
)

func TestTaskService_Create(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("create error")

	tests := []struct {
		name       string
		task       string
		isDone     bool
		setupMock  func(repo *TaskRepositoryMock)
		wantErr    error
		assertTask func(t *testing.T, got *domain.Task)
	}{
		{
			name:   "success",
			task:   "write unit tests",
			isDone: false,
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("Create", mock.MatchedBy(func(task *domain.Task) bool {
					return task.Task == "write unit tests" && !task.IsDone
				})).
					Return(nil).
					Once()
			},
			assertTask: func(t *testing.T, got *domain.Task) {
				require.NotNil(t, got)
				assert.Equal(t, "write unit tests", got.Task)
				assert.False(t, got.IsDone)
			},
		},
		{
			name:    "invalid input",
			task:    "   ",
			isDone:  true,
			wantErr: ErrInvalidInput,
		},
		{
			name:   "repository returns error",
			task:   "create task",
			isDone: true,
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("Create", mock.MatchedBy(func(task *domain.Task) bool {
					return task.Task == "create task" && task.IsDone
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

			repo := &TaskRepositoryMock{}
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			svc := NewTaskService(repo)

			got, err := svc.Create(tt.task, tt.isDone)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				if tt.assertTask != nil {
					tt.assertTask(t, got)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestTaskService_List(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo list error")
	expectedTasks := []domain.Task{
		{ID: 1, Task: "Write tests", IsDone: false},
		{ID: 2, Task: "Review code", IsDone: true},
	}

	tests := []struct {
		name      string
		setupMock func(repo *TaskRepositoryMock)
		wantTasks []domain.Task
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("List").Return(expectedTasks, nil).Once()
			},
			wantTasks: expectedTasks,
		},
		{
			name: "repository returns error",
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("List").Return([]domain.Task(nil), repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &TaskRepositoryMock{}
			tt.setupMock(repo)
			svc := NewTaskService(repo)

			got, err := svc.List()

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

func TestTaskService_Patch(t *testing.T) {
	t.Parallel()

	getErr := errors.New("get error")
	updateErr := errors.New("update error")

	strPtr := func(s string) *string { return &s }
	boolPtr := func(v bool) *bool { return &v }

	tests := []struct {
		name       string
		id         uint
		task       *string
		isDone     *bool
		setupMock  func(repo *TaskRepositoryMock)
		wantErr    error
		assertTask func(t *testing.T, got *domain.Task)
	}{
		{
			name:   "task not found",
			id:     7,
			task:   strPtr("new"),
			isDone: boolPtr(true),
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("GetByID", uint(7)).Return((*domain.Task)(nil), getErr).Once()
			},
			wantErr: ErrNotFound,
		},
		{
			name:   "invalid task input",
			id:     1,
			task:   strPtr("   "),
			isDone: nil,
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("GetByID", uint(1)).
					Return(&domain.Task{ID: 1, Task: "old", IsDone: false}, nil).
					Once()
			},
			wantErr: ErrInvalidInput,
		},
		{
			name:   "repository update error",
			id:     3,
			task:   strPtr("updated"),
			isDone: boolPtr(true),
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("GetByID", uint(3)).
					Return(&domain.Task{ID: 3, Task: "old", IsDone: false}, nil).
					Once()
				repo.On("Update", mock.MatchedBy(func(task *domain.Task) bool {
					return task.ID == 3 && task.Task == "updated" && task.IsDone
				})).
					Return(updateErr).
					Once()
			},
			wantErr: updateErr,
		},
		{
			name:   "successfully updates task",
			id:     9,
			task:   strPtr("new text"),
			isDone: boolPtr(true),
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("GetByID", uint(9)).
					Return(&domain.Task{ID: 9, Task: "old text", IsDone: false}, nil).
					Once()
				repo.On("Update", mock.MatchedBy(func(task *domain.Task) bool {
					return task.ID == 9 && task.Task == "new text" && task.IsDone
				})).
					Return(nil).
					Once()
			},
			assertTask: func(t *testing.T, got *domain.Task) {
				require.NotNil(t, got)
				assert.Equal(t, uint(9), got.ID)
				assert.Equal(t, "new text", got.Task)
				assert.True(t, got.IsDone)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &TaskRepositoryMock{}
			tt.setupMock(repo)
			svc := NewTaskService(repo)

			got, err := svc.Patch(tt.id, tt.task, tt.isDone)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				if tt.assertTask != nil {
					tt.assertTask(t, got)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("delete error")

	tests := []struct {
		name      string
		id        uint
		setupMock func(repo *TaskRepositoryMock)
		wantErr   error
	}{
		{
			name: "repository returns error",
			id:   5,
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("DeleteByID", uint(5)).Return(int64(0), repoErr).Once()
			},
			wantErr: repoErr,
		},
		{
			name: "task not found by id",
			id:   6,
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("DeleteByID", uint(6)).Return(int64(0), nil).Once()
			},
			wantErr: ErrNotFound,
		},
		{
			name: "success",
			id:   8,
			setupMock: func(repo *TaskRepositoryMock) {
				repo.On("DeleteByID", uint(8)).Return(int64(1), nil).Once()
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &TaskRepositoryMock{}
			tt.setupMock(repo)
			svc := NewTaskService(repo)

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
