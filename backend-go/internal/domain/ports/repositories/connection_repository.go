package repositories

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
)

// ConnectionRepository maneja las conexiones activas y su historial
type ConnectionRepository interface {
	// Conexiones activas
	CreateActive(conn *entities.ActiveConnection) error
	GetActiveByUserID(userID uint) (*entities.ActiveConnection, error)
	UpdateActive(conn *entities.ActiveConnection) error
	DeleteActive(userID uint) error
	ListActive() ([]*entities.ActiveConnection, error)

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
