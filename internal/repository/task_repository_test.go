package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"task-api/internal/domain"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	cleanup := func() {
		mock.ExpectClose()
		require.NoError(t, sqlDB.Close())
		require.NoError(t, mock.ExpectationsWereMet())
	}

	return gdb, mock, cleanup
}

func TestGormTaskRepository_List(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "task", "is_done", "deleted_at"}).
			AddRow(1, "Write tests", false, nil).
			AddRow(2, "Review PR", true, nil)

		mock.ExpectQuery(`SELECT .* FROM "tasks" WHERE "tasks"\."deleted_at" IS NULL`).
			WillReturnRows(rows)

		repo := NewGormTaskRepository(db)
		got, err := repo.List()

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, uint(1), got[0].ID)
		assert.Equal(t, "Write tests", got[0].Task)
		assert.False(t, got[0].IsDone)
		assert.Equal(t, uint(2), got[1].ID)
		assert.Equal(t, "Review PR", got[1].Task)
		assert.True(t, got[1].IsDone)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		expectedErr := errors.New("db list error")
		mock.ExpectQuery(`SELECT .* FROM "tasks" WHERE "tasks"\."deleted_at" IS NULL`).
			WillReturnError(expectedErr)

		repo := NewGormTaskRepository(db)
		got, err := repo.List()

		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, got)
	})
}

func TestGormTaskRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		task := &domain.Task{ID: 3, Task: "Updated task", IsDone: true}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tasks" SET .* WHERE .*"id" = \$[0-9]+`).
			WithArgs("Updated task", true, sqlmock.AnyArg(), 3).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		repo := NewGormTaskRepository(db)
		err := repo.Update(task)

		require.NoError(t, err)
	})

	t.Run("exec error", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		task := &domain.Task{ID: 4, Task: "Task", IsDone: false}
		expectedErr := errors.New("update failed")

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tasks" SET .* WHERE .*"id" = \$[0-9]+`).
			WithArgs("Task", false, sqlmock.AnyArg(), 4).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		repo := NewGormTaskRepository(db)
		err := repo.Update(task)

		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestGormTaskRepository_DeleteByID(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"=.* WHERE "tasks"\."id" = \$[0-9]+ AND "tasks"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), 5).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		repo := NewGormTaskRepository(db)
		rows, err := repo.DeleteByID(5)

		require.NoError(t, err)
		assert.Equal(t, int64(1), rows)
	})

	t.Run("delete error", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		expectedErr := errors.New("delete failed")
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"=.* WHERE "tasks"\."id" = \$[0-9]+ AND "tasks"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), 6).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		repo := NewGormTaskRepository(db)
		rows, err := repo.DeleteByID(6)

		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, int64(0), rows)
	})
}
