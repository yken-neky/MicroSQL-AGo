package connection

import (
	"context"
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
)

// ListConnectionHistoryUseCase maneja la obtención del historial de conexiones
type ListConnectionHistoryUseCase struct {
	connRepo repositories.ConnectionRepository
}

func NewListConnectionHistoryUseCase(
	cr repositories.ConnectionRepository,
) *ListConnectionHistoryUseCase {
	return &ListConnectionHistoryUseCase{
		connRepo: cr,
	}
}

// Execute obtiene el historial de conexiones del usuario
func (uc *ListConnectionHistoryUseCase) Execute(ctx context.Context, userID uint, filters HistoryFilters) ([]*entities.ConnectionLog, error) {
	return uc.connRepo.GetLogsByUserID(userID, filters.Limit, filters.Offset)
}

// HistoryFilters contiene los filtros para obtener el historial
type HistoryFilters struct {
	StartDate time.Time
	EndDate   time.Time
	Limit     int
	Offset    int
}
