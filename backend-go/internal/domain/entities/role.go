package entities

// Role representa un rol de usuario en el sistema
type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"unique;not null"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

// HasPermission verifica si el rol tiene un permiso específico
func (r *Role) HasPermission(permission string) bool {
	for _, p := range r.Permissions {
		if p.Name == permission {
			return true
		}
	}
	return false
}

// UserRole representa la relación entre usuarios y roles
type UserRole struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`
	RoleID uint `json:"role_id" gorm:"primaryKey"`
}
