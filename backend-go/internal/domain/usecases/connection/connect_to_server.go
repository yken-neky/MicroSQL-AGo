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

// ConnectToServerUseCase maneja la lógica de conexión a SQL Server
type ConnectToServerUseCase struct {
	connRepo   repositories.ConnectionRepository
	sqlService services.SQLServerService
	encryptSvc services.EncryptionService
}

func NewConnectToServerUseCase(
	cr repositories.ConnectionRepository,
	ss services.SQLServerService,
	es services.EncryptionService,
) *ConnectToServerUseCase {
	return &ConnectToServerUseCase{
		connRepo:   cr,
		sqlService: ss,
		encryptSvc: es,
	}
}

// Execute intenta establecer una conexión a SQL Server
func (uc *ConnectToServerUseCase) Execute(ctx context.Context, userID uint, req ConnectRequest) (*entities.ActiveConnection, error) {
	// Verificar si ya existe una conexión activa
	if active, _ := uc.connRepo.GetActiveByUserID(userID); active != nil {
		return nil, errors.New("user already has an active connection")
	}

	// Encriptar contraseña antes de almacenar
	encryptedPass, err := uc.encryptSvc.Encrypt(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Intentar conectar a SQL Server
	cfg := services.SQLServerConfig{
		Driver:   req.Driver,
		Server:   req.Server,
		Port:     req.Port,
		User:     req.DBUser,
		Password: req.Password,
		Database: "master", // Siempre conectamos a master primero
		Options: map[string]string{
			"encrypt":                "true",
			"TrustServerCertificate": "true",
		},
	}

	db, err := uc.sqlService.Connect(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	// Crear registro de conexión activa
	conn := &entities.ActiveConnection{
		UserID:        userID,
		Driver:        req.Driver,
		Server:        req.Server,
		DBUser:        req.DBUser,
		Password:      encryptedPass,
		IsConnected:   true,
		LastConnected: time.Now(),
	}

	if err := uc.connRepo.CreateActive(conn); err != nil {
		uc.sqlService.Close(db)
		return nil, fmt.Errorf("failed to save connection: %w", err)
	}

	// Registrar en el historial
	log := &entities.ConnectionLog{
		UserID:    userID,
		Driver:    req.Driver,
		Server:    req.Server,
		DBUser:    req.DBUser,
		Timestamp: time.Now(),
		Status:    "connected",
	}

	if err := uc.connRepo.LogConnection(log); err != nil {
		// Solo loggeamos el error, no fallamos la operación
		fmt.Printf("failed to log connection: %v\n", err)
	}

	return conn, nil
}

// ConnectRequest representa los datos necesarios para conectar
type ConnectRequest struct {
	Driver   string
	Server   string
	Port     string
	DBUser   string
	Password string
}
