package queries

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

// ExecuteQueryUseCase maneja la ejecución de consultas SQL
type ExecuteQueryUseCase struct {
	queryRepo  repositories.QueryRepository
	sqlService services.SQLServerService
	queryExec  services.QueryExecutor
	connRepo   repositories.ConnectionRepository
}

func NewExecuteQueryUseCase(
	qr repositories.QueryRepository,
	ss services.SQLServerService,
	qe services.QueryExecutor,
	cr repositories.ConnectionRepository,
) *ExecuteQueryUseCase {
	return &ExecuteQueryUseCase{
		queryRepo:  qr,
		sqlService: ss,
		queryExec:  qe,
		connRepo:   cr,
	}
}

// QueryRequest representa la solicitud para ejecutar una consulta
type QueryRequest struct {
	SQL      string
	Database string
	PageSize int
}

// Execute ejecuta una consulta SQL
func (uc *ExecuteQueryUseCase) Execute(ctx context.Context, userID uint, req QueryRequest) (*entities.QueryResult, error) {
	// Verificar conexión activa
	conn, err := uc.connRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}
	if conn == nil || !conn.IsConnected {
		return nil, errors.New("no active connection")
	}

	// Obtener datos de la conexión activa y crear una conexión física
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

	// Crear registro de consulta
	query := &entities.Query{
		UserID:       userID,
		ConnectionID: conn.ID,
		SQL:          req.SQL,
		Database:     req.Database,
		Status:       "running",
		StartTime:    time.Now(),
	}

	if err := uc.queryRepo.Create(query); err != nil {
		return nil, err
	}

	// Validar y determinar tipo de consulta
	if err := uc.queryExec.ValidateQuery(req.SQL); err != nil {
		return uc.handleQueryError(query, err)
	}

	queryType, err := uc.queryExec.GetQueryType(req.SQL)
	if err != nil {
		return uc.handleQueryError(query, err)
	}

	var result *entities.QueryResult

	// Ejecutar según el tipo de consulta
	switch queryType {
	case "SELECT":
		result, err = uc.executeSelect(ctx, db, query, req.PageSize)
	case "INSERT", "UPDATE", "DELETE":
		result, err = uc.executeNonQuery(ctx, db, query)
	default:
		err = errors.New("unsupported query type")
	}

	if err != nil {
		return uc.handleQueryError(query, err)
	}

	// Actualizar estado de la consulta
	query.Status = "completed"
	query.EndTime = time.Now()
	if err := uc.queryRepo.Update(query); err != nil {
		return nil, err
	}

	return result, nil
}

// executeSelect ejecuta una consulta SELECT y procesa sus resultados
func (uc *ExecuteQueryUseCase) executeSelect(
	ctx context.Context,
	db *sql.DB,
	query *entities.Query,
	pageSize int,
) (*entities.QueryResult, error) {
	rows, err := uc.queryExec.ExecuteQuery(ctx, db, query.SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Obtener metadatos de las columnas
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	types := make([]string, len(columnTypes))
	for i, ct := range columnTypes {
		types[i] = ct.DatabaseTypeName()
	}

	// Procesar resultados
	var resultRows [][]any
	count := 0
	hasMore := false

	for rows.Next() {
		if pageSize > 0 && count >= pageSize {
			hasMore = true
			break
		}

		// Crear slice para valores de columnas
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		resultRows = append(resultRows, values)
		count++
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Crear resultado
	result := &entities.QueryResult{
		QueryID:     query.ID,
		Columns:     columns,
		Types:       types,
		Rows:        resultRows,
		HasMoreRows: hasMore,
		PageSize:    pageSize,
		PageNumber:  1,
	}

	// Guardar resultado
	if err := uc.queryRepo.SaveResult(result); err != nil {
		return nil, err
	}

	return result, nil
}

// executeNonQuery ejecuta una consulta que no retorna resultados
func (uc *ExecuteQueryUseCase) executeNonQuery(
	ctx context.Context,
	db *sql.DB,
	query *entities.Query,
) (*entities.QueryResult, error) {
	rowsAffected, err := uc.queryExec.ExecuteNonQuery(ctx, db, query.SQL)
	if err != nil {
		return nil, err
	}

	query.RowsAffected = rowsAffected

	result := &entities.QueryResult{
		QueryID: query.ID,
	}

	return result, nil
}

// handleQueryError maneja errores de ejecución de consultas
func (uc *ExecuteQueryUseCase) handleQueryError(query *entities.Query, err error) (*entities.QueryResult, error) {
	query.Status = "failed"
	query.EndTime = time.Now()
	query.Error = err.Error()

	if updateErr := uc.queryRepo.Update(query); updateErr != nil {
		return nil, errors.New(err.Error() + ": " + updateErr.Error())
	}

	return nil, err
}
