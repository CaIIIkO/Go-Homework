package postrgesql

import (
	"context"
	"database/sql"
	"errors"
	"homework-3/internal/pkg/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_GetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,address,contact FROM pvz where id=$1", gomock.Any()).
			Return(nil)
		// act
		user, err := s.repo.GetByID(ctx, id)
		// assert
		require.NoError(t, err)
		assert.Equal(t, int64(0), user.ID)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()
			s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,address,contact FROM pvz where id=$1", gomock.Any()).
				Return(sql.ErrNoRows)
			// act
			user, err := s.repo.GetByID(ctx, id)
			// assert
			require.EqualError(t, err, "not found")
			require.True(t, errors.Is(err, repository.ErrObjectNotFound))
			assert.Nil(t, user)
		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()
			s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,address,contact FROM pvz where id=$1", gomock.Any()).
				Return(assert.AnError)
			// act
			user, err := s.repo.GetByID(ctx, id)
			// assert
			require.EqualError(t, err, "assert.AnError general error for testing")
			assert.Nil(t, user)
		})
	})
}

func Test_Update(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		pvz = &repository.Pvz{
			ID:      1,
			Name:    "Updated Name",
			Address: "Updated Address",
			Contact: "Updated Contact",
		}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()

		s.mockDB.EXPECT().Exec(ctx, `UPDATE pvz SET name = $2, address = $3, contact = $4 WHERE id = $1;`, pvz.ID, pvz.Name, pvz.Address, pvz.Contact).
			Return(nil, nil)
		// act
		err := s.repo.Update(ctx, pvz)
		// assert
		require.NoError(t, err)
	})
}

func Test_DeleteByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockDB.EXPECT().Exec(ctx, `DELETE FROM pvz WHERE id = $1`, id).Return(nil, nil)
		// act
		err := s.repo.DeleteByID(ctx, id)
		// assert
		require.NoError(t, err)
	})
}
