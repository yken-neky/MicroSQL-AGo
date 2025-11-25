package entities

import "time"

// Session represents an authenticated user session for the API
type Session struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	// store token as varchar with enough length for jwt to allow indexing in MySQL
	Token     string     `gorm:"type:varchar(512);not null;index" json:"token"`
	ExpiresAt *time.Time `gorm:"index" json:"expires_at"`
	IsActive  bool       `gorm:"default:true;index" json:"is_active"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
