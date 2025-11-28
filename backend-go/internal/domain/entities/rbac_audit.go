package entities

import "time"

// AdminActionLog represents an append-only audit record for admin/RBAC changes
type AdminActionLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ActorID    uint      `gorm:"index;not null" json:"actor_id"` // who performed the action
	ActorName  string    `gorm:"size:255" json:"actor_name"`
	Action     string    `gorm:"size:100;not null" json:"action"`       // e.g., role.create, permission.delete, role.assign
	TargetType string    `gorm:"size:50;not null" json:"target_type"`   // e.g., role, permission, user_role, role_permission
	TargetID   *uint     `gorm:"index" json:"target_id,omitempty"`      // id of target resource if available
	TargetName string    `gorm:"size:255" json:"target_name,omitempty"` // human friendly name
	Details    string    `gorm:"type:text" json:"details,omitempty"`    // JSON or free text with extra context
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
