package userService

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newMockDB creates sqlmock-backed gorm database.
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

func TestGormRepository_List(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "user1@example.com", "secret1", now, now, nil).
			AddRow(2, "user2@example.com", "secret2", now, now, nil)

		mock.ExpectQuery(`SELECT .* FROM "users" WHERE "users"\."deleted_at" IS NULL`).
			WillReturnRows(rows)

		repo := NewGormRepository(db)
		got, err := repo.List()

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, uint(1), got[0].ID)
		assert.Equal(t, "user1@example.com", got[0].Email)
		assert.Equal(t, "secret1", got[0].Password)
		assert.Equal(t, uint(2), got[1].ID)
		assert.Equal(t, "user2@example.com", got[1].Email)
		assert.Equal(t, "secret2", got[1].Password)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		expectedErr := errors.New("db list error")
		mock.ExpectQuery(`SELECT .* FROM "users" WHERE "users"\."deleted_at" IS NULL`).
			WillReturnError(expectedErr)

		repo := NewGormRepository(db)
		got, err := repo.List()

		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, got)
	})
}

func TestGormRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		now := time.Now()
		user := &User{
			ID:        3,
			Email:     "updated@example.com",
			Password:  "updated-password",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET .* WHERE .*"id" = \$[0-9]+`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		repo := NewGormRepository(db)
		err := repo.Update(user)

		require.NoError(t, err)
	})

	t.Run("exec error", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		now := time.Now()
		user := &User{
			ID:        4,
			Email:     "user@example.com",
			Password:  "password",
			CreatedAt: now,
			UpdatedAt: now,
		}
		expectedErr := errors.New("update failed")

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET .* WHERE .*"id" = \$[0-9]+`).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		repo := NewGormRepository(db)
		err := repo.Update(user)

		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestGormRepository_DeleteByID(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "deleted_at"=.* WHERE "users"\."id" = \$[0-9]+ AND "users"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), 5).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		repo := NewGormRepository(db)
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
		mock.ExpectExec(`UPDATE "users" SET "deleted_at"=.* WHERE "users"\."id" = \$[0-9]+ AND "users"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), 6).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		repo := NewGormRepository(db)
		rows, err := repo.DeleteByID(6)

		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, int64(0), rows)
	})
}
