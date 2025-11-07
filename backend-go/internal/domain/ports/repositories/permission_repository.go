package repositories

import "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"

// PermissionRepository define las operaciones para gestionar permisos
type PermissionRepository interface {
	// Crear un nuevo permiso
	Create(permission *entities.Permission) error

	// Actualizar un permiso existente
	Update(permission *entities.Permission) error

	// Eliminar un permiso por ID
	Delete(id uint) error

	// Obtener un permiso por ID
	GetByID(id uint) (*entities.Permission, error)

	// Obtener un permiso por nombre
	GetByName(name string) (*entities.Permission, error)

	// Listar todos los permisos
	List() ([]entities.Permission, error)

	// Asignar permiso a rol
	AssignToRole(roleID uint, permissionID uint) error

	// Revocar permiso de rol
	RevokeFromRole(roleID uint, permissionID uint) error

	// Obtener permisos de rol
	GetRolePermissions(roleID uint) ([]entities.Permission, error)

	// Verificar si un rol tiene un permiso específico
	HasPermission(roleID uint, permissionName string) (bool, error)
}
