package repositories

import "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"

// RoleRepository define las operaciones para gestionar roles
type RoleRepository interface {
	// Crear un nuevo rol
	Create(role *entities.Role) error

	// Actualizar un rol existente
	Update(role *entities.Role) error

	// Eliminar un rol por ID
	Delete(id uint) error

	// Obtener un rol por ID
	GetByID(id uint) (*entities.Role, error)

	// Obtener un rol por nombre
	GetByName(name string) (*entities.Role, error)

	// Listar todos los roles
	List() ([]entities.Role, error)

	// Asignar rol a usuario
	AssignToUser(userID uint, roleID uint) error

	// Revocar rol de usuario
	RevokeFromUser(userID uint, roleID uint) error

	// Obtener roles de usuario
	GetUserRoles(userID uint) ([]entities.Role, error)
}
