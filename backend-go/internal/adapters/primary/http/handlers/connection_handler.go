package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	dto "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/primary/http/dto"
	connectionuc "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/usecases/connection"
)

// ConnectionHandler maneja endpoints para conectar / desconectar a SQL Server
type ConnectionHandler struct {
	connectUC     *connectionuc.ConnectToServerUseCase
	disconnectUC  *connectionuc.DisconnectFromServerUseCase
	getActiveUC   *connectionuc.GetActiveConnectionUseCase
	listActiveUC  *connectionuc.ListActiveConnectionsUseCase
	listHistoryUC *connectionuc.ListConnectionHistoryUseCase
	logger        *zap.Logger
}

func NewConnectionHandler(cu *connectionuc.ConnectToServerUseCase, du *connectionuc.DisconnectFromServerUseCase, gu *connectionuc.GetActiveConnectionUseCase, lu *connectionuc.ListActiveConnectionsUseCase, lhu *connectionuc.ListConnectionHistoryUseCase, logger *zap.Logger) *ConnectionHandler {
	return &ConnectionHandler{connectUC: cu, disconnectUC: du, getActiveUC: gu, listActiveUC: lu, listHistoryUC: lhu, logger: logger}
}

// Connect crea una conexión activa para el usuario y la persiste
func (h *ConnectionHandler) Connect(c *gin.Context) {
	var req dto.ConnectRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, _ := c.Get("userID")
	userID := uid.(uint)

	// manager/driver is taken from path param
	manager := c.Param("manager")
	// validate supported managers
	allowed := map[string]bool{"pgsql": true, "oracle": true, "mysql": true, "mssql": true, "otro": true}
	if manager == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing manager in path"})
		return
	}
	if !allowed[manager] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported manager"})
		return
	}
	input := connectionuc.ConnectRequest{
		Manager:  req.Manager,
		Driver:   req.Driver,
		Server:   req.Server,
		Port:     req.Port,
		DBUser:   req.DBUser,
		Password: req.Password,
	}

	// override driver with manager from path
	// input.Driver = manager

	conn, err := h.connectUC.Execute(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ConnectionResponseDTO{
		ID:            conn.ID,
		UserID:        conn.UserID,
		Manager:       conn.Manager,
		Driver:        conn.Driver,
		Server:        conn.Server,
		DBUser:        conn.DBUser,
		IsConnected:   conn.IsConnected,
		LastConnected: conn.LastConnected.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, gin.H{"connection": resp})
}

// Disconnect cierra la conexión activa del usuario
func (h *ConnectionHandler) Disconnect(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint)
	manager := c.Param("manager")
	allowed := map[string]bool{"pgsql": true, "oracle": true, "mysql": true, "mssql": true, "otro": true}
	if manager == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing manager in path"})
		return
	}
	if !allowed[manager] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported manager"})
		return
	}

	if err := h.disconnectUC.Execute(c.Request.Context(), userID, manager); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disconnected"})
}

// GetActive devuelve la conexión activa del usuario
func (h *ConnectionHandler) GetActive(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint)

	manager := c.Param("manager")
	allowed := map[string]bool{"pgsql": true, "oracle": true, "mysql": true, "mssql": true, "otro": true}
	if manager == "" {
		// list all active connections for the user
		if h.listActiveUC == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "list usecase not configured"})
			return
		}
		list, err := h.listActiveUC.Execute(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// transform to DTOs
		var resp []dto.ConnectionResponseDTO
		for _, conn := range list {
			var lastDisconnected string
			if conn.LastDisconnected != nil {
				lastDisconnected = conn.LastDisconnected.Format(time.RFC3339)
			}
			resp = append(resp, dto.ConnectionResponseDTO{
				ID:               conn.ID,
				UserID:           conn.UserID,
				Manager:          conn.Manager,
				Driver:           conn.Driver,
				Server:           conn.Server,
				DBUser:           conn.DBUser,
				IsConnected:      conn.IsConnected,
				LastConnected:    conn.LastConnected.Format(time.RFC3339),
				LastDisconnected: lastDisconnected,
			})
		}
		c.JSON(http.StatusOK, gin.H{"connections": resp})
		return
	}

	if !allowed[manager] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported manager"})
		return
	}

	conn, err := h.getActiveUC.Execute(c.Request.Context(), userID, manager)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if conn == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active connection"})
		return
	}

	var lastDisconnected string
	if conn.LastDisconnected != nil {
		lastDisconnected = conn.LastDisconnected.Format(time.RFC3339)
	}

	resp := dto.ConnectionResponseDTO{
		ID:               conn.ID,
		UserID:           conn.UserID,
		Driver:           conn.Driver,
		Server:           conn.Server,
		DBUser:           conn.DBUser,
		IsConnected:      conn.IsConnected,
		LastConnected:    conn.LastConnected.Format(time.RFC3339),
		LastDisconnected: lastDisconnected,
	}

	c.JSON(http.StatusOK, gin.H{"connection": resp})
}

// GetHistory devuelve logs de conexión del usuario filtrados por manager (si se pasa)
func (h *ConnectionHandler) GetHistory(c *gin.Context) {
	if h.listHistoryUC == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "history usecase not configured"})
		return
	}

	uid, _ := c.Get("userID")
	userID := uid.(uint)

	manager := c.Param("manager")
	// Accept limit/offset query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	list, err := h.listHistoryUC.Execute(c.Request.Context(), userID, connectionuc.HistoryFilters{Limit: limit, Offset: offset})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by manager if provided
	var out []dto.ConnectionLogDTO
	for _, l := range list {
		if manager != "" && manager != l.Driver {
			continue
		}
		out = append(out, dto.ConnectionLogDTO{
			ID:        l.ID,
			UserID:    l.UserID,
			Driver:    l.Driver,
			Server:    l.Server,
			DBUser:    l.DBUser,
			Timestamp: l.Timestamp.Format(time.RFC3339),
			Status:    l.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"logs": out})
}
