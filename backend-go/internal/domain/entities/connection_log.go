package entities

import "time"

// ConnectionLog representa un registro de conexión/desconexión
type ConnectionLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Driver    string    `gorm:"not null" json:"driver"`
	Server    string    `gorm:"not null" json:"server"`
	DBUser    string    `gorm:"column:db_user;not null" json:"db_user"`
	Timestamp time.Time `gorm:"autoCreateTime;index" json:"timestamp"`
	Status    string    `gorm:"not null;index" json:"status"` // connected, disconnected, reconnected
}
