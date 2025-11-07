package entities

import "time"

// User represents an application user
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:150;unique;not null;index" json:"username"`
	Email     string    `gorm:"size:254;unique;not null;index" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	FirstName string    `gorm:"size:150" json:"first_name"`
	LastName  string    `gorm:"size:150" json:"last_name"`
	Role      string    `gorm:"size:20;default:cliente;index" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	LastLogin time.Time `json:"last_login"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
}

// ActiveConnection represents a user's active SQL Server connection
type ActiveConnection struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Driver           string    `gorm:"size:255;not null" json:"driver"`
	Server           string    `gorm:"size:255;not null" json:"server"`
	DBUser           string    `gorm:"column:db_user;size:255;not null" json:"db_user"`
	Password         string    `gorm:"size:500;not null" json:"-"` // encrypted
	IsConnected      bool      `gorm:"default:false;index" json:"is_connected"`
	LastConnected    time.Time `json:"last_connected"`
	LastDisconnected time.Time `json:"last_disconnected"`
}

// ControlsInformation minimal entity for migration
type ControlsInformation struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Idx         int    `gorm:"not null;index" json:"idx"`
	Chapter     string `gorm:"size:10;not null;index" json:"chapter"`
	Name        string `gorm:"size:255;default:'Control Name'" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

// Additional entities will be added progressively in their own files
