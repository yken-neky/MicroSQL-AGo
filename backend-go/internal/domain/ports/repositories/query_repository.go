package repositories

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
)

// QueryRepository define las operaciones para gestionar consultas
type QueryRepository interface {
	// Crear una nueva consulta
	Create(query *entities.Query) error

	// Actualizar una consulta existente
	Update(query *entities.Query) error

	// Obtener una consulta por ID
	GetByID(id uint) (*entities.Query, error)

	// Listar consultas por usuario con paginación
	ListByUser(userID uint, page, pageSize int) ([]entities.Query, error)

	// Guardar resultado de consulta SELECT
	SaveResult(result *entities.QueryResult) error

	// Obtener resultado de consulta por ID
	GetResult(queryID uint, page, pageSize int) (*entities.QueryResult, error)

	// Guardar estadísticas de ejecución
	SaveStats(stats *entities.ExecutionStats) error

	// Obtener estadísticas de ejecución por ID de consulta
	GetStats(queryID uint) (*entities.ExecutionStats, error)

	// Obtener total de consultas por usuario
	GetUserQueryCount(userID uint) (int64, error)

	// Eliminar consultas antiguas
	CleanOldQueries(olderThan string) error
}
