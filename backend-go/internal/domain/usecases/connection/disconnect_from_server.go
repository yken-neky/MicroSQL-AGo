package connection

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

// DisconnectFromServerUseCase maneja la desconexión de SQL Server
type DisconnectFromServerUseCase struct {
	connRepo   repositories.ConnectionRepository
	sqlService services.SQLServerService
}

func NewDisconnectFromServerUseCase(
	cr repositories.ConnectionRepository,
	ss services.SQLServerService,
) *DisconnectFromServerUseCase {
	return &DisconnectFromServerUseCase{
		connRepo:   cr,
		sqlService: ss,
	}
}

// Execute realiza la desconexión de SQL Server
func (uc *DisconnectFromServerUseCase) Execute(ctx context.Context, userID uint) error {
	// Obtener conexión activa
	active, err := uc.connRepo.GetActiveByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get active connection: %w", err)
	}

	if active == nil {
		return errors.New("no active connection found")
	}

	// Actualizar registro de conexión activa (no gestionamos cierre físico del pool aquí)
	active.IsConnected = false
	active.LastDisconnected = time.Now()

	if err := uc.connRepo.UpdateActive(active); err != nil {
		return fmt.Errorf("failed to update connection status: %w", err)
	}

	// Registrar en el historial
	log := &entities.ConnectionLog{
		UserID:    userID,
		Driver:    active.Driver,
		Server:    active.Server,
		DBUser:    active.DBUser,
		Timestamp: time.Now(),
		Status:    "disconnected",
	}

	if err := uc.connRepo.LogConnection(log); err != nil {
		// Solo loggeamos el error, no fallamos la operación
		fmt.Printf("failed to log disconnection: %v\n", err)
	}

	return nil
}
