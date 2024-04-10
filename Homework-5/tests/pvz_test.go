package tests

import (
	"context"
	"homework-3/internal/pkg/repository"
	"homework-3/internal/pkg/repository/postrgesql"
	"homework-3/tests/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePvz(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		db.SetUp(t, "pvz")
		defer db.TearDown()
		//arrenge
		repo := postrgesql.NewPvz(db.DB)
		//act
		resp, err := repo.Add(ctx, fixtures.Pvz().Valid().P())
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, int64(1))
	})
}

func TestGetPvz(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		db.SetUp(t, "pvz")
		defer db.TearDown()
		//arrenge
		repo := postrgesql.NewPvz(db.DB)
		respAdd, err := repo.Add(ctx, fixtures.Pvz().Valid().P())
		require.NoError(t, err)
		//act
		resp, err := repo.GetByID(ctx, respAdd)
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, &repository.Pvz{
			ID:      respAdd,
			Name:    "asd",
			Address: "asd",
			Contact: "asd",
		})
	})
}

func TestUpdatePvz(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		db.SetUp(t, "pvz")
		defer db.TearDown()
		//arrenge
		repo := postrgesql.NewPvz(db.DB)
		respAdd, err := repo.Add(ctx, fixtures.Pvz().Valid().P())
		require.NoError(t, err)
		//act
		err = repo.Update(ctx, &repository.Pvz{
			ID:      respAdd,
			Name:    "asdUp",
			Address: "asdUp",
			Contact: "asdUp"})
		resp, _ := repo.GetByID(ctx, respAdd)
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, &repository.Pvz{
			ID:      respAdd,
			Name:    "asdUp",
			Address: "asdUp",
			Contact: "asdUp",
		})
	})
}

func TestDeletePvz(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		db.SetUp(t, "pvz")
		defer db.TearDown()
		//arrenge
		repo := postrgesql.NewPvz(db.DB)
		respAdd, err := repo.Add(ctx, fixtures.Pvz().Valid().P())
		require.NoError(t, err)
		//act
		err = repo.DeleteByID(ctx, respAdd)
		// assert
		require.NoError(t, err)
		assert.NoError(t, err)
	})
}
