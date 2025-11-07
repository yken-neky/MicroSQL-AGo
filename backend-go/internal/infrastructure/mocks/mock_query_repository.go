package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
)

type MockQueryRepository struct {
	mock.Mock
}

func (m *MockQueryRepository) Create(q *entities.Query) error {
	args := m.Called(q)
	return args.Error(0)
}

func (m *MockQueryRepository) Update(q *entities.Query) error {
	args := m.Called(q)
	return args.Error(0)
}

func (m *MockQueryRepository) GetByID(id uint) (*entities.Query, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Query), args.Error(1)
}

func (m *MockQueryRepository) ListByUser(userID uint, page, pageSize int) ([]entities.Query, error) {
	args := m.Called(userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Query), args.Error(1)
}

func (m *MockQueryRepository) SaveResult(r *entities.QueryResult) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockQueryRepository) GetResult(queryID uint, page, pageSize int) (*entities.QueryResult, error) {
	args := m.Called(queryID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.QueryResult), args.Error(1)
}

func (m *MockQueryRepository) SaveStats(s *entities.ExecutionStats) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockQueryRepository) GetStats(queryID uint) (*entities.ExecutionStats, error) {
	args := m.Called(queryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ExecutionStats), args.Error(1)
}

func (m *MockQueryRepository) GetUserQueryCount(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueryRepository) CleanOldQueries(olderThan string) error {
	args := m.Called(olderThan)
	return args.Error(0)
}
