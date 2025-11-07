package entities

// Permission representa un permiso en el sistema
type Permission struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
	Resource    string `json:"resource" gorm:"not null"` // ej: "connections", "queries", "users"
	Action      string `json:"action" gorm:"not null"`   // ej: "create", "read", "update", "delete"
}

// PermissionKey devuelve la clave única del permiso en formato "resource:action"
func (p *Permission) PermissionKey() string {
	return p.Resource + ":" + p.Action
}
