package dto

// ConnectRequestDTO representa el payload para POST /api/db/{gestor}/open
type ConnectRequestDTO struct {
	Driver   string `json:"driver" binding:"required"`   // e.g. "mssql"
	Server   string `json:"server" binding:"required"`   // host or IP
	Port     string `json:"port" binding:"required"`     // port e.g. "1433"
	DBUser   string `json:"db_user" binding:"required"`  // DB username
	Password string `json:"password" binding:"required"` // DB password (plain-text over TLS expected)
}

// ConnectionResponseDTO shape returned to client (redact password)
type ConnectionResponseDTO struct {
	ID               uint   `json:"id"`
	UserID           uint   `json:"user_id"`
	Driver           string `json:"driver"`
	Server           string `json:"server"`
	DBUser           string `json:"db_user"`
	IsConnected      bool   `json:"is_connected"`
	LastConnected    string `json:"last_connected"`
	LastDisconnected string `json:"last_disconnected,omitempty"`
}
