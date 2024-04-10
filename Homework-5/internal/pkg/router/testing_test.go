package router

import (
	"homework-3/config"
	mock_repository "homework-3/internal/pkg/repository/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

type pvzRepoFixtures struct {
	ctrl    *gomock.Controller
	srv     Server
	mockPvz *mock_repository.MockPvzRepo
}

func setUp(t *testing.T) pvzRepoFixtures {
	ctrl := gomock.NewController(t)
	mockPvz := mock_repository.NewMockPvzRepo(ctrl)
	srv := Server{mockPvz, config.AuthConfig}
	return pvzRepoFixtures{
		ctrl:    ctrl,
		mockPvz: mockPvz,
		srv:     srv,
	}
}

func (a *pvzRepoFixtures) tearDown() {
	a.ctrl.Finish()
}
