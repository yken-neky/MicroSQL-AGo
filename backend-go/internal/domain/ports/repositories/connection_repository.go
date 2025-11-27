package repositories

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
)

// ConnectionRepository maneja las conexiones activas y su historial
type ConnectionRepository interface {
	// Conexiones activas
	CreateActive(conn *entities.ActiveConnection) error
	GetActiveByUserID(userID uint) (*entities.ActiveConnection, error)
	// GetActiveByUserIDAndManager returns active connection for specific user and manager (gestor)
	GetActiveByUserIDAndManager(userID uint, manager string) (*entities.ActiveConnection, error)
	UpdateActive(conn *entities.ActiveConnection) error
	// DeleteActive removes all active connections for a user (used to close all)
	DeleteActive(userID uint) error
	// DeleteActiveByUserAndManager removes a specific active connection for a user and manager
	DeleteActiveByUserAndManager(userID uint, manager string) error
	ListActive() ([]*entities.ActiveConnection, error)
	// ListActiveByUser returns all active connections for a given user across drivers
	ListActiveByUser(userID uint) ([]*entities.ActiveConnection, error)

	// Historial de conexiones
	LogConnection(log *entities.ConnectionLog) error
	GetLogsByUserID(userID uint, limit, offset int) ([]*entities.ConnectionLog, error)
	GetLogByID(id uint) (*entities.ConnectionLog, error)
	CountLogsByUserID(userID uint) (int64, error)
}

// ConnectionLogFilter para filtrar logs de conexión
type ConnectionLogFilter struct {
	UserID    *uint
	StartDate *string
	EndDate   *string
	Status    *string
	Server    *string
	Limit     int
	Offset    int
}
