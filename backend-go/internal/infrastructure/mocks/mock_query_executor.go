package mocks

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockQueryExecutor struct {
	mock.Mock
}

func (m *MockQueryExecutor) ExecuteQuery(ctx context.Context, db *sql.DB, query string) (*sql.Rows, error) {
	args := m.Called(ctx, db, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Rows), args.Error(1)
}

func (m *MockQueryExecutor) ExecuteNonQuery(ctx context.Context, db *sql.DB, query string) (int64, error) {
	args := m.Called(ctx, db, query)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueryExecutor) Prepare(ctx context.Context, db *sql.DB, query string) (*sql.Stmt, error) {
	args := m.Called(ctx, db, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Stmt), args.Error(1)
}

func (m *MockQueryExecutor) BeginTx(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	args := m.Called(ctx, db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockQueryExecutor) ValidateQuery(query string) error {
	args := m.Called(query)
	return args.Error(0)
}

func (m *MockQueryExecutor) GetQueryType(query string) (string, error) {
	args := m.Called(query)
	return args.String(0), args.Error(1)
}

func (m *MockQueryExecutor) ExtractTables(query string) ([]string, error) {
	args := m.Called(query)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockQueryExecutor) GetQueryPlan(ctx context.Context, db *sql.DB, query string) (string, error) {
	args := m.Called(ctx, db, query)
	return args.String(0), args.Error(1)
}
