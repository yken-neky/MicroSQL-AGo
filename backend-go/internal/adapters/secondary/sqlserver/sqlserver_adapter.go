package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"go.uber.org/zap"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

// SQLServerAdapter implementa SQLServerService con pool de conexiones
type SQLServerAdapter struct {
	logger *zap.Logger
	mu     sync.RWMutex
	pools  map[string]*sql.DB // key: DSN
}

func NewSQLServerAdapter(logger *zap.Logger) *SQLServerAdapter {
	return &SQLServerAdapter{
		logger: logger,
		pools:  make(map[string]*sql.DB),
	}
}

func (a *SQLServerAdapter) Connect(ctx context.Context, cfg services.SQLServerConfig) (*sql.DB, error) {
	dsn := buildDSN(cfg)

	// Check if we already have a pool for this DSN
	a.mu.RLock()
	if db, exists := a.pools[dsn]; exists {
		a.mu.RUnlock()
		if err := a.ValidateConnection(ctx, db); err == nil {
			return db, nil
		}
		// Connection is dead, remove it and create new one
		a.mu.Lock()
		delete(a.pools, dsn)
		a.mu.Unlock()
	} else {
		a.mu.RUnlock()
	}

	// Create new connection pool
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, &services.ConnectionError{
			Message: "failed to open connection",
			Cause:   err,
		}
	}

	// Configure pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Validate connection
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, &services.ConnectionError{
			Message: "failed to ping server",
			Cause:   err,
		}
	}

	// Store in pool map
	a.mu.Lock()
	a.pools[dsn] = db
	a.mu.Unlock()

	return db, nil
}

func (a *SQLServerAdapter) ExecuteQuery(ctx context.Context, db *sql.DB, query string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var result bool
	err := db.QueryRowContext(ctx, query).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("query execution failed: %w", err)
	}

	return result, nil
}

func (a *SQLServerAdapter) ValidateConnection(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

func (a *SQLServerAdapter) Close(db *sql.DB) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Find and remove from pools
	for dsn, pool := range a.pools {
		if pool == db {
			delete(a.pools, dsn)
			break
		}
	}

	return db.Close()
}

// buildDSN construye la cadena de conexión para SQL Server
func buildDSN(cfg services.SQLServerConfig) string {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		cfg.User, cfg.Password, cfg.Server, cfg.Port, cfg.Database)

	// Add additional options
	for k, v := range cfg.Options {
		// Some callers may provide yes/no; the driver expects boolean strings parseable by strconv.ParseBool
		// Normalize common synonyms: yes -> true, no -> false
		lower := strings.ToLower(v)
		if lower == "yes" {
			v = "yes"
		} else if lower == "no" {
			v = "no"
		}
		dsn += fmt.Sprintf("&%s=%s", k, v)
	}

	return dsn
}
