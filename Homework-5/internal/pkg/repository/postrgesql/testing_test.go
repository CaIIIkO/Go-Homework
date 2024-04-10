package postrgesql

import (
	mock_db "homework-3/internal/pkg/db/mocks"
	"homework-3/internal/pkg/repository"
	"testing"

	"go.uber.org/mock/gomock"
)

type pvzRepoFixtures struct {
	ctrl   *gomock.Controller
	repo   repository.PvzRepo
	mockDB *mock_db.MockDBops
}

func setUp(t *testing.T) pvzRepoFixtures {
	ctrl := gomock.NewController(t)
	mockDB := mock_db.NewMockDBops(ctrl)
	repo := NewPvz(mockDB)
	return pvzRepoFixtures{
		ctrl:   ctrl,
		repo:   repo,
		mockDB: mockDB,
	}
}

func (a *pvzRepoFixtures) tearDown() {
	a.ctrl.Finish()
}
