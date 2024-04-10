package router

import (
	"context"
	"homework-3/internal/pkg/repository"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Create(t *testing.T) {
	t.Parallel()
	var (
		ctx   = context.Background()
		repos = &repository.Pvz{
			Name:    "test",
			Address: "test",
			Contact: "test",
		}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()

		s.mockPvz.EXPECT().Add(gomock.Any(), repos).Return(int64(1), nil)
		// act
		result, status := s.srv.crt(ctx, repos)
		// assert
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"name\":\"test\",\"address\":\"test\",\"contact\":\"test\"}", string(result))
	})
}

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

		s.mockPvz.EXPECT().GetByID(gomock.Any(), id).Return(&repository.Pvz{
			ID:      1,
			Name:    "test",
			Address: "test",
			Contact: "test",
		}, nil)
		// act
		result, status := s.srv.get(ctx, id)
		// assert
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"Name\":\"test\",\"Address\":\"test\",\"Contact\":\"test\"}", string(result))
	})
}

func Test_Update(t *testing.T) {
	t.Parallel()
	var (
		ctx   = context.Background()
		repos = &repository.Pvz{
			ID:      1,
			Name:    "test",
			Address: "test",
			Contact: "test",
		}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()

		s.mockPvz.EXPECT().Update(gomock.Any(), repos).Return(nil)
		// act
		result, status := s.srv.upd(ctx, repos)
		// assert
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"Name\":\"test\",\"Address\":\"test\",\"Contact\":\"test\"}", string(result))
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
		s := setUp(t)
		defer s.tearDown()

		s.mockPvz.EXPECT().DeleteByID(gomock.Any(), id).Return(nil)
		status := s.srv.del(ctx, id)

		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, int(200), status)
	})
}
