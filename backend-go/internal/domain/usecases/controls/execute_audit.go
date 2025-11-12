package controls

import (
	"context"
	"fmt"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

// ExecuteAuditUseCase ejecuta scripts de control predefinidos (auditorías completas o parciales)
type ExecuteAuditUseCase struct {
	controlRepo repositories.ControlRepository
	sqlService  services.SQLServerService
	queryExec   services.QueryExecutor
	connRepo    repositories.ConnectionRepository
}

func NewExecuteAuditUseCase(
	cr repositories.ControlRepository,
	ss services.SQLServerService,
	qe services.QueryExecutor,
	conn repositories.ConnectionRepository,
) *ExecuteAuditUseCase {
	return &ExecuteAuditUseCase{
		controlRepo: cr,
		sqlService:  ss,
		queryExec:   qe,
		connRepo:    conn,
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
}

// Execute ejecuta una auditoría con controles o scripts indicados
func (uc *ExecuteAuditUseCase) Execute(ctx context.Context, userID uint, req AuditRequest) (*AuditResult, error) {
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
			continue
		}

		ok, err := uc.sqlService.ExecuteQuery(ctx, db, sc.QuerySQL)
		if err != nil {
			sr.Passed = false
			sr.Error = err.Error()
			res.Scripts = append(res.Scripts, sr)
			res.Failed++
			continue
		}

		sr.Passed = ok
		if ok {
			res.Passed++
		} else {
			res.Failed++
		}
		res.Scripts = append(res.Scripts, sr)
	}

	res.Total = len(res.Scripts)

	return res, nil
}
