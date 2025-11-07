package entities

import (
	"time"
)

// Query representa una consulta SQL y su resultado
type Query struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	ConnectionID uint      `json:"connection_id" gorm:"not null"`
	SQL          string    `json:"sql" gorm:"type:text;not null"`
	Status       string    `json:"status" gorm:"not null"` // pending, running, completed, failed
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time,omitempty"`
	RowsAffected int64     `json:"rows_affected"`
	Error        string    `json:"error,omitempty"`
	Database     string    `json:"database" gorm:"not null"` // Base de datos actual
}

// QueryResult representa el resultado de una consulta SELECT
type QueryResult struct {
	QueryID     uint     `json:"query_id"`
	Columns     []string `json:"columns"`
	Types       []string `json:"types"`     // Tipos de datos de las columnas
	Rows        [][]any  `json:"rows"`      // Datos de las filas
	HasMoreRows bool     `json:"has_more"`  // Indica si hay más filas disponibles
	PageSize    int      `json:"page_size"` // Tamaño de página actual
	PageNumber  int      `json:"page"`      // Número de página actual
}

// ExecutionStats representa estadísticas de ejecución de la consulta
type ExecutionStats struct {
	QueryID      uint    `json:"query_id"`
	Duration     float64 `json:"duration_ms"` // Duración en milisegundos
	RowsAffected int64   `json:"rows_affected"`
	CPU          float64 `json:"cpu_time_ms,omitempty"`
	IO           float64 `json:"io_time_ms,omitempty"`
	Memory       int64   `json:"memory_kb,omitempty"`
}
