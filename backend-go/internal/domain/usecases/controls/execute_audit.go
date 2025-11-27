package controls

import (
	"context"
	"fmt"
	"strings"
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
	encryptSvc  services.EncryptionService
}

// NewExecuteAuditUseCase crea una nueva instancia con todas las dependencias
func NewExecuteAuditUseCase(
	cr repositories.ControlRepository,
	ss services.SQLServerService,
	qe services.QueryExecutor,
	conn repositories.ConnectionRepository,
	ar repositories.AuditRepository,
	enc services.EncryptionService,
) *ExecuteAuditUseCase {
	return &ExecuteAuditUseCase{
		controlRepo: cr,
		sqlService:  ss,
		queryExec:   qe,
		connRepo:    conn,
		auditRepo:   ar,
		encryptSvc:  enc,
	}
}

// AuditRequest representa la petición para ejecutar una auditoría
type AuditRequest struct {
	ControlIDs []uint `json:"control_ids,omitempty"`
	ScriptIDs  []uint `json:"script_ids,omitempty"`
	Database   string `json:"database"`
	FullAudit  bool   `json:"full_audit,omitempty"`
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
	Total      int            `json:"total"`
	Passed     int            `json:"passed"`
	Manual     int            `json:"manual_count,omitempty"`
	Failed     int            `json:"failed"`
	Scripts    []ScriptResult `json:"scripts"`
	AuditRunID uint           `json:"audit_run_id,omitempty"`
}

// Execute ejecuta una auditoría con controles o scripts indicados
func (uc *ExecuteAuditUseCase) Execute(ctx context.Context, userID uint, manager string, req AuditRequest) (*AuditResult, error) {
	// Prepare and persist AuditRun (mode: partial|full)
	mode := "partial"
	if req.FullAudit {
		mode = "full"
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
	if req.FullAudit {
		run.Controls = "ALL"
	}

	if uc.auditRepo != nil {
		if err := uc.auditRepo.CreateAuditRun(run); err != nil {
			return nil, err
		}
	}
	// Verificar conexión activa — preferir la conexión activa para el gestor/driver pedido
	// Intentar obtener la conexión activa específica por user+driver
	conn, err := uc.connRepo.GetActiveByUserIDAndManager(userID, manager)
	if err != nil {
		return nil, err
	}
	// Si no existe conexión específica, buscar en todas las activas y aplicar heurísticas
	if conn == nil {
		conns, err := uc.connRepo.ListActiveByUser(userID)
		if err != nil {
			return nil, err
		}
		if len(conns) == 0 {
			return nil, fmt.Errorf("no active connection")
		}

		// Preferir conexiones que parezcan apuntar a SQL Server (driver contiene sql/mssql/odbc)
		var candidates []*entities.ActiveConnection
		for _, c := range conns {
			if !c.IsConnected {
				continue
			}
			low := strings.ToLower(c.Driver)
			if strings.Contains(low, "mssql") {
				candidates = append(candidates, c)
			}
		}

		// Si no hay candidatos específicos, usar todos los conectados
		if len(candidates) == 0 {
			for _, c := range conns {
				if c.IsConnected {
					candidates = append(candidates, c)
				}
			}
		}

		// Elegir la más reciente (LastConnected)
		var latestConn *entities.ActiveConnection
		var latest time.Time
		for _, c := range candidates {
			if c.LastConnected.After(latest) {
				latest = c.LastConnected
				latestConn = c
			}
		}
		if latestConn == nil {
			return nil, fmt.Errorf("no active connection")
		}
		// assign selected connection to outer variable
		conn = latestConn
	}
	// Recolectar scripts desde full audit OR controlIDs/scriptIDs
	// If FullAudit==true we ignore control_ids/script_ids and load all scripts
	// from the control repository.
	scriptsMap := make(map[uint]repositories.ControlsScript)

	if req.FullAudit {
		all, err := uc.controlRepo.GetAllScripts()
		if err != nil {
			return nil, err
		}
		for _, sc := range all {
			scriptsMap[sc.ID] = sc
		}
	}

	if !req.FullAudit && len(req.ControlIDs) > 0 {
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

	if !req.FullAudit && len(req.ScriptIDs) > 0 {
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
	// Desencriptar contraseña si está cifrada (fallback: usarla tal cual)
	password := conn.Password
	if uc.encryptSvc != nil {
		if dec, derr := uc.encryptSvc.Decrypt(conn.Password); derr == nil && dec != "" {
			password = dec
		}
	}

	cfg := services.SQLServerConfig{
		Driver:   conn.Driver,
		Server:   conn.Server,
		Port:     "1433",
		User:     conn.DBUser,
		Password: password,
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

		// If script is manual, treat it as passed (manual checks are external)
		if strings.ToLower(strings.TrimSpace(sc.ControlType)) == "manual" {
			sr.Passed = true
			// Count manual as passed
			res.Passed++
			res.Manual++
			res.Scripts = append(res.Scripts, sr)
			if uc.auditRepo != nil {
				resRow := &entities.AuditScriptResult{
					AuditRunID: run.ID,
					ScriptID:   sc.ID,
					ControlID:  sc.ControlScriptRef,
					QuerySQL:   sc.QuerySQL,
					Passed:     true,
				}
				_ = uc.auditRepo.CreateScriptResult(resRow)
			}
			continue
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
		Total:      run.Total,
		Passed:     run.Passed,
		Failed:     run.Failed,
		Scripts:    make([]ScriptResult, 0, len(results)),
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
		if r.QuerySQL == "" {
			res.Manual++
		}
		res.Scripts = append(res.Scripts, sr)
	}

	return res, run, nil
}
