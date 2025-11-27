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
	return &GetActiveConnectionUseCase{connRepo: cr}
}

// Execute obtiene la conexión activa del usuario para el driver/gestor indicado
func (uc *GetActiveConnectionUseCase) Execute(ctx context.Context, userID uint, manager string) (*entities.ActiveConnection, error) {
	return uc.connRepo.GetActiveByUserIDAndManager(userID, manager)
}

// ListActiveConnectionsUseCase lista todas las conexiones activas de un usuario (todos los gestores)
type ListActiveConnectionsUseCase struct {
	connRepo repositories.ConnectionRepository
}

func NewListActiveConnectionsUseCase(cr repositories.ConnectionRepository) *ListActiveConnectionsUseCase {
	return &ListActiveConnectionsUseCase{connRepo: cr}
}

func (uc *ListActiveConnectionsUseCase) Execute(ctx context.Context, userID uint) ([]*entities.ActiveConnection, error) {
	return uc.connRepo.ListActiveByUser(userID)
}
