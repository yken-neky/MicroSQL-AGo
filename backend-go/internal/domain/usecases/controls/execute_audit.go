package controls

import (
	"context"
	"fmt"
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

// ExecuteAuditUseCase ejecuta scripts de control predefinidos (auditorías completas o parciales)
type ExecuteAuditUseCase struct {
	controlRepo repositories.ControlRepository
	sqlService  services.SQLServerService
	queryExec   services.QueryExecutor
	connRepo    repositories.ConnectionRepository
	auditRepo   repositories.AuditRepository
}

func NewExecuteAuditUseCase(
	cr repositories.ControlRepository,
	ss services.SQLServerService,
	qe services.QueryExecutor,
	conn repositories.ConnectionRepository,
	ar repositories.AuditRepository,
) *ExecuteAuditUseCase {
	return &ExecuteAuditUseCase{
		controlRepo: cr,
		sqlService:  ss,
		queryExec:   qe,
		connRepo:    conn,
		auditRepo:   ar,
	}
}

// AuditRequest representa la petición para ejecutar una auditoría
type AuditRequest struct {
	ControlIDs []uint `json:"control_ids,omitempty"`
	ScriptIDs  []uint `json:"script_ids,omitempty"`
	Database   string `json:"database"`
}

// ScriptResult es el resultado de ejecutar un script de control
type ScriptResult struct {
	ScriptID    uint   `json:"script_id"`
	ControlID   uint   `json:"control_id"`
	ControlType string `json:"control_type"`
	QuerySQL    string `json:"query_sql"`
	Passed      bool   `json:"passed"`
	Error       string `json:"error,omitempty"`
}

// AuditResult agrega el resumen de la auditoría
type AuditResult struct {
	Total   int            `json:"total"`
	Passed  int            `json:"passed"`
	Failed  int            `json:"failed"`
	Scripts []ScriptResult `json:"scripts"`
	AuditRunID uint        `json:"audit_run_id,omitempty"`
}

// Execute ejecuta una auditoría con controles o scripts indicados
func (uc *ExecuteAuditUseCase) Execute(ctx context.Context, userID uint, req AuditRequest) (*AuditResult, error) {
	// Prepare and persist AuditRun
	mode := "partial"
	if len(req.ControlIDs) == 0 && len(req.ScriptIDs) > 0 {
		mode = "partial"
	}
	// if user provided control IDs but empty scriptIDs, still partial; in future full mode could be explicit
	run := &entities.AuditRun{
		UserID:   userID,
		Mode:     mode,
		Database: req.Database,
		Status:   "running",
		Controls: "",
	}
	// store control IDs as a simple CSV or JSON string
	if len(req.ControlIDs) > 0 {
		// use fmt to build simple comma-separated list
		run.Controls = fmt.Sprintf("%v", req.ControlIDs)
	}

	if uc.auditRepo != nil {
		if err := uc.auditRepo.CreateAuditRun(run); err != nil {
			return nil, err
		}
	}
	// Verificar conexión activa
	conn, err := uc.connRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}
	if conn == nil || !conn.IsConnected {
		return nil, fmt.Errorf("no active connection")
	}

	// Recolectar scripts desde controlIDs y scriptIDs
	scriptsMap := make(map[uint]repositories.ControlsScript)

	if len(req.ControlIDs) > 0 {
		for _, cid := range req.ControlIDs {
			s, err := uc.controlRepo.GetControlScripts(cid)
			if err != nil {
				return nil, err
			}
			for _, sc := range s {
				scriptsMap[sc.ID] = sc
			}
		}
	}

	if len(req.ScriptIDs) > 0 {
		s, err := uc.controlRepo.GetScriptsByIDs(req.ScriptIDs)
		if err != nil {
			return nil, err
		}
		for _, sc := range s {
			scriptsMap[sc.ID] = sc
		}
	}

	// Si no hubo scripts obtenidos, devolver error
	if len(scriptsMap) == 0 {
		return nil, fmt.Errorf("no scripts found for given control_ids or script_ids")
	}

	// Conectar a SQL Server usando la conexión activa
	cfg := services.SQLServerConfig{
		Driver:   conn.Driver,
		Server:   conn.Server,
		Port:     "1433",
		User:     conn.DBUser,
		Password: conn.Password,
		Database: req.Database,
		Options:  map[string]string{"TrustServerCertificate": "true"},
	}

	db, err := uc.sqlService.Connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Ejecutar scripts y recolectar resultados
	res := &AuditResult{}
	for _, sc := range scriptsMap {
		sr := ScriptResult{
			ScriptID:    sc.ID,
			ControlID:   sc.ControlScriptRef,
			ControlType: sc.ControlType,
			QuerySQL:    sc.QuerySQL,
		}

		// Validar script
		if err := uc.queryExec.ValidateQuery(sc.QuerySQL); err != nil {
			sr.Passed = false
			sr.Error = err.Error()
			res.Scripts = append(res.Scripts, sr)
			res.Failed++
			// persist per-script result
			if uc.auditRepo != nil {
				resRow := &entities.AuditScriptResult{
					AuditRunID: run.ID,
					ScriptID:   sc.ID,
					ControlID:  sc.ControlScriptRef,
					QuerySQL:   sc.QuerySQL,
					Passed:     false,
					Error:      sr.Error,
				}
				_ = uc.auditRepo.CreateScriptResult(resRow)
			}
			continue
		}

		// execute and measure
		start := time.Now()
		ok, err := uc.sqlService.ExecuteQuery(ctx, db, sc.QuerySQL)
		duration := time.Since(start)
		if err != nil {
			sr.Passed = false
			sr.Error = err.Error()
			res.Scripts = append(res.Scripts, sr)
			res.Failed++
			if uc.auditRepo != nil {
				resRow := &entities.AuditScriptResult{
					AuditRunID: run.ID,
					ScriptID:   sc.ID,
					ControlID:  sc.ControlScriptRef,
					QuerySQL:   sc.QuerySQL,
					Passed:     false,
					Error:      sr.Error,
					DurationMs: duration.Milliseconds(),
				}
				_ = uc.auditRepo.CreateScriptResult(resRow)
			}
			continue
		}

		sr.Passed = ok
		if ok {
			res.Passed++
		} else {
			res.Failed++
		}
		res.Scripts = append(res.Scripts, sr)
		if uc.auditRepo != nil {
			resRow := &entities.AuditScriptResult{
				AuditRunID: run.ID,
				ScriptID:   sc.ID,
				ControlID:  sc.ControlScriptRef,
				QuerySQL:   sc.QuerySQL,
				Passed:     sr.Passed,
				DurationMs: duration.Milliseconds(),
			}
			_ = uc.auditRepo.CreateScriptResult(resRow)
		}
	}

	res.Total = len(res.Scripts)
	// Finalize audit run
	if uc.auditRepo != nil {
		run.Total = res.Total
		run.Passed = res.Passed
		run.Failed = res.Failed
		now := time.Now()
		run.FinishedAt = &now
		run.Status = "completed"
		_ = uc.auditRepo.UpdateAuditRun(run)
		res.AuditRunID = run.ID
	}

	return res, nil
}

// GetAuditRun fetches an audit run and its script results if the user has access
func (uc *ExecuteAuditUseCase) GetAuditRun(ctx context.Context, userID uint, auditID uint) (*AuditResult, *entities.AuditRun, error) {
	if uc.auditRepo == nil {
		return nil, nil, fmt.Errorf("audit repository not configured")
	}

	run, err := uc.auditRepo.GetAuditRunByID(auditID)
	if err != nil {
		return nil, nil, err
	}

	// ensure the requesting user is owner (simple authorization)
	if run.UserID != userID {
		return nil, nil, fmt.Errorf("forbidden")
	}

	// load script results
	results, err := uc.auditRepo.ListScriptResultsByAuditRun(run.ID)
	if err != nil {
		return nil, nil, err
	}

	res := &AuditResult{
		Total:   run.Total,
		Passed:  run.Passed,
		Failed:  run.Failed,
		Scripts: make([]ScriptResult, 0, len(results)),
		AuditRunID: run.ID,
	}

	for _, r := range results {
		sr := ScriptResult{
			ScriptID:    r.ScriptID,
			ControlID:   r.ControlID,
			ControlType: "", // not persisted here (could be fetched from controlRepo)
			QuerySQL:    r.QuerySQL,
			Passed:      r.Passed,
			Error:       r.Error,
		}
		res.Scripts = append(res.Scripts, sr)
	}

	return res, run, nil
}
