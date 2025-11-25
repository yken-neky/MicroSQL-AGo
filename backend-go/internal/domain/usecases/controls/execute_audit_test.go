package controls

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	repoport "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/mocks"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeControlRepo struct{}

func (f *fakeControlRepo) ListControls() ([]entities.ControlsInformation, error) { return nil, nil }
func (f *fakeControlRepo) GetControlScripts(controlID uint) ([]repoport.ControlsScript, error) {
	return []repoport.ControlsScript{{ID: 1, ControlScriptRef: controlID, QuerySQL: "SELECT 1", ControlType: "type1"}}, nil
}
func (f *fakeControlRepo) GetScriptsByIDs(ids []uint) ([]repoport.ControlsScript, error) {
	return []repoport.ControlsScript{{ID: 1, ControlScriptRef: 2, QuerySQL: "SELECT 1", ControlType: "type1"}}, nil
}

func TestExecuteAuditPersistsAuditRun(t *testing.T) {
	// in-memory sqlite for audit repo
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&entities.AuditRun{}, &entities.AuditScriptResult{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	auditRepo := repositories.NewGormAuditRepository(db)

	// mocks
	sqlSvc := &mocks.MockSQLServerService{}
	qExec := &mocks.MockQueryExecutor{}
	connRepo := &mocks.MockConnectionRepository{}

	// active connection
	connRepo.On("GetActiveByUserID", uint(1)).Return(&entities.ActiveConnection{UserID: 1, Driver: "sqlserver", Server: "localhost", DBUser: "sa", Password: "x", IsConnected: true}, nil)
	sqlSvc.On("Connect", mock.Anything, mock.Anything).Return(nil, nil)
	sqlSvc.On("ExecuteQuery", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	qExec.On("ValidateQuery", mock.Anything).Return(nil)

	// usecase
	uc := NewExecuteAuditUseCase(&fakeControlRepo{}, sqlSvc, qExec, connRepo, auditRepo)

	req := AuditRequest{ControlIDs: []uint{2}, Database: "db1"}
	res, err := uc.Execute(context.Background(), 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Total)

	// check audit run persisted
	runs := []entities.AuditRun{}
	if err := db.Find(&runs).Error; err != nil {
		t.Fatalf("find runs: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 audit run, got %d", len(runs))
	}
}
