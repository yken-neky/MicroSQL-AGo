package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
)

type MockConnectionRepository struct {
	mock.Mock
}

// CreateActive mocks adding or updating an active connection
func (m *MockConnectionRepository) CreateActive(conn *entities.ActiveConnection) error {
	args := m.Called(conn)
	return args.Error(0)
}

func (m *MockConnectionRepository) GetActiveByUserID(userID uint) (*entities.ActiveConnection, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ActiveConnection), args.Error(1)
}

func (m *MockConnectionRepository) GetActiveByUserIDAndManager(userID uint, manager string) (*entities.ActiveConnection, error) {
	args := m.Called(userID, manager)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ActiveConnection), args.Error(1)
}

func (m *MockConnectionRepository) UpdateActive(conn *entities.ActiveConnection) error {
	args := m.Called(conn)
	return args.Error(0)
}

func (m *MockConnectionRepository) DeleteActive(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockConnectionRepository) DeleteActiveByUserAndManager(userID uint, manager string) error {
	args := m.Called(userID, manager)
	return args.Error(0)
}

func (m *MockConnectionRepository) ListActive() ([]*entities.ActiveConnection, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ActiveConnection), args.Error(1)
}

func (m *MockConnectionRepository) ListActiveByUser(userID uint) ([]*entities.ActiveConnection, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ActiveConnection), args.Error(1)
}

// LogConnection mocks recording a connection log
func (m *MockConnectionRepository) LogConnection(log *entities.ConnectionLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockConnectionRepository) GetLogsByUserID(userID uint, limit, offset int) ([]*entities.ConnectionLog, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ConnectionLog), args.Error(1)
}

func (m *MockConnectionRepository) GetLogByID(id uint) (*entities.ConnectionLog, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ConnectionLog), args.Error(1)
}

func (m *MockConnectionRepository) CountLogsByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
