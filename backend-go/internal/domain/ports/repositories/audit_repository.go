package repositories

import "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"

// AuditRepository provides persistence for audit runs and script results
type AuditRepository interface {
	CreateAuditRun(run *entities.AuditRun) error
	UpdateAuditRun(run *entities.AuditRun) error
	GetAuditRunByID(id uint) (*entities.AuditRun, error)

	CreateScriptResult(res *entities.AuditScriptResult) error
	ListScriptResultsByAuditRun(auditRunID uint) ([]entities.AuditScriptResult, error)
}
