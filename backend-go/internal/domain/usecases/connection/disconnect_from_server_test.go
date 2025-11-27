package connection

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
    "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/mocks"
)

func TestDisconnectFromServer_DeletesActiveConnection(t *testing.T) {
    // prepare an active connection
    conn := &entities.ActiveConnection{
        ID:            1,
        UserID:        6,
        Manager:       "mssql",
        Driver:        "mssql",
        Server:        "host",
        DBUser:        "sa",
        Password:      "enc",
        IsConnected:   true,
        LastConnected: time.Now(),
    }

    mconn := &mocks.MockConnectionRepository{}
    // expect GetActive to find it
    mconn.On("GetActiveByUserIDAndManager", uint(6), "mssql").Return(conn, nil)
    // expect DeleteActiveByUserAndManager to be called
    mconn.On("DeleteActiveByUserAndManager", uint(6), "mssql").Return(nil)
    // expect log to be called
    mconn.On("LogConnection", mock.Anything).Return(nil)

    msql := &mocks.MockSQLServerService{}

    uc := NewDisconnectFromServerUseCase(mconn, msql)

    err := uc.Execute(context.Background(), 6, "mssql")
    assert.NoError(t, err)

    mconn.AssertExpectations(t)
}

func TestDisconnectFromServer_NoActiveReturnsError(t *testing.T) {
    mconn := &mocks.MockConnectionRepository{}
    mconn.On("GetActiveByUserIDAndManager", uint(6), "mssql").Return(nil, nil)

    msql := &mocks.MockSQLServerService{}

    uc := NewDisconnectFromServerUseCase(mconn, msql)

    err := uc.Execute(context.Background(), 6, "mssql")
    assert.Error(t, err)
    assert.Equal(t, "no active connection found for this driver", err.Error())

    mconn.AssertExpectations(t)
}
