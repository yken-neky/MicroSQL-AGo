package controls

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/encryption"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/mocks"
)

// fake control repo returns a single script for ids
type fakeControlRepo struct{}

func (f *fakeControlRepo) ListControls() ([]entities.ControlsInformation, error) { return nil, nil }
func (f *fakeControlRepo) GetControlScripts(controlID uint) ([]repositories.ControlsScript, error) {
	return nil, nil
}
func (f *fakeControlRepo) GetScriptsByIDs(ids []uint) ([]repositories.ControlsScript, error) {
	return []repositories.ControlsScript{{ID: 1, ControlType: "automatic", QuerySQL: "SELECT 1", ControlScriptRef: 1}}, nil
}
func (f *fakeControlRepo) GetAllScripts() ([]repositories.ControlsScript, error) {
	return []repositories.ControlsScript{{ID: 1, ControlType: "automatic", QuerySQL: "SELECT 1", ControlScriptRef: 1}}, nil
}

// manualRepo returns manual and automatic scripts for tests
type manualRepo struct{}

func (r *manualRepo) ListControls() ([]entities.ControlsInformation, error) { return nil, nil }
func (r *manualRepo) GetControlScripts(controlID uint) ([]repositories.ControlsScript, error) {
	if controlID == 100 { // manual control
		return []repositories.ControlsScript{{ID: 100, ControlType: "manual", QuerySQL: "", ControlScriptRef: 100}}, nil
	}
	return []repositories.ControlsScript{{ID: 1, ControlType: "automatic", QuerySQL: "SELECT 1", ControlScriptRef: 1}}, nil
}
func (r *manualRepo) GetScriptsByIDs(ids []uint) ([]repositories.ControlsScript, error) {
	// return the automatic script for any id
	return []repositories.ControlsScript{{ID: 1, ControlType: "automatic", QuerySQL: "SELECT 1", ControlScriptRef: 1}}, nil
}
func (r *manualRepo) GetAllScripts() ([]repositories.ControlsScript, error) { return nil, nil }

// fake audit repo minimal implementation collecting results
type fakeAuditRepo struct {
	createdRun *entities.AuditRun
}

func (f *fakeAuditRepo) CreateAuditRun(run *entities.AuditRun) error { f.createdRun = run; return nil }
func (f *fakeAuditRepo) UpdateAuditRun(run *entities.AuditRun) error { f.createdRun = run; return nil }
func (f *fakeAuditRepo) GetAuditRunByID(id uint) (*entities.AuditRun, error) {
	return f.createdRun, nil
}
func (f *fakeAuditRepo) CreateScriptResult(res *entities.AuditScriptResult) error { return nil }
func (f *fakeAuditRepo) ListScriptResultsByAuditRun(auditRunID uint) ([]entities.AuditScriptResult, error) {
	return nil, nil
}

func TestExecuteAudit_usesDecryptedPasswordAndLatestConnection(t *testing.T) {
	// prepare encryption and encrypt password
	enc := encryption.NewAESGCMService("your-32-byte-encryption-key-here")
	encrypted, err := enc.Encrypt("STRONG.PASS123sql")
	assert.NoError(t, err)

	// prepare active connection for user 6
	conn := &entities.ActiveConnection{
		ID:            1,
		UserID:        6,
		Driver:        "mssql",
		Server:        "host.docker.internal",
		DBUser:        "sa",
		Password:      encrypted,
		IsConnected:   true,
		LastConnected: time.Now(),
	}

	// mocks
	mconn := &mocks.MockConnectionRepository{}
	mconn.On("GetActiveByUserIDAndManager", uint(6), "mssql").Return(nil, nil)
	mconn.On("ListActiveByUser", uint(6)).Return([]*entities.ActiveConnection{conn}, nil)

	msql := &mocks.MockSQLServerService{}
	// Connect should be called with decrypted password
	msql.On("Connect", mock.Anything, mock.MatchedBy(func(cfg services.SQLServerConfig) bool {
		return cfg.Password == "STRONG.PASS123sql"
	})).Return((*sql.DB)(nil), nil)
	msql.On("ExecuteQuery", mock.Anything, (*sql.DB)(nil), "SELECT 1").Return(true, nil)

	mq := &mocks.MockQueryExecutor{}
	mq.On("ValidateQuery", "SELECT 1").Return(nil)

	aud := &fakeAuditRepo{}

	uc := NewExecuteAuditUseCase(&fakeControlRepo{}, msql, mq, mconn, aud, enc)

	res, err := uc.Execute(context.Background(), 6, "mssql", AuditRequest{ScriptIDs: []uint{1}, Database: "master"})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Total)
	assert.Equal(t, 1, res.Passed)

	mconn.AssertExpectations(t)
	msql.AssertExpectations(t)
	mq.AssertExpectations(t)
}

func TestExecuteAudit_fullAudit_executesAllScripts(t *testing.T) {
	// prepare encryption
	enc := encryption.NewAESGCMService("your-32-byte-encryption-key-here")

	// prepare active connection for user 6
	conn := &entities.ActiveConnection{
		ID:            1,
		UserID:        6,
		Driver:        "mssql",
		Server:        "host.docker.internal",
		DBUser:        "sa",
		Password:      "plain",
		IsConnected:   true,
		LastConnected: time.Now(),
	}

	mconn := &mocks.MockConnectionRepository{}
	mconn.On("GetActiveByUserIDAndManager", uint(6), "mssql").Return(conn, nil)

	msql := &mocks.MockSQLServerService{}
	msql.On("Connect", mock.Anything, mock.Anything).Return((*sql.DB)(nil), nil)
	msql.On("ExecuteQuery", mock.Anything, (*sql.DB)(nil), "SELECT 1").Return(true, nil)

	mq := &mocks.MockQueryExecutor{}
	mq.On("ValidateQuery", "SELECT 1").Return(nil)

	aud := &fakeAuditRepo{}

	uc := NewExecuteAuditUseCase(&fakeControlRepo{}, msql, mq, mconn, aud, enc)

	res, err := uc.Execute(context.Background(), 6, "mssql", AuditRequest{FullAudit: true, Database: "master"})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Total)
	assert.Equal(t, 1, res.Passed)

	// audit run persisted should be full and controls set to ALL
	if aud.createdRun == nil {
		t.Fatalf("expected audit run to be created")
	}
	assert.Equal(t, "full", aud.createdRun.Mode)
	assert.Equal(t, "ALL", aud.createdRun.Controls)

	mconn.AssertExpectations(t)
	msql.AssertExpectations(t)
	mq.AssertExpectations(t)
}

func TestExecuteAudit_manualScripts_areMarkedPassed(t *testing.T) {
	// prepare encryption
	enc := encryption.NewAESGCMService("your-32-byte-encryption-key-here")

	// prepare active connection for user 6
	conn := &entities.ActiveConnection{
		ID:            1,
		UserID:        6,
		Driver:        "mssql",
		Server:        "host.docker.internal",
		DBUser:        "sa",
		Password:      "plain",
		IsConnected:   true,
		LastConnected: time.Now(),
	}

	mconn := &mocks.MockConnectionRepository{}
	mconn.On("GetActiveByUserIDAndManager", uint(6), "mssql").Return(conn, nil)

	// make control repo return one manual and one automatic (see manualRepo implementation)

	msql := &mocks.MockSQLServerService{}
	msql.On("Connect", mock.Anything, mock.Anything).Return((*sql.DB)(nil), nil)
	msql.On("ExecuteQuery", mock.Anything, (*sql.DB)(nil), "SELECT 1").Return(true, nil)

	mq := &mocks.MockQueryExecutor{}
	mq.On("ValidateQuery", "SELECT 1").Return(nil)

	aud := &fakeAuditRepo{}

	uc := NewExecuteAuditUseCase(&manualRepo{}, msql, mq, mconn, aud, enc)

	// Request with both control IDs (one manual) and script IDs
	req := AuditRequest{ControlIDs: []uint{100}, ScriptIDs: []uint{1}, Database: "master"}

	res, err := uc.Execute(context.Background(), 6, "mssql", req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	// we expect 2 scripts total (manual + automatic)
	assert.Equal(t, 2, res.Total)
	// manual counts as passed
	assert.Equal(t, 2, res.Passed)
	assert.Equal(t, 1, res.Manual)
}
