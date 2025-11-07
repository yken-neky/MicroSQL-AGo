package connection

import (
	"context"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
)

// GetActiveConnectionUseCase maneja la obtención de la conexión activa
type GetActiveConnectionUseCase struct {
	connRepo repositories.ConnectionRepository
}

func NewGetActiveConnectionUseCase(
	cr repositories.ConnectionRepository,
) *GetActiveConnectionUseCase {
	return &GetActiveConnectionUseCase{
		connRepo: cr,
	}
}

// Execute obtiene la conexión activa del usuario si existe
func (uc *GetActiveConnectionUseCase) Execute(ctx context.Context, userID uint) (*entities.ActiveConnection, error) {
	return uc.connRepo.GetActiveByUserID(userID)
}
