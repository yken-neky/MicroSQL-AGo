package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// SQLServerQueryExecutor implementa QueryExecutor para SQL Server
type SQLServerQueryExecutor struct{}

func NewSQLServerQueryExecutor() *SQLServerQueryExecutor {
	return &SQLServerQueryExecutor{}
}

// ExecuteQuery ejecuta una consulta SQL que retorna resultados
func (e *SQLServerQueryExecutor) ExecuteQuery(ctx context.Context, db *sql.DB, query string) (*sql.Rows, error) {
	return db.QueryContext(ctx, query)
}

// ExecuteNonQuery ejecuta una consulta que no retorna resultados
func (e *SQLServerQueryExecutor) ExecuteNonQuery(ctx context.Context, db *sql.DB, query string) (int64, error) {
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Prepare prepara una consulta SQL
func (e *SQLServerQueryExecutor) Prepare(ctx context.Context, db *sql.DB, query string) (*sql.Stmt, error) {
	return db.PrepareContext(ctx, query)
}

// BeginTx inicia una transacción
func (e *SQLServerQueryExecutor) BeginTx(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	return db.BeginTx(ctx, nil)
}

// ValidateQuery realiza validación básica de la sintaxis SQL
func (e *SQLServerQueryExecutor) ValidateQuery(query string) error {
	query = strings.TrimSpace(query)
	if query == "" {
		return errors.New("empty query")
	}

	// Verificar palabras clave peligrosas
	// Use word boundaries to reduce false positives (e.g. don't match 'backupset')
	dangerousKeywords := []string{
		`(?i)\bSHUTDOWN\b`,
		`(?i)\bBACKUP\b`,
		`(?i)\bRESTORE\b`,
		`(?i)\bKILL\b`,
		`(?i)\bDROP\s+DATABASE\b`,
		`(?i)\bALTER\s+DATABASE\b`,
		`(?i)\bCREATE\s+DATABASE\b`,
		`(?i)\bxp_cmdshell\b`,
		`(?i)\bsp_configure\b`,
	}

	for _, keyword := range dangerousKeywords {
		if matched, _ := regexp.MatchString(keyword, query); matched {
			// Return a clearer error including the matched keyword for diagnostics
			return fmt.Errorf("query contains forbidden keyword: %s", keyword)
		}
	}

	return nil
}

// GetQueryType determina el tipo de consulta SQL
func (e *SQLServerQueryExecutor) GetQueryType(query string) (string, error) {
	query = strings.TrimSpace(strings.ToUpper(query))

	patterns := map[string]string{
		"SELECT": `^SELECT\s`,
		"INSERT": `^INSERT\s`,
		"UPDATE": `^UPDATE\s`,
		"DELETE": `^DELETE\s`,
		"MERGE":  `^MERGE\s`,
	}

	for queryType, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, query); matched {
			return queryType, nil
		}
	}

	return "", errors.New("unknown query type")
}

// ExtractTables extrae las tablas mencionadas en la consulta
func (e *SQLServerQueryExecutor) ExtractTables(query string) ([]string, error) {
	// Expresión regular para encontrar nombres de tablas
	// Nota: Esta es una implementación básica y puede necesitar mejoras
	tablePattern := `(?i)FROM\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)?)`
	re := regexp.MustCompile(tablePattern)

	matches := re.FindAllStringSubmatch(query, -1)
	if matches == nil {
		return nil, nil
	}

	tables := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			tables = append(tables, match[1])
		}
	}

	return tables, nil
}

// GetQueryPlan obtiene el plan de ejecución de la consulta
func (e *SQLServerQueryExecutor) GetQueryPlan(ctx context.Context, db *sql.DB, query string) (string, error) {
	// Construir consulta para obtener el plan de ejecución
	planQuery := "SET SHOWPLAN_XML ON; " + query + "; SET SHOWPLAN_XML OFF;"

	rows, err := db.QueryContext(ctx, planQuery)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", errors.New("no execution plan available")
	}

	var plan string
	if err := rows.Scan(&plan); err != nil {
		return "", err
	}

	return plan, nil
}
