package repositories

import (
	"testing"
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormAuditRepository_CreateGet(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&entities.AuditRun{}, &entities.AuditScriptResult{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	repo := NewGormAuditRepository(db)

	run := &entities.AuditRun{
		UserID:   42,
		Mode:     "partial",
		Database: "db1",
		Status:   "running",
	}

	if err := repo.CreateAuditRun(run); err != nil {
		t.Fatalf("create audit run: %v", err)
	}
	if run.ID == 0 {
		t.Fatalf("expected run ID to be set")
	}

	got, err := repo.GetAuditRunByID(run.ID)
	if err != nil {
		t.Fatalf("get audit run: %v", err)
	}
	if got.UserID != run.UserID || got.Database != run.Database {
		t.Fatalf("mismatch run data")
	}

	// create a script result
	res := &entities.AuditScriptResult{
		AuditRunID: run.ID,
		ScriptID:   1,
		ControlID:  2,
		QuerySQL:   "SELECT 1",
		Passed:     true,
		DurationMs: 10,
		Rows:       0,
	}
	if err := repo.CreateScriptResult(res); err != nil {
		t.Fatalf("create script result: %v", err)
	}

	list, err := repo.ListScriptResultsByAuditRun(run.ID)
	if err != nil {
		t.Fatalf("list results: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 result, got %d", len(list))
	}

	// update run as completed
	now := time.Now()
	run.FinishedAt = &now
	run.Status = "completed"
	run.Total = 1
	run.Passed = 1
	_ = repo.UpdateAuditRun(run)

	got2, _ := repo.GetAuditRunByID(run.ID)
	if got2.Status != "completed" {
		t.Fatalf("expected status completed, got %s", got2.Status)
	}
}
