package mocks

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

type MockSQLServerService struct {
	mock.Mock
}

func (m *MockSQLServerService) Connect(ctx context.Context, cfg services.SQLServerConfig) (*sql.DB, error) {
	args := m.Called(ctx, cfg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.DB), args.Error(1)
}

func (m *MockSQLServerService) GetConnection(userID uint) (*sql.DB, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.DB), args.Error(1)
}

func (m *MockSQLServerService) ExecuteQuery(ctx context.Context, db *sql.DB, query string) (bool, error) {
	args := m.Called(ctx, db, query)
	return args.Bool(0), args.Error(1)
}

func (m *MockSQLServerService) ValidateConnection(ctx context.Context, db *sql.DB) error {
	args := m.Called(ctx, db)
	return args.Error(0)
}

func (m *MockSQLServerService) Close(db *sql.DB) error {
	args := m.Called(db)
	return args.Error(0)
}
