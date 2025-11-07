package services

import (
	"context"
	"database/sql"
)

// SQLServerService define las operaciones para interactuar con SQL Server
type SQLServerService interface {
	// Connect establece una conexión a SQL Server y la prueba
	Connect(ctx context.Context, config SQLServerConfig) (*sql.DB, error)

	// ExecuteQuery ejecuta una query y retorna un valor booleano (para controles)
	ExecuteQuery(ctx context.Context, db *sql.DB, query string) (bool, error)

	// ValidateConnection verifica si una conexión está viva
	ValidateConnection(ctx context.Context, db *sql.DB) error

	// Close cierra una conexión y libera recursos
	Close(db *sql.DB) error
}

// SQLServerConfig encapsula la configuración de conexión
type SQLServerConfig struct {
	Driver   string
	Server   string
	Port     string
	User     string
	Password string
	Database string
	Options  map[string]string // Opciones adicionales como encrypt=true
}

// ConnectionError representa errores específicos de conexión
type ConnectionError struct {
	Message string
	Cause   error
}

func (e *ConnectionError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}
