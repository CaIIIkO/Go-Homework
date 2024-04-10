// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	pgconn "github.com/jackc/pgconn"
	v4 "github.com/jackc/pgx/v4"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	gomock "go.uber.org/mock/gomock"
)

// MockDBops is a mock of DBops interface.
type MockDBops struct {
	ctrl     *gomock.Controller
	recorder *MockDBopsMockRecorder
}

// MockDBopsMockRecorder is the mock recorder for MockDBops.
type MockDBopsMockRecorder struct {
	mock *MockDBops
}

// NewMockDBops creates a new mock instance.
func NewMockDBops(ctrl *gomock.Controller) *MockDBops {
	mock := &MockDBops{ctrl: ctrl}
	mock.recorder = &MockDBopsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBops) EXPECT() *MockDBopsMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockDBops) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockDBopsMockRecorder) Exec(ctx, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockDBops)(nil).Exec), varargs...)
}

// ExecQueryRow mocks base method.
func (m *MockDBops) ExecQueryRow(ctx context.Context, query string, args ...any) v4.Row {
	m.ctrl.T.Helper()
	varargs := []any{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecQueryRow", varargs...)
	ret0, _ := ret[0].(v4.Row)
	return ret0
}

// ExecQueryRow indicates an expected call of ExecQueryRow.
func (mr *MockDBopsMockRecorder) ExecQueryRow(ctx, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecQueryRow", reflect.TypeOf((*MockDBops)(nil).ExecQueryRow), varargs...)
}

// Get mocks base method.
func (m *MockDBops) Get(ctx context.Context, dest any, query string, args ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockDBopsMockRecorder) Get(ctx, dest, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDBops)(nil).Get), varargs...)
}

// GetPool mocks base method.
func (m *MockDBops) GetPool(arg0 context.Context) *pgxpool.Pool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPool", arg0)
	ret0, _ := ret[0].(*pgxpool.Pool)
	return ret0
}

// GetPool indicates an expected call of GetPool.
func (mr *MockDBopsMockRecorder) GetPool(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPool", reflect.TypeOf((*MockDBops)(nil).GetPool), arg0)
}

// Select mocks base method.
func (m *MockDBops) Select(ctx context.Context, dest any, query string, args ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Select", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Select indicates an expected call of Select.
func (mr *MockDBopsMockRecorder) Select(ctx, dest, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockDBops)(nil).Select), varargs...)
}
