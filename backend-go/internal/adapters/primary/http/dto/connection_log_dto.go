package dto

// ConnectionLogDTO representa un registro de conexión retornado por la API
type ConnectionLogDTO struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Driver    string `json:"driver"`
	Server    string `json:"server"`
	DBUser    string `json:"db_user"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}
