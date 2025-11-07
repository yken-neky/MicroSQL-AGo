package services

import (
	"context"
	"database/sql"
)

// QueryExecutor define las operaciones para ejecutar consultas SQL
type QueryExecutor interface {
	// ExecuteQuery ejecuta una consulta SQL y retorna los resultados
	ExecuteQuery(ctx context.Context, db *sql.DB, query string) (*sql.Rows, error)

	// ExecuteNonQuery ejecuta una consulta que no retorna resultados
	ExecuteNonQuery(ctx context.Context, db *sql.DB, query string) (int64, error)

	// Prepare prepara una consulta SQL para su ejecución
	Prepare(ctx context.Context, db *sql.DB, query string) (*sql.Stmt, error)

	// BeginTx inicia una nueva transacción
	BeginTx(ctx context.Context, db *sql.DB) (*sql.Tx, error)

	// ValidateQuery valida la sintaxis de una consulta SQL
	ValidateQuery(query string) error

	// GetQueryType determina el tipo de consulta (SELECT, INSERT, UPDATE, etc)
	GetQueryType(query string) (string, error)

	// ExtractTables extrae las tablas afectadas por la consulta
	ExtractTables(query string) ([]string, error)

	// GetQueryPlan obtiene el plan de ejecución de la consulta
	GetQueryPlan(ctx context.Context, db *sql.DB, query string) (string, error)
}
